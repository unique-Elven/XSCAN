package ember

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestShannonEntropyUniform256(t *testing.T) {
	data := make([]byte, 256)
	for i := range data {
		data[i] = byte(i)
	}
	var c [256]uint64
	for _, b := range data {
		c[b]++
	}
	h := ShannonEntropyFromCounts(&c, len(data))
	if math.Abs(h-8.0) > 1e-6 {
		t.Fatalf("expected entropy 8.0, got %v", h)
	}
}

func TestByteHistogramNormalizedSum(t *testing.T) {
	data := []byte{0x00, 0x00, 0xFF}
	c := ByteCounts256(data)
	raw := RawByteHistogram(&c)
	dst := make([]float32, ByteHistogramDim)
	ProcessByteHistogram(&raw, len(data), dst)
	var sum float32
	for _, v := range dst {
		sum += v
	}
	if math.Abs(float64(sum-1.0)) > 1e-5 {
		t.Fatalf("histogram should sum to 1, got %v", sum)
	}
	if dst[0] != 2.0/3.0 || dst[0xFF] != 1.0/3.0 {
		t.Fatalf("unexpected bins: [0]=%v [255]=%v", dst[0], dst[0xFF])
	}
}

func TestExtractFeaturesFromBytes_MinimalNonPE(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	vec, err := ExtractFeaturesFromBytes(data, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(vec) != FeatureDimFull {
		t.Fatalf("len %d want %d", len(vec), FeatureDimFull)
	}
	if vec[0] != 5 || vec[2] != 0 { // size, not PE
		t.Fatalf("general slice unexpected: %v", vec[:7])
	}
}

func TestExtractFeatures_SelfExecutable(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Skip(err)
	}
	vec, err := ExtractFeatures(exe, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(vec) != FeatureDimFull {
		t.Fatalf("len %d", len(vec))
	}
	var sum float32
	for _, v := range vec[offHist : offHist+DimHistogram] {
		sum += v
	}
	if math.Abs(float64(sum-1.0)) > 1e-3 {
		t.Fatalf("byte histogram should sum to ~1, got %v", sum)
	}
}

func TestExtractFeatures_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty.dat")
	if err := os.WriteFile(p, nil, 0600); err != nil {
		t.Fatal(err)
	}
	vec, err := ExtractFeatures(p, false)
	if err != nil {
		t.Fatal(err)
	}
	if vec[0] != 0 || vec[1] != 0 || vec[2] != 0 {
		t.Fatalf("empty file general: %v", vec[:3])
	}
}
