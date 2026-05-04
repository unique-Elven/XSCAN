package lgbm

import (
	"fmt"
	"sync"

	"github.com/dmitryikh/leaves"
)

// Engine loads one LightGBM text model and runs inference (pure Go via leaves).
type Engine struct {
	mu    sync.Mutex
	model *leaves.Ensemble
	path  string
}

// EnsureLoaded loads modelPath if not already the active file.
func (e *Engine) EnsureLoaded(modelPath string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if modelPath == "" {
		return fmt.Errorf("empty model path")
	}
	if e.model != nil && e.path == modelPath {
		return nil
	}
	m, err := ensembleFromPath(modelPath, true)
	if err != nil {
		return fmt.Errorf("load lightgbm model: %w", err)
	}
	e.model = m
	e.path = modelPath
	return nil
}

// LoadedPath returns the path of the currently loaded model, or empty if none.
func (e *Engine) LoadedPath() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.path
}

// NFeatures returns expected feature count for the loaded model, or 0 if none.
func (e *Engine) NFeatures() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.model == nil {
		return 0
	}
	return e.model.NFeatures()
}

// Predict runs inference after aligning the vector to the loaded model’s NFeatures()
// (padding with 0 or truncating), matching leaves’ expected input width.
func (e *Engine) Predict(features []float32) (float64, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.model == nil {
		return 0, fmt.Errorf("no lightgbm model loaded")
	}
	n := e.model.NFeatures()
	aligned := AlignFeatures(features, n)
	fvals := make([]float64, len(aligned))
	for i := range aligned {
		fvals[i] = float64(aligned[i])
	}
	score := e.model.PredictSingle(fvals, 0)
	return score, nil
}

// Unload releases the current model.
func (e *Engine) Unload() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.model = nil
	e.path = ""
}
