package ember

import (
	"strconv"

	"github.com/spaolacci/murmur3"
)

const sklearnHashSeed uint32 = 0

func murmurHash32Sklearn(key []byte) int32 {
	u := murmur3.Sum32WithSeed(key, sklearnHashSeed)
	return int32(u)
}

func sklearnBucketIndex(h int32, nFeatures int) int {
	const minInt32 = -2147483648
	if h == minInt32 {
		return int((2147483647 - int64(nFeatures-1)) % int64(nFeatures))
	}
	abs := int64(h)
	if abs < 0 {
		abs = -abs
	}
	return int(abs % int64(nFeatures))
}

// accumulatePairHashes sklearn FeatureHasher input_type="pair" with numeric value as weight.
func accumulatePairHashes(names []string, vals []float64, dim int, alternateSign bool, acc []float64) {
	for i := range names {
		key := []byte(names[i])
		h := murmurHash32Sklearn(key)
		idx := sklearnBucketIndex(h, dim)
		v := vals[i]
		if alternateSign && h < 0 {
			v = -v
		}
		acc[idx] += v
	}
}

// accumulateStringHashes sklearn FeatureHasher input_type="string", weight 1 per feature.
func accumulateStringHashes(features []string, dim int, alternateSign bool, acc []float64) {
	for _, s := range features {
		key := []byte(s)
		h := murmurHash32Sklearn(key)
		idx := sklearnBucketIndex(h, dim)
		v := 1.0
		if alternateSign && h < 0 {
			v = -1.0
		}
		acc[idx] += v
	}
}

func pairKeyStrInt(name string, val int) string {
	return name + "=" + strconv.Itoa(val)
}

func pairKeyStrUint32(name string, val uint32) string {
	return name + "=" + strconv.FormatUint(uint64(val), 10)
}

func denseToFloat32(acc []float64, dst []float32) {
	for i := range acc {
		dst[i] = float32(acc[i])
	}
}
