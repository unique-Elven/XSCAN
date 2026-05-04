package lgbm

// AlignFeatures resizes raw to exactly expectedDim for LightGBM inference:
// shorter vectors are zero-padded; longer vectors are truncated.
func AlignFeatures(raw []float32, expectedDim int) []float32 {
	if expectedDim <= 0 {
		return raw
	}
	switch {
	case len(raw) == expectedDim:
		return raw
	case len(raw) > expectedDim:
		return raw[:expectedDim]
	default:
		out := make([]float32, expectedDim)
		copy(out, raw)
		return out
	}
}
