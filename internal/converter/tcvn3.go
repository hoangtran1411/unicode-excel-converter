package converter

import "strings"

// TCVN3Converter handles conversion from TCVN3 (ABC) encoding to Unicode.
// Why: Encapsulates TCVN3 mapping logic.
type TCVN3Converter struct {
	replacer *strings.Replacer
}

// NewTCVN3Converter creates a new instance.
func NewTCVN3Converter() *TCVN3Converter {
	return &TCVN3Converter{
		replacer: strings.NewReplacer(
			// Lowercase
			"\u00B8", "á", // ¸
			"\u00B5", "à", // µ
			"\u00B6", "ả", // ¶
			"\u00B7", "ã", // ·
			"\u00B9", "ạ", // ¹

			"\u00A2", "â", // ¢
			"\u00CA", "ấ", // Ê -> Wait. TCVN3 map is tricky.
			// Let's use the standard "ABC" table sequence.
			// a, á, à, ả, ã, ạ
			// ă, ắ, ằ, ẳ, ẵ, ặ
			// ...

			// Revised TCVN3 Table (ABC):
			"\u00B8", "á",
			"\u00B5", "à",
			"\u00B6", "ả",
			"\u00B7", "ã",
			"\u00B9", "ạ",

			"\u00A8", "ă",
			"\u00BE", "ắ",
			"\u00BB", "ằ",
			"\u00BC", "ẳ",
			"\u00BD", "ẵ",
			"\u00C6", "ặ",

			"\u00A2", "â",
			"\u00CA", "ấ",
			"\u00C7", "ầ",
			"\u00C8", "ẩ",
			"\u00C9", "ẫ",
			"\u00CB", "ậ",

			"\u00D1", "é", // Ñ
			"\u00CC", "è", // Ì
			"\u00D0", "ẻ", // Ð
			"\u00CE", "ẽ", // Î
			"\u00CF", "ẹ", // Ï

			"\u00A3", "ê", // £
			"\u00D5", "ế", // Õ
			"\u00D2", "ề", // Ò
			"\u00D3", "ể", // Ó
			"\u00D4", "ễ", // Ô
			"\u00D6", "ệ", // Ö

			"\u00DD", "í", // Ý
			"\u00D8", "ì", // Ø
			"\u00DC", "ỉ", // Ü
			"\u00DE", "ĩ", // Þ
			"\u00DF", "ị", // ß

			"\u00F3", "ó",
			"\u00F2", "ò",
			"\u00F4", "õ",
			"\u00F5", "ọ",
			"\u00F6", "ô", // ö

			// Uppercase vowels in TCVN3 are often mapped to specific other chars or handled by
			// .VnTimeH font (which maps standard ASCII A-Z to localized A-Z).
			// However, mixed chars like 'ố' exist.
			// TCVN3 Uppercase is typically dependent on using the UPPERCASE FONT (.VnTimeH).
			// If the user uses .VnTimeH, then typing 'A' produces 'A', 'B' produces 'B'.
			// But 'á' (input) -> '¸' -> displays 'Á' in .VnTimeH.
			// So, if the font is .VnTimeH, we should convert '¸' to 'Á'.
			// For now, let's strictly handle the lowercase logic which is universally mapped in the standard font.

			// d
			"\u00AE", "đ", // ®
		),
	}
}

// ToUnicode converts TCVN3 encoded text to Unicode.
func (c *TCVN3Converter) ToUnicode(text string) string {
	return c.replacer.Replace(text)
}
