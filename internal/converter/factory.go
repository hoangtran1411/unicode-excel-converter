package converter

import "fmt"

// NewConverter creates a converter based on the encoding type.
// Returns an error for unsupported encodings instead of nil (idiomatic Go).
func NewConverter(encoding EncodingType) (Converter, error) {
	switch encoding {
	case EncodingVNI:
		return NewVNIConverter(), nil
	case EncodingTCVN3:
		return NewTCVN3Converter(), nil
	default:
		return nil, fmt.Errorf("unsupported encoding type: %s", encoding)
	}
}

// NoOpConverter is a pass-through converter that returns text unchanged.
// Useful for Unknown encoding to avoid nil handling.
type NoOpConverter struct{}

// ToUnicode returns the text unchanged (identity function).
func (NoOpConverter) ToUnicode(text string) string { return text }

// NewConverterOrNoop returns a converter, falling back to NoOpConverter for unknown types.
// Why: Prevents nil pointer dereference and simplifies caller code.
func NewConverterOrNoop(encoding EncodingType) Converter {
	c, err := NewConverter(encoding)
	if err != nil {
		return NoOpConverter{}
	}
	return c
}
