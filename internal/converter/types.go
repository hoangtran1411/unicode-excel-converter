// Package converter provides functions to convert legacy Vietnamese encodings to Unicode.
package converter

// EncodingType represents the source font encoding.
// Why: Using a typed string constant ensures type safety and prevents magic strings.
type EncodingType string

const (
	// EncodingVNI represents VNI-Windows encoding
	EncodingVNI EncodingType = "VNI"
	// EncodingTCVN3 represents TCVN3 (ABC) encoding
	EncodingTCVN3 EncodingType = "TCVN3"
	// EncodingAuto represents automatic encoding detection
	EncodingAuto EncodingType = "AUTO"
	// EncodingUnknown represents an unknown encoding
	EncodingUnknown EncodingType = "UNKNOWN"
)

// Converter is the interface that all specific encoding converters must implement.
// Why: Allows for polymorphism, making it easy to swap or add new converters without changing the consuming code.
type Converter interface {
	// ToUnicode converts the given legacy encoded string to a Unicode string.
	ToUnicode(text string) string
}
