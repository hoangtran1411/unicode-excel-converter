package converter

// NewConverterFactory creates a converter based on the encoding type.
// Why: Centralizes converter creation logic.
func NewConverterFactory(encoding EncodingType) Converter {
	switch encoding {
	case EncodingVNI:
		return NewVNIConverter()
	case EncodingTCVN3:
		return NewTCVN3Converter()
	default:
		return nil
	}
}
