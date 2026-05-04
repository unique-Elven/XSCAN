package lgbm

import (
	"fmt"
	"sync"

	"github.com/dmitryikh/leaves"
)

// ScanModelPool keeps Ember2018 and Ember2024 LightGBM ensembles in memory so scans do not
// re-read model files when switching routes. Reload is skipped when paths are unchanged.
type ScanModelPool struct {
	mu sync.RWMutex

	unsigned *leaves.Ensemble // no-cert / 2018 pipeline
	signed   *leaves.Ensemble // cert / 2024 pipeline

	pathUnsigned string
	pathSigned   string
}

// Reload loads both models from disk. On failure the previous pool state is left unchanged.
func (p *ScanModelPool) Reload(pathUnsigned, pathSigned string) error {
	if pathUnsigned == "" || pathSigned == "" {
		return fmt.Errorf("empty model path")
	}

	u, err := ensembleFromPath(pathUnsigned, true)
	if err != nil {
		return fmt.Errorf("load Ember2018 model: %w", err)
	}
	s, err := ensembleFromPath(pathSigned, true)
	if err != nil {
		return fmt.Errorf("load Ember2024 model: %w", err)
	}

	p.mu.Lock()
	p.unsigned, p.signed = u, s
	p.pathUnsigned, p.pathSigned = pathUnsigned, pathSigned
	p.mu.Unlock()
	return nil
}

// Ensure loads both models if missing or if paths differ from the active pair.
func (p *ScanModelPool) Ensure(pathUnsigned, pathSigned string) error {
	if pathUnsigned == "" || pathSigned == "" {
		return fmt.Errorf("scanner paths not configured")
	}
	p.mu.RLock()
	ok := p.unsigned != nil && p.signed != nil &&
		p.pathUnsigned == pathUnsigned && p.pathSigned == pathSigned
	p.mu.RUnlock()
	if ok {
		return nil
	}
	return p.Reload(pathUnsigned, pathSigned)
}

// Clear drops cached ensembles (e.g. when models are replaced on disk).
func (p *ScanModelPool) Clear() {
	p.mu.Lock()
	p.unsigned, p.signed = nil, nil
	p.pathUnsigned, p.pathSigned = "", ""
	p.mu.Unlock()
}

func predictAligned(m *leaves.Ensemble, raw []float32) (float64, error) {
	if m == nil {
		return 0, fmt.Errorf("lightgbm model not loaded")
	}
	expected := m.NFeatures()
	aligned := AlignFeatures(raw, expected)
	fvals := make([]float64, len(aligned))
	for i := range aligned {
		fvals[i] = float64(aligned[i])
	}
	return m.PredictSingle(fvals, 0), nil
}

// PredictUnsigned runs inference with the no-certificate (Ember2018) ensemble after aligning dims.
func (p *ScanModelPool) PredictUnsigned(raw []float32) (float64, error) {
	p.mu.RLock()
	m := p.unsigned
	p.mu.RUnlock()
	return predictAligned(m, raw)
}

// PredictSigned runs inference with the certificate-aware (Ember2024) ensemble after aligning dims.
func (p *ScanModelPool) PredictSigned(raw []float32) (float64, error) {
	p.mu.RLock()
	m := p.signed
	p.mu.RUnlock()
	return predictAligned(m, raw)
}

// NFeaturesUnsigned / NFeaturesSigned expose expected dims from leaves (0 if not loaded).
func (p *ScanModelPool) NFeaturesUnsigned() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.unsigned == nil {
		return 0
	}
	return p.unsigned.NFeatures()
}

func (p *ScanModelPool) NFeaturesSigned() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.signed == nil {
		return 0
	}
	return p.signed.NFeatures()
}
