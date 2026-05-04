package ember

import (
	"strconv"
	"strings"

	"github.com/saferwall/pe"

	_ "embed"
)

//go:embed pefile_warnings.txt
var pefileWarningsTxt string

var (
	warningPrefixes map[string]struct{}
	warningSuffixes map[string]struct{}
	warningLineIDs  map[string]int
)

// Ember DataDirectories._name_order (features.py).
var dataDirOrder = []string{
	"EXPORT",
	"IMPORT",
	"RESOURCE",
	"EXCEPTION",
	"SECURITY",
	"BASERELOC",
	"DEBUG",
	"COPYRIGHT",
	"GLOBALPTR",
	"TLS",
	"LOAD_CONFIG",
	"BOUND_IMPORT",
	"IAT",
	"DELAY_IMPORT",
	"COM_DESCRIPTOR",
	"RESERVED",
}

var dataDirNameIndex map[string]int

func init() {
	dataDirNameIndex = make(map[string]int, len(dataDirOrder))
	for i, n := range dataDirOrder {
		dataDirNameIndex[n] = i
	}

	warningPrefixes = make(map[string]struct{})
	warningSuffixes = make(map[string]struct{})
	warningLineIDs = make(map[string]int)
	lines := strings.Split(pefileWarningsTxt, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		warningLineIDs[line] = i
		if strings.HasPrefix(line, "...") {
			warningSuffixes[line[3:]] = struct{}{}
		} else {
			if strings.HasSuffix(line, "...") {
				warningPrefixes[line[:len(line)-3]] = struct{}{}
			}
		}
	}
}

func optionalDataDirectories(p *pe.File) *[16]pe.DataDirectory {
	switch oh := p.NtHeader.OptionalHeader.(type) {
	case pe.ImageOptionalHeader32:
		return &oh.DataDirectory
	case pe.ImageOptionalHeader64:
		return &oh.DataDirectory
	default:
		return nil
	}
}

func hasDynamicRelocs(p *pe.File) float32 {
	if p.LoadConfig.DVRT != nil && len(p.LoadConfig.DVRT.Entries) > 0 {
		return 1
	}
	return 0
}

// ProcessDataDirectories writes DataDirectories.dim (ember_cert). Matches pefile loop skipping RESERVED slot.
func ProcessDataDirectories(p *pe.File, dst []float32) {
	if len(dst) < DimDataDir {
		return
	}
	clear(dst[:DimDataDir])
	if p == nil {
		return
	}
	dd := optionalDataDirectories(p)
	if dd == nil {
		return
	}

	dst[32] = boolToFloat(p.HasReloc)
	dst[33] = hasDynamicRelocs(p)

	// Python uses range(1, len(raw_obj)-1) — skips last directory entry (RESERVED).
	for i := 0; i < 15; i++ {
		name := dataDirOrder[i]
		idx, ok := dataDirNameIndex[name]
		if !ok {
			continue
		}
		entry := dd[i]
		dst[2*idx] = float32(entry.Size)
		dst[2*idx+1] = float32(entry.VirtualAddress)
	}
}

func boolToFloat(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

// ProcessRichHeader writes RichHeader.dim (ember_cert).
func ProcessRichHeader(p *pe.File, dst []float32) {
	if len(dst) < DimRich {
		return
	}
	clear(dst[:DimRich])
	if p == nil || !p.HasRichHdr || len(p.RichHeader.CompIDs) == 0 {
		return
	}
	rh := p.RichHeader
	dst[0] = float32(len(rh.CompIDs))

	names := make([]string, len(rh.CompIDs))
	vals := make([]float64, len(rh.CompIDs))
	for i, c := range rh.CompIDs {
		names[i] = strconv.FormatUint(uint64(c.Unmasked), 10)
		vals[i] = float64(c.Count)
	}
	var acc [32]float64
	accumulatePairHashes(names, vals, 32, true, acc[:])
	denseToFloat32(acc[:], dst[1:])
}

func normalizeWarnings(anoms []string) []string {
	var out []string
	seen := make(map[string]struct{})
	for _, w := range anoms {
		found := false
		for suf := range warningSuffixes {
			if strings.HasSuffix(w, suf) {
				norm := "..." + suf
				if _, ok := seen[norm]; !ok {
					out = append(out, norm)
					seen[norm] = struct{}{}
				}
				found = true
				break
			}
		}
		if found {
			continue
		}
		for pre := range warningPrefixes {
			if strings.HasPrefix(w, pre) {
				norm := pre + "..."
				if _, ok := seen[norm]; !ok {
					out = append(out, norm)
					seen[norm] = struct{}{}
				}
				found = true
				break
			}
		}
		if !found {
			continue
		}
	}
	return out
}

// ProcessPEFormatWarnings writes PEFormatWarnings.dim (ember_cert).
func ProcessPEFormatWarnings(p *pe.File, dst []float32) {
	if len(dst) < DimWarnings {
		return
	}
	clear(dst[:DimWarnings])
	if p == nil {
		return
	}
	norms := normalizeWarnings(p.Anomalies)
	for _, w := range norms {
		if idx, ok := warningLineIDs[w]; ok {
			dst[idx] = 1
		}
	}
	dst[DimWarnings-1] = float32(len(norms))
}
