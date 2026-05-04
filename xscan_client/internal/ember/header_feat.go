package ember

import (
	"github.com/saferwall/pe"
)

// Categorical orders copied from ember_cert HeaderFileInfo (features.py).
var emberMachineOrder = []string{
	"IMAGE_FILE_MACHINE_UNKNOWN",
	"IMAGE_FILE_MACHINE_I386",
	"IMAGE_FILE_MACHINE_R3000",
	"IMAGE_FILE_MACHINE_R4000",
	"IMAGE_FILE_MACHINE_R10000",
	"IMAGE_FILE_MACHINE_WCEMIPSV2",
	"IMAGE_FILE_MACHINE_ALPHA",
	"IMAGE_FILE_MACHINE_SH3",
	"IMAGE_FILE_MACHINE_SH3DSP",
	"IMAGE_FILE_MACHINE_SH3E",
	"IMAGE_FILE_MACHINE_SH4",
	"IMAGE_FILE_MACHINE_SH5",
	"IMAGE_FILE_MACHINE_ARM",
	"IMAGE_FILE_MACHINE_THUMB",
	"IMAGE_FILE_MACHINE_ARMNT",
	"IMAGE_FILE_MACHINE_AM33",
	"IMAGE_FILE_MACHINE_POWERPC",
	"IMAGE_FILE_MACHINE_POWERPCFP",
	"IMAGE_FILE_MACHINE_IA64",
	"IMAGE_FILE_MACHINE_MIPS16",
	"IMAGE_FILE_MACHINE_ALPHA64",
	"IMAGE_FILE_MACHINE_AXP64",
	"IMAGE_FILE_MACHINE_MIPSFPU",
	"IMAGE_FILE_MACHINE_MIPSFPU16",
	"IMAGE_FILE_MACHINE_TRICORE",
	"IMAGE_FILE_MACHINE_CEF",
	"IMAGE_FILE_MACHINE_EBC",
	"IMAGE_FILE_MACHINE_RISCV32",
	"IMAGE_FILE_MACHINE_RISCV64",
	"IMAGE_FILE_MACHINE_RISCV128",
	"IMAGE_FILE_MACHINE_LOONGARCH32",
	"IMAGE_FILE_MACHINE_LOONGARCH64",
	"IMAGE_FILE_MACHINE_AMD64",
	"IMAGE_FILE_MACHINE_M32R",
	"IMAGE_FILE_MACHINE_ARM64",
	"IMAGE_FILE_MACHINE_CEE",
}

var emberSubsystemOrder = []string{
	"IMAGE_SUBSYSTEM_UNKNOWN",
	"IMAGE_SUBSYSTEM_NATIVE",
	"IMAGE_SUBSYSTEM_WINDOWS_GUI",
	"IMAGE_SUBSYSTEM_WINDOWS_CUI",
	"IMAGE_SUBSYSTEM_OS2_CUI",
	"IMAGE_SUBSYSTEM_POSIX_CUI",
	"IMAGE_SUBSYSTEM_NATIVE_WINDOWS",
	"IMAGE_SUBSYSTEM_WINDOWS_CE_GUI",
	"IMAGE_SUBSYSTEM_EFI_APPLICATION",
	"IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER",
	"IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER",
	"IMAGE_SUBSYSTEM_EFI_ROM",
	"IMAGE_SUBSYSTEM_XBOX",
	"IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION",
}

var coffFlagOrder = []struct {
	short string
	mask  uint16
}{
	{"RELOCS_STRIPPED", pe.ImageFileRelocsStripped},
	{"EXECUTABLE_IMAGE", pe.ImageFileExecutableImage},
	{"LINE_NUMS_STRIPPED", pe.ImageFileLineNumsStripped},
	{"LOCAL_SYMS_STRIPPED", pe.ImageFileLocalSymsStripped},
	{"AGGRESIVE_WS_TRIM", pe.ImageFileAggressiveWSTrim},
	{"LARGE_ADDRESS_AWARE", pe.ImageFileLargeAddressAware},
	{"16BIT_MACHINE", 0x0040},
	{"BYTES_REVERSED_LO", pe.ImageFileBytesReservedLow},
	{"32BIT_MACHINE", pe.ImageFile32BitMachine},
	{"DEBUG_STRIPPED", pe.ImageFileDebugStripped},
	{"REMOVABLE_RUN_FROM_SWAP", pe.ImageFileRemovableRunFromSwap},
	{"NET_RUN_FROM_SWAP", pe.ImageFileNetRunFromSwap},
	{"SYSTEM", pe.ImageFileSystem},
	{"DLL", pe.ImageFileDLL},
	{"UP_SYSTEM_ONLY", pe.ImageFileUpSystemOnly},
	{"BYTES_REVERSED_HI", pe.ImageFileBytesReservedHigh},
}

var dllFlagOrder = []struct {
	short string
	mask  uint16
}{
	{"HIGH_ENTROPY_VA", pe.ImageDllCharacteristicsHighEntropyVA},
	{"DYNAMIC_BASE", pe.ImageDllCharacteristicsDynamicBase},
	{"FORCE_INTEGRITY", pe.ImageDllCharacteristicsForceIntegrity},
	{"NX_COMPAT", pe.ImageDllCharacteristicsNXCompact},
	{"NO_ISOLATION", pe.ImageDllCharacteristicsNoIsolation},
	{"NO_SEH", pe.ImageDllCharacteristicsNoSEH},
	{"NO_BIND", pe.ImageDllCharacteristicsNoBind},
	{"APPCONTAINER", pe.ImageDllCharacteristicsAppContainer},
	{"WDM_DRIVER", pe.ImageDllCharacteristicsWdmDriver},
	{"GUARD_CF", pe.ImageDllCharacteristicsGuardCF},
	{"TERMINAL_SERVER_AWARE", pe.ImageDllCharacteristicsTerminalServiceAware},
}

var emberMachineCat map[string]int
var emberSubsystemCat map[string]int

// pefile MACHINE_TYPE numeric → IMAGE_FILE_MACHINE_* (last duplicate wins for 0x284).
var peMachineNumToName = map[uint16]string{
	0x0:    "IMAGE_FILE_MACHINE_UNKNOWN",
	0x014c: "IMAGE_FILE_MACHINE_I386",
	0x0160: "IMAGE_FILE_MACHINE_R3000BE",
	0x0162: "IMAGE_FILE_MACHINE_R3000",
	0x0166: "IMAGE_FILE_MACHINE_R4000",
	0x0168: "IMAGE_FILE_MACHINE_R10000",
	0x0169: "IMAGE_FILE_MACHINE_WCEMIPSV2",
	0x0184: "IMAGE_FILE_MACHINE_ALPHA",
	0x01a2: "IMAGE_FILE_MACHINE_SH3",
	0x01a3: "IMAGE_FILE_MACHINE_SH3DSP",
	0x01a4: "IMAGE_FILE_MACHINE_SH3E",
	0x01a6: "IMAGE_FILE_MACHINE_SH4",
	0x01a8: "IMAGE_FILE_MACHINE_SH5",
	0x01c0: "IMAGE_FILE_MACHINE_ARM",
	0x01c2: "IMAGE_FILE_MACHINE_THUMB",
	0x01c4: "IMAGE_FILE_MACHINE_ARMNT",
	0x01d3: "IMAGE_FILE_MACHINE_AM33",
	0x01f0: "IMAGE_FILE_MACHINE_POWERPC",
	0x01f1: "IMAGE_FILE_MACHINE_POWERPCFP",
	0x0200: "IMAGE_FILE_MACHINE_IA64",
	0x0266: "IMAGE_FILE_MACHINE_MIPS16",
	0x0284: "IMAGE_FILE_MACHINE_AXP64",
	0x0366: "IMAGE_FILE_MACHINE_MIPSFPU",
	0x0466: "IMAGE_FILE_MACHINE_MIPSFPU16",
	0x0520: "IMAGE_FILE_MACHINE_TRICORE",
	0x0cef: "IMAGE_FILE_MACHINE_CEF",
	0x0ebc: "IMAGE_FILE_MACHINE_EBC",
	0x5032: "IMAGE_FILE_MACHINE_RISCV32",
	0x5064: "IMAGE_FILE_MACHINE_RISCV64",
	0x5128: "IMAGE_FILE_MACHINE_RISCV128",
	0x6232: "IMAGE_FILE_MACHINE_LOONGARCH32",
	0x6264: "IMAGE_FILE_MACHINE_LOONGARCH64",
	0x8664: "IMAGE_FILE_MACHINE_AMD64",
	0x9041: "IMAGE_FILE_MACHINE_M32R",
	0xaa64: "IMAGE_FILE_MACHINE_ARM64",
	0xc0ee: "IMAGE_FILE_MACHINE_CEE",
}

var peSubsystemNumToName = map[uint16]string{
	0:  "IMAGE_SUBSYSTEM_UNKNOWN",
	1:  "IMAGE_SUBSYSTEM_NATIVE",
	2:  "IMAGE_SUBSYSTEM_WINDOWS_GUI",
	3:  "IMAGE_SUBSYSTEM_WINDOWS_CUI",
	5:  "IMAGE_SUBSYSTEM_OS2_CUI",
	7:  "IMAGE_SUBSYSTEM_POSIX_CUI",
	8:  "IMAGE_SUBSYSTEM_NATIVE_WINDOWS",
	9:  "IMAGE_SUBSYSTEM_WINDOWS_CE_GUI",
	10: "IMAGE_SUBSYSTEM_EFI_APPLICATION",
	11: "IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER",
	12: "IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER",
	13: "IMAGE_SUBSYSTEM_EFI_ROM",
	14: "IMAGE_SUBSYSTEM_XBOX",
	16: "IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION",
}

func init() {
	emberMachineCat = make(map[string]int, len(emberMachineOrder))
	for i, s := range emberMachineOrder {
		emberMachineCat[s] = i
	}
	emberSubsystemCat = make(map[string]int, len(emberSubsystemOrder))
	for i, s := range emberSubsystemOrder {
		emberSubsystemCat[s] = i
	}
}

func coffCharacteristicTags(ch uint16) []string {
	var out []string
	for _, row := range coffFlagOrder {
		if ch&row.mask == row.mask {
			out = append(out, row.short)
		}
	}
	return out
}

func dllCharacteristicTags(ch uint16) []string {
	var out []string
	for _, row := range dllFlagOrder {
		if ch&row.mask == row.mask {
			out = append(out, row.short)
		}
	}
	return out
}

func machineCategory(machine uint16) float32 {
	name, ok := peMachineNumToName[machine]
	if !ok {
		name = "IMAGE_FILE_MACHINE_UNKNOWN"
	}
	idx, ok := emberMachineCat[name]
	if !ok {
		return 0
	}
	return float32(idx)
}

func subsystemCategory(sub uint16) float32 {
	name, ok := peSubsystemNumToName[sub]
	if !ok {
		name = "IMAGE_SUBSYSTEM_UNKNOWN"
	}
	idx, ok := emberSubsystemCat[name]
	if !ok {
		return 0
	}
	return float32(idx)
}

// ProcessHeaderFileInfo writes HeaderFileInfo.dim floats (ember_cert).
func ProcessHeaderFileInfo(p *pe.File, dst []float32) {
	if len(dst) < DimHeader || p == nil {
		if len(dst) >= DimHeader {
			clear(dst[:DimHeader])
		}
		return
	}
	clear(dst[:DimHeader])

	fh := p.NtHeader.FileHeader
	ch := uint16(fh.Characteristics)

	var oh32 pe.ImageOptionalHeader32
	var oh64 pe.ImageOptionalHeader64
	switch oh := p.NtHeader.OptionalHeader.(type) {
	case pe.ImageOptionalHeader32:
		oh32 = oh
	case pe.ImageOptionalHeader64:
		oh64 = oh
	default:
		return
	}

	var (
		subsys                     uint16
		dllCh                      uint16
		majImg, minImg             uint16
		majLnk, minLnk             uint8
		majOS, minOS               uint16
		majSubV, minSubV           uint16
		sizeCode, sizeHdr, sizeImg uint32
		sizeInit, sizeUninit       uint32
		stackRes, stackCom         uint64
		heapRes, heapCom           uint64
		entry, baseCode            uint32
		imgBase                    uint64
		sectAlign                  uint32
		checksum                   uint32
		numRVA                     uint32
	)

	if p.Is64 {
		subsys = uint16(oh64.Subsystem)
		dllCh = uint16(oh64.DllCharacteristics)
		majImg, minImg = oh64.MajorImageVersion, oh64.MinorImageVersion
		majLnk, minLnk = oh64.MajorLinkerVersion, oh64.MinorLinkerVersion
		majOS, minOS = oh64.MajorOperatingSystemVersion, oh64.MinorOperatingSystemVersion
		majSubV, minSubV = oh64.MajorSubsystemVersion, oh64.MinorSubsystemVersion
		sizeCode = oh64.SizeOfCode
		sizeHdr = oh64.SizeOfHeaders
		sizeImg = oh64.SizeOfImage
		sizeInit = oh64.SizeOfInitializedData
		sizeUninit = oh64.SizeOfUninitializedData
		stackRes, stackCom = oh64.SizeOfStackReserve, oh64.SizeOfStackCommit
		heapRes, heapCom = oh64.SizeOfHeapReserve, oh64.SizeOfHeapCommit
		entry = oh64.AddressOfEntryPoint
		baseCode = oh64.BaseOfCode
		imgBase = oh64.ImageBase
		sectAlign = oh64.SectionAlignment
		checksum = oh64.CheckSum
		numRVA = oh64.NumberOfRvaAndSizes
	} else {
		subsys = uint16(oh32.Subsystem)
		dllCh = uint16(oh32.DllCharacteristics)
		majImg, minImg = oh32.MajorImageVersion, oh32.MinorImageVersion
		majLnk, minLnk = oh32.MajorLinkerVersion, oh32.MinorLinkerVersion
		majOS, minOS = oh32.MajorOperatingSystemVersion, oh32.MinorOperatingSystemVersion
		majSubV, minSubV = oh32.MajorSubsystemVersion, oh32.MinorSubsystemVersion
		sizeCode = oh32.SizeOfCode
		sizeHdr = oh32.SizeOfHeaders
		sizeImg = oh32.SizeOfImage
		sizeInit = oh32.SizeOfInitializedData
		sizeUninit = oh32.SizeOfUninitializedData
		stackRes, stackCom = uint64(oh32.SizeOfStackReserve), uint64(oh32.SizeOfStackCommit)
		heapRes, heapCom = uint64(oh32.SizeOfHeapReserve), uint64(oh32.SizeOfHeapCommit)
		entry = oh32.AddressOfEntryPoint
		baseCode = oh32.BaseOfCode
		imgBase = uint64(oh32.ImageBase)
		sectAlign = oh32.SectionAlignment
		checksum = oh32.CheckSum
		numRVA = oh32.NumberOfRvaAndSizes
	}

	coffTags := coffCharacteristicTags(ch)
	coffSet := make(map[string]struct{}, len(coffTags))
	for _, t := range coffTags {
		coffSet[t] = struct{}{}
	}
	dllTags := dllCharacteristicTags(dllCh)
	dllSet := make(map[string]struct{}, len(dllTags))
	for _, t := range dllTags {
		dllSet[t] = struct{}{}
	}

	i := 0
	put := func(v float32) {
		if i < DimHeader {
			dst[i] = v
			i++
		}
	}
	putF := func(v uint32) { put(float32(v)) }
	putU64 := func(v uint64) {
		put(float32(float64(v))) // match numpy float32 truncation for large values
	}

	putF(fh.TimeDateStamp)
	putF(uint32(fh.NumberOfSections))
	putF(fh.NumberOfSymbols)
	putF(uint32(fh.SizeOfOptionalHeader))
	putF(fh.PointerToSymbolTable)
	put(machineCategory(uint16(fh.Machine)))
	put(subsystemCategory(subsys))
	put(float32(majImg))
	put(float32(minImg))
	put(float32(majLnk))
	put(float32(minLnk))
	put(float32(majOS))
	put(float32(minOS))
	put(float32(majSubV))
	put(float32(minSubV))
	putF(sizeCode)
	putF(sizeHdr)
	putF(sizeImg)
	putF(sizeInit)
	putF(sizeUninit)
	putU64(stackRes)
	putU64(stackCom)
	putU64(heapRes)
	putU64(heapCom)
	putF(entry)
	putF(baseCode)
	putU64(imgBase)
	putF(sectAlign)
	putF(checksum)
	putF(numRVA)

	for _, row := range coffFlagOrder {
		if _, ok := coffSet[row.short]; ok {
			put(1)
		} else {
			put(0)
		}
	}
	for _, row := range dllFlagOrder {
		if _, ok := dllSet[row.short]; ok {
			put(1)
		} else {
			put(0)
		}
	}

	d := p.DOSHeader
	put(float32(d.Magic))
	put(float32(d.BytesOnLastPageOfFile))
	put(float32(d.PagesInFile))
	put(float32(d.Relocations))
	put(float32(d.SizeOfHeader))
	put(float32(d.MinExtraParagraphsNeeded))
	put(float32(d.MaxExtraParagraphsNeeded))
	put(float32(d.InitialSS))
	put(float32(d.InitialSP))
	put(float32(d.Checksum))
	put(float32(d.InitialIP))
	put(float32(d.InitialCS))
	put(float32(d.AddressOfRelocationTable))
	put(float32(d.OverlayNumber))
	put(float32(d.OEMIdentifier))
	put(float32(d.OEMInformation))
	put(float32(d.AddressOfNewEXEHeader))
}
