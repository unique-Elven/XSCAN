package ember

import (
	"github.com/saferwall/pe"
)

// OpenParsedPE parses a PE image from bytes for full Ember extraction.
// When security parsing is omitted, certificate-based features still get dimension placeholders.
func OpenParsedPE(data []byte, includeCert bool) (*pe.File, error) {
	opts := &pe.Options{
		Fast:                       false,
		SectionEntropy:             true,
		OmitSecurityDirectory:      !includeCert,
		DisableCertValidation:      true,
		DisableSignatureValidation: true,
	}
	p, err := pe.NewBytes(data, opts)
	if err != nil {
		return nil, err
	}
	if err := p.Parse(); err != nil {
		return nil, err
	}
	return p, nil
}

// OverlayBytes returns overlay bytes from an in-memory PE (saferwall Overlay() requires *os.File).
func OverlayBytes(p *pe.File, data []byte) []byte {
	if p == nil || len(data) == 0 {
		return nil
	}
	off := int(p.OverlayOffset)
	if off < 0 || off >= len(data) {
		return nil
	}
	return data[off:]
}

func entropyOfBytes(b []byte) float64 {
	if len(b) == 0 {
		return 0
	}
	var c [256]uint64
	for _, x := range b {
		c[x]++
	}
	return ShannonEntropyFromCounts(&c, len(b))
}
