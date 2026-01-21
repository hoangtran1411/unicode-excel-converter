package converter

import "strings"

// VNIConverter handles conversion from VNI-Windows encoding to Unicode.
// Why: Encapsulates the specific mapping logic for VNI legacy font.
type VNIConverter struct {
	replacer *strings.Replacer
}

// NewVNIConverter creates a new instance of VNIConverter.
// Why: Initializes the reusable replacer for efficient string transformation.
func NewVNIConverter() *VNIConverter {
	// Standard VNI-Windows mapping based on mapping VNI-specific byte values to Unicode.
	// When Excel reads VNI text, it typically sees the ANSI characters corresponding to the raw bytes.
	// E.g., 'á' in VNI is 0xE1. In Unicode/ANSI, 0xE1 is also 'á'.
	// But 'ả' in VNI is 0x1E. In Unicode, 0x1E is a Control Character.
	// We map these specific "garbage" characters to the correct Vietnamese Unicode.
	return &VNIConverter{
		replacer: strings.NewReplacer(
			// Lowercase
			"\u00E1", "á", // a1
			"\u00E0", "à", // a2
			"\u1E03", "ả", // a3 - Note: 1E03 is Unicode 'ả'. Wait. VNI byte is 0x1E.
			"\u001E", "ả", // Correct mapping for byte 0x1E
			"\u00E3", "ã", // a4
			"\u00B5", "ạ", // a5 - 0xB5 is µ

			"\u00E2", "â", // a6
			"\u00CA", "ấ", // a61
			"\u00C7", "ầ", // a62
			"\u00C8", "ẩ", // a63
			"\u00C9", "ẫ", // a64
			"\u00CB", "ậ", // a65

			"\u00E5", "ă", // a8
			"\u00D0", "ắ", // a81 - 0xD0 is Đ (Eth)
			"\u00CC", "ằ", // a82
			"\u00CE", "ẳ", // a83
			"\u00CF", "ẵ", // a84
			"\u00D1", "ặ", // a85 - 0xD1 is Ñ

			"\u00E9", "é", // e1
			"\u00E8", "è", // e2
			"\u001F", "ẻ", // e3 - 0x1F (Unit Separator)
			"\u00E4", "ẽ", // e4
			"\u00B9", "ẹ", // e5 - 0xB9 is ¹

			"\u00EA", "ê", // e6
			"\u00D5", "ế", // e61 - 0xD5 is Õ
			"\u00D2", "ề", // e62
			"\u00D3", "ể", // e63
			"\u00D4", "ễ", // e64
			"\u00D6", "ệ", // e65

			"\u00ED", "í", // i1
			"\u00EC", "ì", // i2
			"\u0021", "ỉ", // i3 - !
			"\u00EF", "ĩ", // i4
			"\u00F7", "ị", // i5 - Division sign

			"\u00F3", "ó", // o1
			"\u00F2", "ò", // o2
			"\u0023", "ỏ", // o3 - #
			"\u00F5", "õ", // o4
			"\u00F4", "ọ", // o5

			"\u00F6", "ô", // o6
			"\u0092", "ố", // o61 - CP1252 0x92 is ’ (Right single quote)
			"\u0093", "ồ", // o62 - “
			"\u0094", "ổ", // o63 - ”
			"\u0095", "ỗ", // o64 - •
			"\u0096", "ộ", // o65 - –

			"\u00F8", "ơ", // o7
			"\u00B6", "ớ", // o71 - ¶
			"\u00B7", "ờ", // o72 - ·
			"\u0080", "ở", // o73 - €
			"\u0081", "ỡ", // o74 - (Unused in 1252? mapped to specific control)
			"\u0082", "ợ", // o75 - ‚ (Single low-9 quote)

			"\u00FA", "ú", // u1
			"\u00F9", "ù", // u2
			"\u001A", "ủ", // u3 - SUB
			"\u0169", "ũ", // u4 - Wait, VNI u4 is 0xFB? 0xFB is û.
			// Correction needed for 'ũ'.
			// Standard VNI: u4 (ũ) -> 0xFB (û)
			"\u00FB", "ũ",

			"\u00D7", "ụ", // u5 - 0xD7 is ×

			"\u00FC", "ư", // u7
			"\u00BE", "ứ", // u71 - copy from previous
			"\u00BB", "ừ", // u72
			"\u00BC", "ử", // u73
			"\u00BD", "ữ", // u74
			"\u00AC", "ự", // u75

			"\u00FD", "ý", // y1
			"\u00D8", "ỳ", // y2 - Ø
			"\u002A", "ỷ", // y3 - *
			"\u00D9", "ỹ", // y4 - Ù
			"\u00DA", "ỵ", // y5 - Ú

			"\u00F1", "đ", // d9 - ñ
			"\u00AE", "Đ", // D9 - ®
		),
	}
}

func (c *VNIConverter) ToUnicode(text string) string {
	return c.replacer.Replace(text)
}
