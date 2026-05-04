package ember

import (
	"os"
	"sync"
)

// FeatureDimGeneralHistogram is the legacy prefix length (general ∥ histogram only).
const FeatureDimGeneralHistogram = GeneralFileInfoDim + ByteHistogramDim

// ExtractFeatures reads a file and returns the full Ember PE feature vector (~2568 floats).
//
// includeCert: when false, the security directory is omitted (faster); Authenticode block is mostly zeros.
func ExtractFeatures(filePath string, includeCert bool) ([]float32, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return ExtractFeaturesFromBytes(data, includeCert)
}

// ExtractFeaturesFromBytes computes the full Ember vector.
// Byte-level work (entropy histogram + string features) runs concurrently; PE-derived blocks run
// sequentially on *pe.File because github.com/saferwall/pe is not safe for concurrent use.
func ExtractFeaturesFromBytes(data []byte, includeCert bool) ([]float32, error) {
	out := make([]float32, FeatureDimFull)

	counts := ByteCounts256(data)
	size := len(data)
	ent := ShannonEntropyFromCounts(&counts, size)

	p, err := OpenParsedPE(data, includeCert)
	isPE := err == nil && p != nil

	ProcessGeneralFileInfo(RawGeneralFileInfo(data, ent, isPE), out[offGeneral:offGeneral+DimGeneral])
	rawHist := RawByteHistogram(&counts)
	ProcessByteHistogram(&rawHist, size, out[offHist:offHist+DimHistogram])

	var prepWG sync.WaitGroup
	prepWG.Add(2)
	go func() {
		defer prepWG.Done()
		ProcessByteEntropyHistogram(data, out[offByteEnt:offByteEnt+DimByteEntropy])
	}()
	go func() {
		defer prepWG.Done()
		ProcessStringExtractor(data, out[offStrings:offStrings+DimStrings])
	}()
	prepWG.Wait()

	if !isPE {
		DumpFeaturesIfEnv(out)
		return out, nil
	}

	ProcessHeaderFileInfo(p, out[offHeader:offHeader+DimHeader])
	ProcessSectionInfo(p, data, out[offSection:offSection+DimSection])
	ProcessImports(p, out[offImports:offImports+DimImports])
	ProcessExports(p, out[offExports:offExports+DimExports])
	ProcessDataDirectories(p, out[offDataDir:offDataDir+DimDataDir])
	ProcessRichHeader(p, out[offRich:offRich+DimRich])
	ProcessAuthenticodeSignature(p, includeCert, out[offAuth:offAuth+DimAuthenticode])
	ProcessPEFormatWarnings(p, out[offWarn:offWarn+DimWarnings])

	DumpFeaturesIfEnv(out)
	return out, nil
}
