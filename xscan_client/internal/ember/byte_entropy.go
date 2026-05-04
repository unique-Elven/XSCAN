package ember

import "math"

const byteEntropyWindow = 2048
const byteEntropyStep = 1024

// ByteEntropyHistogramDim is ByteEntropyHistogram.dim (ember_cert).
const ByteEntropyHistogramDim = 256

func entropyBinCounts(block []byte) (hbin int, counts [16]int) {
	w := len(block)
	if w == 0 {
		return 0, counts
	}
	for _, b := range block {
		counts[b>>4]++
	}
	var H float64
	wf := float64(w)
	for _, c := range counts {
		if c == 0 {
			continue
		}
		p := float64(c) / wf
		H -= p * math.Log2(p)
	}
	H *= 2
	hbin = int(H * 2)
	if hbin == 16 {
		hbin = 15
	}
	return hbin, counts
}

// ProcessByteEntropyHistogram matches ByteEntropyHistogram.process_raw_features (normalized counts).
func ProcessByteEntropyHistogram(data []byte, dst []float32) {
	var grid [16][16]int
	n := len(data)
	if n < byteEntropyWindow {
		hbin, c := entropyBinCounts(data)
		for j := 0; j < 16; j++ {
			grid[hbin][j] += c[j]
		}
	} else {
		lastStart := n - byteEntropyWindow
		for start := 0; start <= lastStart; start += byteEntropyStep {
			block := data[start : start+byteEntropyWindow]
			hbin, c := entropyBinCounts(block)
			for j := 0; j < 16; j++ {
				grid[hbin][j] += c[j]
			}
		}
	}
	var sum float64
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			sum += float64(grid[i][j])
		}
	}
	if sum <= 0 {
		return
	}
	inv := 1.0 / sum
	k := 0
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			dst[k] = float32(float64(grid[i][j]) * inv)
			k++
		}
	}
}
