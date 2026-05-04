package ember

// ByteHistogramDim is the output length of [ByteHistogram.ProcessRawFeatures] (ember_cert features.py).
const ByteHistogramDim = 256

// ByteHistogramRaw is the Python raw_features list: 256 integer counts (non-normalized).
type ByteHistogramRaw [256]uint64

// RawByteHistogram copies counts into ByteHistogramRaw (no allocation).
func RawByteHistogram(counts *[256]uint64) ByteHistogramRaw {
	return *(*ByteHistogramRaw)(counts)
}

// ProcessByteHistogram normalizes counts by the total byte count (same as numpy: counts / sum).
// If size == 0, dst is zero-filled.
func ProcessByteHistogram(raw *ByteHistogramRaw, size int, dst []float32) {
	if len(dst) < ByteHistogramDim {
		return
	}
	if size <= 0 {
		clear(dst[:ByteHistogramDim])
		return
	}
	inv := 1.0 / float32(size)
	for i := 0; i < ByteHistogramDim; i++ {
		dst[i] = float32(raw[i]) * inv
	}
}
