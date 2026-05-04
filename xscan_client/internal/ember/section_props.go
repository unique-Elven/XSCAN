package ember

import "github.com/saferwall/pe"

// sectionCharacteristicTags mirrors pefile.section_characteristics order (short names after IMAGE_SCN_),
// excluding TYPE_REG (mask 0) and ALIGN_MASK. Align entries compare (ch & 0x00F00000) == mask.
type sectionCharRow struct {
	short string
	mask  uint32
	align bool
}

var sectionCharTable = []sectionCharRow{
	{"TYPE_DSECT", 0x00000001, false},
	{"TYPE_NOLOAD", 0x00000002, false},
	{"TYPE_GROUP", 0x00000004, false},
	{"TYPE_NO_PAD", 0x00000008, false},
	{"TYPE_COPY", 0x00000010, false},
	{"CNT_CODE", pe.ImageSectionCntCode, false},
	{"CNT_INITIALIZED_DATA", pe.ImageSectionCntInitializedData, false},
	{"CNT_UNINITIALIZED_DATA", pe.ImageSectionCntUninitializedData, false},
	{"LNK_OTHER", pe.ImageSectionLnkOther, false},
	{"LNK_INFO", pe.ImageSectionLnkInfo, false},
	{"LNK_OVER", pe.ImageSectionReserved6, false},
	{"LNK_REMOVE", pe.ImageSectionLnkRemove, false},
	{"LNK_COMDAT", pe.ImageSectionLnkCOMDAT, false},
	{"MEM_PROTECTED", 0x00004000, false},
	{"NO_DEFER_SPEC_EXC", 0x00004000, false},
	{"GPREL", pe.ImageSectionGpRel, false},
	{"MEM_FARDATA", pe.ImageSectionGpRel, false}, // same mask as GPREL in pefile
	{"MEM_SYSHEAP", 0x00010000, false},
	{"MEM_PURGEABLE", pe.ImageSectionMemPurgeable, false},
	{"MEM_16BIT", pe.ImageSectionMem16Bit, false},
	{"MEM_LOCKED", pe.ImageSectionMemLocked, false},
	{"MEM_PRELOAD", pe.ImageSectionMemPreload, false},
	{"ALIGN_1BYTES", pe.ImageSectionAlign1Bytes, true},
	{"ALIGN_2BYTES", pe.ImageSectionAlign2Bytes, true},
	{"ALIGN_4BYTES", pe.ImageSectionAlign4Bytes, true},
	{"ALIGN_8BYTES", pe.ImageSectionAlign8Bytes, true},
	{"ALIGN_16BYTES", pe.ImageSectionAlign16Bytes, true},
	{"ALIGN_32BYTES", pe.ImageSectionAlign32Bytes, true},
	{"ALIGN_64BYTES", pe.ImageSectionAlign64Bytes, true},
	{"ALIGN_128BYTES", pe.ImageSectionAlign128Bytes, true},
	{"ALIGN_256BYTES", pe.ImageSectionAlign256Bytes, true},
	{"ALIGN_512BYTES", pe.ImageSectionAlign512Bytes, true},
	{"ALIGN_1024BYTES", pe.ImageSectionAlign1024Bytes, true},
	{"ALIGN_2048BYTES", pe.ImageSectionAlign2048Bytes, true},
	{"ALIGN_4096BYTES", pe.ImageSectionAlign4096Bytes, true},
	{"ALIGN_8192BYTES", pe.ImageSectionAlign8192Bytes, true},
	{"LNK_NRELOC_OVFL", pe.ImageSectionLnkNRelocOvfl, false},
	{"MEM_DISCARDABLE", pe.ImageSectionMemDiscardable, false},
	{"MEM_NOT_CACHED", pe.ImageSectionMemNotCached, false},
	{"MEM_NOT_PAGED", pe.ImageSectionMemNotPaged, false},
	{"MEM_SHARED", pe.ImageSectionMemShared, false},
	{"MEM_EXECUTE", pe.ImageSectionMemExecute, false},
	{"MEM_READ", pe.ImageSectionMemRead, false},
	{"MEM_WRITE", pe.ImageSectionMemWrite, false},
}

func sectionCharacteristicTags(ch uint32) []string {
	const alignMask = 0x00F00000
	am := ch & alignMask
	var out []string
	for _, row := range sectionCharTable {
		if row.align {
			if am == row.mask {
				out = append(out, row.short)
			}
			continue
		}
		if ch&row.mask == row.mask {
			out = append(out, row.short)
		}
	}
	return out
}
