package ember

import (
	"encoding/json"
	"os"
)

// DumpFeaturesIfEnv writes out when EMBER_DUMP_GO_FEATURES is set (path).
// If the value is exactly "1", writes ./go_features.json.
func DumpFeaturesIfEnv(out []float32) {
	p := os.Getenv("EMBER_DUMP_GO_FEATURES")
	if p == "" {
		return
	}
	if p == "1" {
		p = "go_features.json"
	}
	_ = WriteFeaturesJSON(p, out)
}

// WriteFeaturesJSON writes feat as a JSON array of numbers (e.g. [0.1, 0.0, 1.2, ...]).
func WriteFeaturesJSON(path string, feat []float32) error {
	out := make([]float64, len(feat))
	for i, v := range feat {
		out[i] = float64(v)
	}
	b, err := json.Marshal(out)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
