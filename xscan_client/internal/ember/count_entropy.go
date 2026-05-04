package ember

import "math"

// ByteCounts256 counts byte frequencies in a single pass (O(n), stack-allocated).
func ByteCounts256(data []byte) (counts [256]uint64) {
	for i := range data {
		counts[data[i]]++
	}
	return counts
}

// ShannonEntropyFromCounts computes base-2 Shannon entropy from a frequency table.
// size must equal the total number of samples (sum of counts). If size == 0, returns 0.
func ShannonEntropyFromCounts(counts *[256]uint64, size int) float64 {
	if size <= 0 {
		return 0
	}
	sz := float64(size)
	var h float64
	for i := 0; i < 256; i++ {
		c := counts[i]
		if c == 0 {
			continue
		}
		p := float64(c) / sz
		h -= p * math.Log2(p)
	}
	return h
}
