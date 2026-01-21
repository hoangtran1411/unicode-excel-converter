package converter

// EncodingType represents the source font encoding.
// Why: Using a typed string constant ensures type safety and prevents magic strings.
type EncodingType string

const (
	EncodingVNI    EncodingType = "VNI"
	EncodingTCVN3  EncodingType = "TCVN3"
	EncodingAuto   EncodingType = "AUTO"
	EncodingUnkown EncodingType = "UNKNOWN"
)

// Converter is the interface that all specific encoding converters must implement.
// Why: Allows for polymorphism, making it easy to swap or add new converters without changing the consuming code.
type Converter interface {
	// ToUnicode converts the given legacy encoded string to a Unicode string.
	ToUnicode(text string) string
}
