package ember

// parsePESuccess returns true when saferwall/pe can parse the buffer as a PE (full data dirs).
func parsePESuccess(data []byte, includeCert bool) bool {
	_, err := OpenParsedPE(data, includeCert)
	return err == nil
}
