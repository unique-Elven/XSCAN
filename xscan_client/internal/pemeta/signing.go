// Package pemeta classifies PE binaries for scan routing (signed vs unsigned pipeline).
package pemeta

import (
	"os"

	"github.com/saferwall/pe"
)

// ShouldUseSignedPipeline returns true when the sample should use the signed / cert-aware
// feature pipeline (Ember2024 + includeCert). It treats a non‑PE or unreadable file as unsigned.
func ShouldUseSignedPipeline(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	opts := &pe.Options{
		DisableCertValidation:      true,
		DisableSignatureValidation: true,
		OmitSecurityDirectory:      false,
	}
	p, err := pe.NewBytes(data, opts)
	if err != nil {
		return false
	}
	if err := p.Parse(); err != nil {
		return certificateTableLooksPresent(p)
	}
	if p.HasCertificate {
		return true
	}
	return certificateTableLooksPresent(p)
}

// ShouldUseSignedPipelineFromPath reads path and applies ShouldUseSignedPipeline.
func ShouldUseSignedPipelineFromPath(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return ShouldUseSignedPipeline(data)
}

func certificateTableLooksPresent(p *pe.File) bool {
	rva, size := certificateDataDirectory(p)
	if rva == 0 || size == 0 {
		return false
	}
	// Minimum size for a plausible WIN_CERTIFICATE header + payload stub
	if size < 8 {
		return false
	}
	return true
}

func certificateDataDirectory(p *pe.File) (uint32, uint32) {
	switch h := p.NtHeader.OptionalHeader.(type) {
	case *pe.ImageOptionalHeader32:
		if int(pe.ImageDirectoryEntryCertificate) < len(h.DataDirectory) {
			d := h.DataDirectory[pe.ImageDirectoryEntryCertificate]
			return d.VirtualAddress, d.Size
		}
	case *pe.ImageOptionalHeader64:
		if int(pe.ImageDirectoryEntryCertificate) < len(h.DataDirectory) {
			d := h.DataDirectory[pe.ImageDirectoryEntryCertificate]
			return d.VirtualAddress, d.Size
		}
	}
	return 0, 0
}
