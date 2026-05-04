package ember

import "github.com/spaolacci/murmur3"

// FeatureHashBucket maps a string key to a bucket in [0, dim), using MurmurHash3 32-bit.
// Intended for sklearn FeatureHasher-style bucketing in upcoming extractors (sections, imports, etc.).
// dim must be > 0.
func FeatureHashBucket(s string, dim int) int {
	if dim <= 0 {
		return 0
	}
	h := murmur3.Sum32([]byte(s))
	return int(uint32(h) % uint32(dim))
}

// FeatureHashBuckets accumulates signed counts into dst for a set of string features (no alternate_sign).
// Same string always maps to the same bucket; collisions sum.
func FeatureHashBuckets(features []string, dst []float32) {
	for _, s := range features {
		i := FeatureHashBucket(s, len(dst))
		dst[i] += 1
	}
}
