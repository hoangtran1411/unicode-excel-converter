package engine

import (
	"convert-vni-to-unicode/internal/converter"
	"strings"
)

// DetectEncoding attempts to identify the encoding based on font name and content.
// Why: Allows for "Auto" mode where the system guesses the encoding.
func DetectEncoding(fontName string, text string) converter.EncodingType {
	// 1. Check Font Name (Strongest indicator)
	if strings.HasPrefix(fontName, "VNI-") {
		return converter.EncodingVNI
	}
	if strings.HasPrefix(fontName, ".Vn") {
		return converter.EncodingTCVN3
	}

	// 2. Check content (Heuristic)
	// VNI uses unique chars like \u00D6 (Ö), \u00DB (Û), etc. found in vni.go
	if strings.ContainsAny(text, "\u00D6\u00DB\u00DC\u00DD\u00DE\u00DF\u00E6\u00E7\u00E8\u00E9") {
		return converter.EncodingVNI
	}

	// TCVN3 uses specific high-byte chars. Example: \u00B9, \u00AE, \u00A9 ...
	// Cöng ty -> 'ö' is \u00F6. 'ô' in TCVN3 is \u00F4.
	// Check for common TCVN3 vowels that differ from Unicode/VNI.
	// TCVN3 map: \u00F6 -> ô.
	if strings.ContainsAny(text, "\u00F6\u00F4\u00E2\u00EA\u00EE\u00B9") {
		return converter.EncodingTCVN3
	}

	return converter.EncodingUnkown
}
