package ember

import (
	"fmt"
	"strings"

	"github.com/saferwall/pe"
)

func addressOfEntryPoint(p *pe.File) uint32 {
	switch oh := p.NtHeader.OptionalHeader.(type) {
	case pe.ImageOptionalHeader32:
		return oh.AddressOfEntryPoint
	case pe.ImageOptionalHeader64:
		return oh.AddressOfEntryPoint
	default:
		return 0
	}
}

func sectionNameLower(sec pe.Section) string {
	raw := strings.Replace(string(sec.Header.Name[:]), "\x00", "", -1)
	return strings.ToLower(strings.TrimSpace(raw))
}

func sectionEntropy(sec *pe.Section, p *pe.File) float64 {
	if sec.Entropy != nil {
		return *sec.Entropy
	}
	return sec.CalculateEntropy(p)
}

func hasProp(props []string, tag string) bool {
	for _, p := range props {
		if p == tag {
			return true
		}
	}
	return false
}

func maxMinSlice(vals []float64) (maxV float64, minV float64) {
	if len(vals) == 0 {
		return 0, 0
	}
	maxV, minV = vals[0], vals[0]
	for _, v := range vals[1:] {
		if v > maxV {
			maxV = v
		}
		if v < minV {
			minV = v
		}
	}
	return maxV, minV
}

// ProcessSectionInfo writes SectionInfo.dim floats (ember_cert).
func ProcessSectionInfo(p *pe.File, fileBytes []byte, dst []float32) {
	if len(dst) < DimSection || p == nil {
		if len(dst) >= DimSection {
			clear(dst[:DimSection])
		}
		return
	}
	clear(dst[:DimSection])

	nfile := len(fileBytes)
	if nfile <= 0 {
		return
	}
	inv := 1.0 / float64(nfile)

	if len(p.Sections) == 0 {
		return
	}

	aoep := addressOfEntryPoint(p)
	entryName := ""
	for i := range p.Sections {
		if p.Sections[i].Contains(aoep, p) {
			entryName = sectionNameLower(p.Sections[i])
			break
		}
	}
	if entryName == "" {
		for i := range p.Sections {
			if p.Sections[i].Header.Characteristics&pe.ImageSectionMemExecute != 0 {
				entryName = sectionNameLower(p.Sections[i])
				break
			}
		}
	}

	type secRaw struct {
		name       string
		size       uint32
		vsize      uint32
		entropy    float64
		sizeRatio  float64
		vsizeRatio float64
		props      []string
	}
	raw := make([]secRaw, 0, len(p.Sections))
	for i := range p.Sections {
		s := &p.Sections[i]
		sz := s.Header.SizeOfRawData
		vs := s.Header.VirtualSize
		denomVs := max(vs, uint32(1))
		raw = append(raw, secRaw{
			name:       sectionNameLower(*s),
			size:       sz,
			vsize:      vs,
			entropy:    sectionEntropy(s, p),
			sizeRatio:  float64(sz) * inv,
			vsizeRatio: float64(sz) / float64(denomVs),
			props:      sectionCharacteristicTags(s.Header.Characteristics),
		})
	}

	ov := OverlayBytes(p, fileBytes)
	var ovSize int
	var ovRatio float64
	var ovEnt float64
	if len(ov) > 0 {
		ovSize = len(ov)
		ovRatio = float64(ovSize) * inv
		ovEnt = entropyOfBytes(ov)
	}

	nSec := len(raw)
	nZero := 0
	nEmptyName := 0
	nRX := 0
	nW := 0
	for _, s := range raw {
		if s.size == 0 {
			nZero++
		}
		if s.name == "" {
			nEmptyName++
		}
		if hasProp(s.props, "MEM_READ") && hasProp(s.props, "MEM_EXECUTE") {
			nRX++
		}
		if hasProp(s.props, "MEM_WRITE") {
			nW++
		}
	}

	entropies := make([]float64, 0, nSec+2)
	sizeRatios := make([]float64, 0, nSec+2)
	vsizeRatios := make([]float64, 0, nSec+1)
	for _, s := range raw {
		entropies = append(entropies, s.entropy)
		sizeRatios = append(sizeRatios, s.sizeRatio)
		vsizeRatios = append(vsizeRatios, s.vsizeRatio)
	}
	entropies = append(entropies, ovEnt, 0)
	sizeRatios = append(sizeRatios, ovRatio, 0)
	vsizeRatios = append(vsizeRatios, 0)

	maxE, minE := maxMinSlice(entropies)
	maxSR, minSR := maxMinSlice(sizeRatios)
	maxVR, minVR := maxMinSlice(vsizeRatios)

	general := []float64{
		float64(nSec),
		float64(nZero),
		float64(nEmptyName),
		float64(nRX),
		float64(nW),
		maxE,
		minE,
		maxSR,
		minSR,
		maxVR,
		minVR,
	}

	namesSizes := make([]string, len(raw))
	valsSizes := make([]float64, len(raw))
	namesVs := make([]string, len(raw))
	valsVs := make([]float64, len(raw))
	namesEnt := make([]string, len(raw))
	valsEnt := make([]float64, len(raw))
	for i := range raw {
		namesSizes[i] = raw[i].name
		valsSizes[i] = float64(raw[i].size)
		namesVs[i] = raw[i].name
		valsVs[i] = float64(raw[i].vsize)
		namesEnt[i] = raw[i].name
		valsEnt[i] = raw[i].entropy
	}

	var charFeatures []string
	for _, s := range raw {
		for _, pr := range s.props {
			charFeatures = append(charFeatures, fmt.Sprintf("%s:%s", s.name, pr))
		}
	}

	off := 0
	for _, g := range general {
		dst[off] = float32(g)
		off++
	}

	var buf50 [50]float64
	clear(buf50[:])
	accumulatePairHashes(namesSizes, valsSizes, 50, true, buf50[:])
	denseToFloat32(buf50[:], dst[off:off+50])
	off += 50

	clear(buf50[:])
	accumulatePairHashes(namesVs, valsVs, 50, true, buf50[:])
	denseToFloat32(buf50[:], dst[off:off+50])
	off += 50

	clear(buf50[:])
	accumulatePairHashes(namesEnt, valsEnt, 50, true, buf50[:])
	denseToFloat32(buf50[:], dst[off:off+50])
	off += 50

	clear(buf50[:])
	accumulateStringHashes(charFeatures, 50, true, buf50[:])
	denseToFloat32(buf50[:], dst[off:off+50])
	off += 50

	var buf10 [10]float64
	clear(buf10[:])
	accumulateStringHashes([]string{entryName}, 10, true, buf10[:])
	denseToFloat32(buf10[:], dst[off:off+10])
	off += 10

	dst[off] = float32(ovSize)
	off++
	dst[off] = float32(ovRatio)
	off++
	dst[off] = float32(ovEnt)
}
