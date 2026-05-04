package ember

import "math"

// GeneralFileInfoDim is the output length of [GeneralFileInfo.ProcessRawFeatures] (ember_cert features.py).
const GeneralFileInfoDim = 7

// GeneralFileInfoRaw matches Python GeneralFileInfo.raw_features JSON-able dict:
// size, entropy, is_pe, start_bytes[4].
type GeneralFileInfoRaw struct {
	Size       int       `json:"size"`
	Entropy    float64   `json:"entropy"`
	IsPE       int       `json:"is_pe"` // 0 or 1
	StartBytes [4]uint32 `json:"start_bytes"`
}

// RawGeneralFileInfo builds raw features from file bytes, precomputed entropy, and PE parse success.
func RawGeneralFileInfo(data []byte, entropy float64, isPE bool) GeneralFileInfoRaw {
	size := len(data)
	r := GeneralFileInfoRaw{
		Size:    size,
		Entropy: entropy,
		IsPE:    0,
	}
	if isPE {
		r.IsPE = 1
	}
	switch {
	case size >= 4:
		r.StartBytes[0] = uint32(data[0])
		r.StartBytes[1] = uint32(data[1])
		r.StartBytes[2] = uint32(data[2])
		r.StartBytes[3] = uint32(data[3])
	case size == 3:
		r.StartBytes[0] = uint32(data[0])
		r.StartBytes[1] = uint32(data[1])
		r.StartBytes[2] = uint32(data[2])
	case size == 2:
		r.StartBytes[0] = uint32(data[0])
		r.StartBytes[1] = uint32(data[1])
	case size == 1:
		r.StartBytes[0] = uint32(data[0])
	}
	return r
}

// ProcessGeneralFileInfo writes the 7-d float32 vector (numpy.hstack order from features.py).
func ProcessGeneralFileInfo(raw GeneralFileInfoRaw, dst []float32) {
	if len(dst) < GeneralFileInfoDim {
		return
	}
	dst[0] = safeFloat32FromInt(raw.Size)
	dst[1] = float32(raw.Entropy)
	dst[2] = float32(raw.IsPE)
	dst[3] = float32(raw.StartBytes[0])
	dst[4] = float32(raw.StartBytes[1])
	dst[5] = float32(raw.StartBytes[2])
	dst[6] = float32(raw.StartBytes[3])
}

func safeFloat32FromInt(v int) float32 {
	if v == 0 {
		return 0
	}
	f := float64(v)
	if f > math.MaxFloat32 {
		return float32(math.MaxFloat32)
	}
	if f < -math.MaxFloat32 {
		return float32(-math.MaxFloat32)
	}
	return float32(v)
}
