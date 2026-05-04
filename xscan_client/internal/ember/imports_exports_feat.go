package ember

import (
	"strconv"
	"strings"

	"github.com/saferwall/pe"
)

func clipRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// ImportLibrariesAndFQ builds the same lowercase-DLL list and fully-qualified import
// strings as Python ember_cert ImportsInfo.process_raw_features (pefile-based).
func ImportLibrariesAndFQ(p *pe.File) (libs []string, fq []string) {
	if p == nil || len(p.Imports) == 0 {
		return nil, nil
	}
	libSet := make(map[string]struct{})
	for _, imp := range p.Imports {
		dllOriginal := strings.TrimSpace(imp.Name)
		dll := strings.ToLower(dllOriginal)
		if dll == "" {
			continue
		}
		if _, ok := libSet[dll]; !ok {
			libSet[dll] = struct{}{}
			libs = append(libs, dll)
		}
		for _, fn := range imp.Functions {
			// saferwall sets both ByOrdinal and Name="#<ord>" for ordinal thunks — match ByOrdinal first
			// so we emit the same string as Python/pefile (`lower(dll):originalDll:ordinalN`).
			switch {
			case fn.ByOrdinal:
				fq = append(fq, dll+":"+dllOriginal+":ordinal"+strconv.FormatUint(uint64(fn.Ordinal), 10))
			case fn.Name != "":
				fq = append(fq, dll+":"+clipRunes(fn.Name, 10000))
			}
		}
	}
	return libs, fq
}

// ProcessImports writes ImportsInfo.dim floats (ember_cert).
func ProcessImports(p *pe.File, dst []float32) {
	if len(dst) < DimImports {
		return
	}
	clear(dst[:DimImports])
	if p == nil || len(p.Imports) == 0 {
		return
	}

	libs, fq := ImportLibrariesAndFQ(p)

	dst[0] = float32(len(fq))
	dst[1] = float32(len(libs))

	var libAcc [256]float64
	clear(libAcc[:])
	accumulateStringHashes(libs, 256, false, libAcc[:])
	denseToFloat32(libAcc[:], dst[2:2+256])

	var impAcc [1024]float64
	clear(impAcc[:])
	accumulateStringHashes(fq, 1024, false, impAcc[:])
	denseToFloat32(impAcc[:], dst[2+256:])
}

// ProcessExports writes ExportsInfo.dim floats (ember_cert).
func ProcessExports(p *pe.File, dst []float32) {
	if len(dst) < DimExports {
		return
	}
	clear(dst[:DimExports])
	if p == nil || len(p.Export.Functions) == 0 {
		return
	}

	var names []string
	for _, fn := range p.Export.Functions {
		switch {
		case fn.Name != "":
			names = append(names, clipRunes(fn.Name, 10000))
		default:
			names = append(names, "ordinal"+strconv.FormatUint(uint64(fn.Ordinal), 10))
		}
	}

	var acc [128]float64
	clear(acc[:])
	accumulateStringHashes(names, 128, true, acc[:])

	// Ember reproduces sklearn quirk: leading scalar is len(hash_vector)==128, not export count.
	dst[0] = 128
	denseToFloat32(acc[:], dst[1:])
}
