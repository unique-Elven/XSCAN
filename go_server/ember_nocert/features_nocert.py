#!/usr/bin/python
import re
import lief
import hashlib
import numpy as np
import os
import json
from sklearn.feature_extraction import FeatureHasher

# 针对 LIEF 0.14.0+ 优化的版本判定
LIEF_VERSION = [int(i) for i in lief.__version__.split('.')[:2]]
IS_NEW_LIEF = LIEF_VERSION[0] > 0 or LIEF_VERSION[1] >= 10

class FeatureType(object):
    name = ''
    dim = 0

    def __repr__(self):
        return '{}({})'.format(self.name, self.dim)

    def raw_features(self, bytez, lief_binary):
        raise (NotImplementedError)

    def process_raw_features(self, raw_obj):
        raise (NotImplementedError)

    def feature_vector(self, bytez, lief_binary):
        return self.process_raw_features(self.raw_features(bytez, lief_binary))

class ByteHistogram(FeatureType):
    name = 'histogram'
    dim = 256
    def __init__(self):
        super(ByteHistogram, self).__init__()

    def raw_features(self, bytez, lief_binary):
        counts = np.bincount(np.frombuffer(bytez, dtype=np.uint8), minlength=256)
        return counts.tolist()

    def process_raw_features(self, raw_obj):
        counts = np.array(raw_obj, dtype=np.float32)
        normalized = counts / counts.sum()
        return normalized

class ByteEntropyHistogram(FeatureType):
    name = 'byteentropy'
    dim = 256
    def __init__(self, step=1024, window=2048):
        super(ByteEntropyHistogram, self).__init__()
        self.window = window
        self.step = step

    def _entropy_bin_counts(self, block):
        c = np.bincount(block >> 4, minlength=16)
        p = c.astype(np.float32) / self.window
        wh = np.where(c)[0]
        H = np.sum(-p[wh] * np.log2(p[wh])) * 2
        Hbin = int(H * 2)
        if Hbin == 16: Hbin = 15
        return Hbin, c

    def raw_features(self, bytez, lief_binary):
        output = np.zeros((16, 16), dtype=int)
        a = np.frombuffer(bytez, dtype=np.uint8)
        if a.shape[0] < self.window:
            Hbin, c = self._entropy_bin_counts(a)
            output[Hbin, :] += c
        else:
            shape = a.shape[:-1] + (a.shape[-1] - self.window + 1, self.window)
            strides = a.strides + (a.strides[-1],)
            blocks = np.lib.stride_tricks.as_strided(a, shape=shape, strides=strides)[::self.step, :]
            for block in blocks:
                Hbin, c = self._entropy_bin_counts(block)
                output[Hbin, :] += c
        return output.flatten().tolist()

    def process_raw_features(self, raw_obj):
        counts = np.array(raw_obj, dtype=np.float32)
        return counts / counts.sum()

class SectionInfo(FeatureType):
    name = 'section'
    dim = 5 + 50 + 50 + 50 + 50 + 50
    def __init__(self):
        super(SectionInfo, self).__init__()

    @staticmethod
    def _properties(s):
        # 适配新版特性的获取方式
        return [str(c).split('.')[-1] for c in s.characteristics_lists]

    def raw_features(self, bytez, lief_binary):
        if lief_binary is None:
            return {"entry": "", "sections": []}

        try:
            # 0.14.0+ 统一使用这种方式获取入口点所在节
            rva = lief_binary.entrypoint - lief_binary.imagebase
            section = lief_binary.section_from_rva(rva)
            entry_section = section.name if section else ""
        except:
            entry_section = ""
            for s in lief_binary.sections:
                if lief.PE.SECTION_CHARACTERISTICS.MEM_EXECUTE in s.characteristics_lists:
                    entry_section = s.name
                    break

        return {
            "entry": entry_section,
            "sections": [{
                'name': s.name,
                'size': s.size,
                'entropy': s.entropy,
                'vsize': s.virtual_size,
                'props': self._properties(s)
            } for s in lief_binary.sections]
        }

    def process_raw_features(self, raw_obj):
        sections = raw_obj['sections']
        general = [
            len(sections),
            sum(1 for s in sections if s['size'] == 0),
            sum(1 for s in sections if s['name'] == ""),
            sum(1 for s in sections if 'MEM_READ' in s['props'] and 'MEM_EXECUTE' in s['props']),
            sum(1 for s in sections if 'MEM_WRITE' in s['props'])
        ]
        section_sizes = [(s['name'], s['size']) for s in sections]
        section_sizes_hashed = FeatureHasher(50, input_type="pair").transform([section_sizes]).toarray()[0]
        section_entropy = [(s['name'], s['entropy']) for s in sections]
        section_entropy_hashed = FeatureHasher(50, input_type="pair").transform([section_entropy]).toarray()[0]
        section_vsize = [(s['name'], s['vsize']) for s in sections]
        section_vsize_hashed = FeatureHasher(50, input_type="pair").transform([section_vsize]).toarray()[0]
        entry_name_hashed = FeatureHasher(50, input_type="string").transform([[raw_obj['entry']]]).toarray()[0]
        characteristics = [p for s in sections for p in s['props'] if s['name'] == raw_obj['entry']]
        characteristics_hashed = FeatureHasher(50, input_type="string").transform([characteristics]).toarray()[0]

        return np.hstack([general, section_sizes_hashed, section_entropy_hashed, section_vsize_hashed, entry_name_hashed, characteristics_hashed]).astype(np.float32)

class ImportsInfo(FeatureType):
    name = 'imports'
    dim = 1280
    def __init__(self):
        super(ImportsInfo, self).__init__()

    def raw_features(self, bytez, lief_binary):
        imports = {}
        if lief_binary is None: return imports
        # 适配 0.14.0+: .imports 已经是属性
        for lib in lief_binary.imports:
            if lib.name not in imports:
                imports[lib.name] = []
            for entry in lib.entries:
                if entry.is_ordinal:
                    imports[lib.name].append("ordinal" + str(entry.ordinal))
                else:
                    imports[lib.name].append(entry.name[:10000])
        return imports

    def process_raw_features(self, raw_obj):
        libraries = list(set([l.lower() for l in raw_obj.keys()]))
        libraries_hashed = FeatureHasher(256, input_type="string").transform([libraries]).toarray()[0]
        imports = [lib.lower() + ':' + e for lib, elist in raw_obj.items() for e in elist]
        imports_hashed = FeatureHasher(1024, input_type="string").transform([imports]).toarray()[0]
        return np.hstack([libraries_hashed, imports_hashed]).astype(np.float32)

class ExportsInfo(FeatureType):
    name = 'exports'
    dim = 128
    def __init__(self):
        super(ExportsInfo, self).__init__()

    def raw_features(self, bytez, lief_binary):
        if lief_binary is None: return []
        # 适配 0.14.0+: 导出现在是对象列表，使用 .name 属性
        clipped_exports = [exp.name[:10000] for exp in lief_binary.exported_functions]
        return clipped_exports

    def process_raw_features(self, raw_obj):
        return FeatureHasher(128, input_type="string").transform([raw_obj]).toarray()[0].astype(np.float32)

class GeneralFileInfo(FeatureType):
    name = 'general'
    dim = 10
    def __init__(self):
        super(GeneralFileInfo, self).__init__()

    def raw_features(self, bytez, lief_binary):
        if lief_binary is None:
            return {'size': len(bytez), 'vsize': 0, 'has_debug': 0, 'exports': 0, 'imports': 0, 'has_relocations': 0, 'has_resources': 0, 'has_signature': 0, 'has_tls': 0, 'symbols': 0}
        
        # 适配新版属性
        return {
            'size': len(bytez),
            'vsize': lief_binary.virtual_size,
            'has_debug': int(lief_binary.has_debug),
            'exports': len(lief_binary.exported_functions),
            'imports': len(lief_binary.imported_functions),
            'has_relocations': int(lief_binary.has_relocations),
            'has_resources': int(lief_binary.has_resources),
            'has_signature': int(lief_binary.has_signatures), # 0.14.0 统一使用复数形式
            'has_tls': int(lief_binary.has_tls),
            'symbols': len(lief_binary.symbols),
        }

    def process_raw_features(self, raw_obj):
        return np.asarray([raw_obj[k] for k in ['size', 'vsize', 'has_debug', 'exports', 'imports', 'has_relocations', 'has_resources', 'has_signature', 'has_tls', 'symbols']], dtype=np.float32)

class HeaderFileInfo(FeatureType):
    name = 'header'
    dim = 62
    def __init__(self):
        super(HeaderFileInfo, self).__init__()

    def raw_features(self, bytez, lief_binary):
        raw_obj = {'coff': {'timestamp': 0, 'machine': "", 'characteristics': []}, 'optional': {'subsystem': "", 'dll_characteristics': [], 'magic': "", 'major_image_version': 0, 'minor_image_version': 0, 'major_linker_version': 0, 'minor_linker_version': 0, 'major_operating_system_version': 0, 'minor_operating_system_version': 0, 'major_subsystem_version': 0, 'minor_subsystem_version': 0, 'sizeof_code': 0, 'sizeof_headers': 0, 'sizeof_heap_commit': 0}}
        if lief_binary is None: return raw_obj

        # 适配新版属性访问 (去掉了所有的括号)
        h = lief_binary.header
        opt = lief_binary.optional_header
        
        raw_obj['coff']['timestamp'] = h.time_date_stamps
        raw_obj['coff']['machine'] = str(h.machine).split('.')[-1]
        raw_obj['coff']['characteristics'] = [str(c).split('.')[-1] for c in h.characteristics_list]
        
        raw_obj['optional']['subsystem'] = str(opt.subsystem).split('.')[-1]
        raw_obj['optional']['dll_characteristics'] = [str(c).split('.')[-1] for c in opt.dll_characteristics_lists]
        raw_obj['optional']['magic'] = str(opt.magic).split('.')[-1]
        raw_obj['optional']['major_image_version'] = opt.major_image_version
        raw_obj['optional']['minor_image_version'] = opt.minor_image_version
        raw_obj['optional']['major_linker_version'] = opt.major_linker_version
        raw_obj['optional']['minor_linker_version'] = opt.minor_linker_version
        raw_obj['optional']['major_operating_system_version'] = opt.major_operating_system_version
        raw_obj['optional']['minor_operating_system_version'] = opt.minor_operating_system_version
        raw_obj['optional']['major_subsystem_version'] = opt.major_subsystem_version
        raw_obj['optional']['minor_subsystem_version'] = opt.minor_subsystem_version
        raw_obj['optional']['sizeof_code'] = opt.sizeof_code
        raw_obj['optional']['sizeof_headers'] = opt.sizeof_headers
        raw_obj['optional']['sizeof_heap_commit'] = opt.sizeof_heap_commit
        return raw_obj

    def process_raw_features(self, raw_obj):
        return np.hstack([
            raw_obj['coff']['timestamp'],
            FeatureHasher(10, input_type="string").transform([[raw_obj['coff']['machine']]]).toarray()[0],
            FeatureHasher(10, input_type="string").transform([raw_obj['coff']['characteristics']]).toarray()[0],
            FeatureHasher(10, input_type="string").transform([[raw_obj['optional']['subsystem']]]).toarray()[0],
            FeatureHasher(10, input_type="string").transform([raw_obj['optional']['dll_characteristics']]).toarray()[0],
            FeatureHasher(10, input_type="string").transform([[raw_obj['optional']['magic']]]).toarray()[0],
            raw_obj['optional']['major_image_version'], raw_obj['optional']['minor_image_version'],
            raw_obj['optional']['major_linker_version'], raw_obj['optional']['minor_linker_version'],
            raw_obj['optional']['major_operating_system_version'], raw_obj['optional']['minor_operating_system_version'],
            raw_obj['optional']['major_subsystem_version'], raw_obj['optional']['minor_subsystem_version'],
            raw_obj['optional']['sizeof_code'], raw_obj['optional']['sizeof_headers'], raw_obj['optional']['sizeof_heap_commit'],
        ]).astype(np.float32)

class StringExtractor(FeatureType):
    name = 'strings'
    dim = 1 + 1 + 1 + 96 + 1 + 1 + 1 + 1 + 1
    def __init__(self):
        super(StringExtractor, self).__init__()
        self._allstrings = re.compile(b'[\x20-\x7f]{5,}')
        self._paths = re.compile(b'c:\\\\', re.IGNORECASE)
        self._urls = re.compile(b'https?://', re.IGNORECASE)
        self._registry = re.compile(b'HKEY_')
        self._mz = re.compile(b'MZ')

    def raw_features(self, bytez, lief_binary):
        allstrings = self._allstrings.findall(bytez)
        if allstrings:
            string_lengths = [len(s) for s in allstrings]
            avlength = sum(string_lengths) / len(string_lengths)
            as_shifted_string = [b - ord(b'\x20') for b in b''.join(allstrings)]
            c = np.bincount(as_shifted_string, minlength=96)
            csum = c.sum()
            p = c.astype(np.float32) / csum
            wh = np.where(c)[0]
            H = np.sum(-p[wh] * np.log2(p[wh]))
        else:
            avlength = 0; c = np.zeros((96,), dtype=np.float32); H = 0; csum = 0
        return {'numstrings': len(allstrings), 'avlength': avlength, 'printabledist': c.tolist(), 'printables': int(csum), 'entropy': float(H), 'paths': len(self._paths.findall(bytez)), 'urls': len(self._urls.findall(bytez)), 'registry': len(self._registry.findall(bytez)), 'MZ': len(self._mz.findall(bytez))}

    def process_raw_features(self, raw_obj):
        hist_divisor = float(raw_obj['printables']) if raw_obj['printables'] > 0 else 1.0
        return np.hstack([raw_obj['numstrings'], raw_obj['avlength'], raw_obj['printables'], np.asarray(raw_obj['printabledist']) / hist_divisor, raw_obj['entropy'], raw_obj['paths'], raw_obj['urls'], raw_obj['registry'], raw_obj['MZ']]).astype(np.float32)

class DataDirectories(FeatureType):
    name = 'datadirectories'
    dim = 15 * 2
    def __init__(self):
        super(DataDirectories, self).__init__()
        self._name_order = ["EXPORT_TABLE", "IMPORT_TABLE", "RESOURCE_TABLE", "EXCEPTION_TABLE", "CERTIFICATE_TABLE", "BASE_RELOCATION_TABLE", "DEBUG", "ARCHITECTURE", "GLOBAL_PTR", "TLS_TABLE", "LOAD_CONFIG_TABLE", "BOUND_IMPORT", "IAT", "DELAY_IMPORT_DESCRIPTOR", "CLR_RUNTIME_HEADER"]

    def raw_features(self, bytez, lief_binary):
        output = []
        if lief_binary is None: return output
        for data_directory in lief_binary.data_directories:
            output.append({"name": str(data_directory.type).replace("DATA_DIRECTORY.", ""), "size": data_directory.size, "virtual_address": data_directory.rva})
        return output

    def process_raw_features(self, raw_obj):
        features = np.zeros(2 * len(self._name_order), dtype=np.float32)
        for i in range(min(len(self._name_order), len(raw_obj))):
            features[2 * i] = raw_obj[i]["size"]
            features[2 * i + 1] = raw_obj[i]["virtual_address"]
        return features

class PEFeatureExtractor(object):
    def __init__(self, feature_version=2, print_feature_warning=False, features_file=''):
        self.features = []
        # 注意：此处移除了所有关于版本比对的 Warning 打印逻辑
        struct = {
            'ByteHistogram': ByteHistogram(),
            'ByteEntropyHistogram': ByteEntropyHistogram(),
            'StringExtractor': StringExtractor(),
            'GeneralFileInfo': GeneralFileInfo(),
            'HeaderFileInfo': HeaderFileInfo(),
            'SectionInfo': SectionInfo(),
            'ImportsInfo': ImportsInfo(),
            'ExportsInfo': ExportsInfo()
        }
        if os.path.exists(features_file):
            with open(features_file, encoding='utf8') as f:
                x = json.load(f)
                self.features = [struct[fe] for fe in x['features'] if fe in struct]
        else:
            self.features = list(struct.values())

        if feature_version == 2:
            self.features.append(DataDirectories())
        self.dim = sum([fe.dim for fe in self.features])

    def raw_features(self, bytez):
        # 适配 0.14.0 的解析器异常
        try:
            # 现代 LIEF 推荐直接传入 bytes 对象
            lief_binary = lief.PE.parse(list(bytez) if isinstance(bytez, bytes) else bytez)
        except Exception as e:
            lief_binary = None

        features = {"sha256": hashlib.sha256(bytez).hexdigest()}
        features.update({fe.name: fe.raw_features(bytez, lief_binary) for fe in self.features})
        return features

    def process_raw_features(self, raw_obj):
        feature_vectors = [fe.process_raw_features(raw_obj[fe.name]) for fe in self.features]
        return np.hstack(feature_vectors).astype(np.float32)

    def feature_vector(self, bytez):
        return self.process_raw_features(self.raw_features(bytez))