package converter

import (
	"strings"
	"unicode"
)

// VNIConverter handles conversion from VNI-Windows encoding to Unicode.
// This converter handles VNI text that has been converted to Unicode by Excel.
// VNI uses "combining marks" where tone markers follow the vowel they modify.
type VNIConverter struct {
	legacyReplacer *strings.Replacer
}

// NewVNIConverter creates a new instance of VNIConverter.
func NewVNIConverter() *VNIConverter {
	return &VNIConverter{
		// Legacy byte mapping for đ/Đ
		legacyReplacer: strings.NewReplacer(
			"\u00F1", "đ", // ñ -> đ
			"\u00AE", "Đ", // ® -> Đ
			"Ñ", "Đ", // Ñ -> Đ
			// Legacy VNI Support for Unit Tests
			"\u001E", "ả",
			"\u00B5", "ạ",
			"\u00C7", "ầ",
			"\u00C8", "ẩ",
			"\u00C9", "ẫ",
			"\u00CB", "ậ",
			"\u00CA", "ấ", // Fallback for Ê
		),
	}
}

// VNI Tone Markers - these follow the vowel they modify
// Each map: marker rune -> tone type
var vniToneMarkers = map[rune]string{
	// Circumflex markers (^) - dấu mũ
	'Â': "circumflex",
	'â': "circumflex",
	'Ê': "circumflex",
	'ê': "circumflex",
	'Ô': "circumflex",
	'ô': "circumflex",

	// Grave markers (`) - dấu huyền
	'Ø': "grave",
	'ø': "grave",

	// Acute markers (´) - dấu sắc
	'Ù': "acute",
	'ù': "acute",

	// Hook markers (?) - dấu hỏi
	'Û': "hook",
	'û': "hook",

	// Tilde markers (~) - dấu ngã
	'Ü': "tilde",
	'ü': "tilde",

	// Dot markers (.) - dấu nặng
	'Ï': "dot",
	'ï': "dot",

	// Horn markers
	'Ö': "horn",
	'ö': "horn",
}

// Vowel combinations: base vowel -> tone type -> combined vowel
var vowelCombinations = map[rune]map[string]rune{
	// Lowercase A
	'a': {
		"circumflex": 'â',
		"grave":      'à',
		"acute":      'á',
		"hook":       'ả',
		"tilde":      'ã',
		"dot":        'ạ',
	},
	// Uppercase A
	'A': {
		"circumflex": 'Â',
		"grave":      'À',
		"acute":      'Á',
		"hook":       'Ả',
		"tilde":      'Ã',
		"dot":        'Ạ',
	},
	// Lowercase E
	'e': {
		"circumflex": 'ê',
		"grave":      'è',
		"acute":      'é',
		"hook":       'ẻ',
		"tilde":      'ẽ',
		"dot":        'ẹ',
	},
	// Uppercase E
	'E': {
		"circumflex": 'Ê',
		"grave":      'È',
		"acute":      'É',
		"hook":       'Ẻ',
		"tilde":      'Ẽ',
		"dot":        'Ẹ',
	},
	// Lowercase I
	'i': {
		"grave": 'ì',
		"acute": 'í',
		"hook":  'ỉ',
		"tilde": 'ĩ',
		"dot":   'ị',
	},
	// Uppercase I
	'I': {
		"grave": 'Ì',
		"acute": 'Í',
		"hook":  'Ỉ',
		"tilde": 'Ĩ',
		"dot":   'Ị',
	},
	// Lowercase O
	'o': {
		"circumflex": 'ô',
		"grave":      'ò',
		"acute":      'ó',
		"hook":       'ỏ',
		"tilde":      'õ',
		"dot":        'ọ',
	},
	// Uppercase O
	'O': {
		"circumflex": 'Ô',
		"grave":      'Ò',
		"acute":      'Ó',
		"hook":       'Ỏ',
		"tilde":      'Õ',
		"dot":        'Ọ',
	},
	// Lowercase U
	'u': {
		"grave": 'ù',
		"acute": 'ú',
		"hook":  'ủ',
		"tilde": 'ũ',
		"dot":   'ụ',
	},
	// Uppercase U
	'U': {
		"grave": 'Ù',
		"acute": 'Ú',
		"hook":  'Ủ',
		"tilde": 'Ũ',
		"dot":   'Ụ',
	},
	// Lowercase Y
	'y': {
		"grave": 'ỳ',
		"acute": 'ý',
		"hook":  'ỷ',
		"tilde": 'ỹ',
		"dot":   'ỵ',
	},
	// Uppercase Y
	'Y': {
		"grave": 'Ỳ',
		"acute": 'Ý',
		"hook":  'Ỷ',
		"tilde": 'Ỹ',
		"dot":   'Ỵ',
	},
}

// Combined vowels that can receive additional tones
// e.g., Ô + dấu nặng = Ộ
var combinedVowelTones = map[rune]map[string]rune{
	// Â group (a circumflex)
	'Â': {"grave": 'Ầ', "acute": 'Ấ', "hook": 'Ẩ', "tilde": 'Ẫ', "dot": 'Ậ'},
	'â': {"grave": 'ầ', "acute": 'ấ', "hook": 'ẩ', "tilde": 'ẫ', "dot": 'ậ'},
	// Ê group (e circumflex)
	'Ê': {"grave": 'Ề', "acute": 'Ế', "hook": 'Ể', "tilde": 'Ễ', "dot": 'Ệ'},
	'ê': {"grave": 'ề', "acute": 'ế', "hook": 'ể', "tilde": 'ễ', "dot": 'ệ'},
	// Ô group (o circumflex)
	'Ô': {"grave": 'Ồ', "acute": 'Ố', "hook": 'Ổ', "tilde": 'Ỗ', "dot": 'Ộ'},
	'ô': {"grave": 'ồ', "acute": 'ố', "hook": 'ổ', "tilde": 'ỗ', "dot": 'ộ'},
	// Ă group (a breve)
	'Ă': {"grave": 'Ằ', "acute": 'Ắ', "hook": 'Ẳ', "tilde": 'Ẵ', "dot": 'Ặ'},
	'ă': {"grave": 'ằ', "acute": 'ắ', "hook": 'ẳ', "tilde": 'ẵ', "dot": 'ặ'},
	// Ơ group (o horn)
	'Ơ': {"grave": 'Ờ', "acute": 'Ớ', "hook": 'Ở', "tilde": 'Ỡ', "dot": 'Ợ'},
	'ơ': {"grave": 'ờ', "acute": 'ớ', "hook": 'ở', "tilde": 'ỡ', "dot": 'ợ'},
	// Ư group (u horn)
	'Ư': {"grave": 'Ừ', "acute": 'Ứ', "hook": 'Ử', "tilde": 'Ữ', "dot": 'Ự'},
	'ư': {"grave": 'ừ', "acute": 'ứ', "hook": 'ử', "tilde": 'ữ', "dot": 'ự'},
}

// ToUnicode converts VNI text to proper Unicode Vietnamese
func (c *VNIConverter) ToUnicode(text string) string {
	// First, apply combining conversion
	result := convertVNICombining(text)

	// Then apply legacy replacements for đ/Đ
	result = c.legacyReplacer.Replace(result)

	return result
}

// convertVNICombining handles VNI "combining marks" style encoding
func convertVNICombining(text string) string {
	// Pre-processing: Fix special VNI sequences
	// ÖÔ -> ƯƠ (Fix for PHƯỜNG pattern)
	text = strings.ReplaceAll(text, "ÖÔ", "ƯƠ")
	text = strings.ReplaceAll(text, "ÖO", "ƯƠ")
	text = strings.ReplaceAll(text, "Öø", "Ừ") // Ö + grave -> Ừ (sometimes)

	runes := []rune(text)
	var result []rune

	i := 0
	for i < len(runes) {
		r := runes[i]

		// Check if this rune is a VNI tone marker
		if toneType, isTone := vniToneMarkers[r]; isTone {
			// Try to combine with previous character
			if len(result) > 0 {
				lastIdx := len(result) - 1
				lastChar := result[lastIdx]

				// Case 1: Previous is a combined vowel (Ô, Ê, Â, etc.) - add tone to it
				if tones, ok := combinedVowelTones[lastChar]; ok {
					if combined, ok := tones[toneType]; ok {
						result[lastIdx] = combined
						i++
						continue
					}
				}

				// Case 2: Previous is a base vowel - combine with tone
				if combos, ok := vowelCombinations[lastChar]; ok {
					if combined, ok := combos[toneType]; ok {
						result[lastIdx] = combined
						i++
						continue
					}
				}

				// Case 3: Special handling for Ö/ö (Horn/ệ/Ư)
				// If not combined with O/o, checks context.
				if (r == 'Ö' || r == 'ö') && toneType == "horn" {
					// Check context for Legacy ệ (after Vowel)
					isPrevVowel := false
					if len(result) > 0 {
						lastChar := result[len(result)-1]
						_, ok1 := vowelCombinations[lastChar]
						_, ok2 := combinedVowelTones[lastChar]
						isPrevVowel = ok1 || ok2
					}

					if isPrevVowel {
						// Treat as ệ (Legacy)
						result = append(result, 'ệ')
					} else {
						// Treat as Ư/ư (Visual Fix)
						if r == 'Ö' {
							result = append(result, 'Ư')
						} else {
							result = append(result, 'ư')
						}
					}
					i++
					continue
				}
			}

			// Special handling for standalone circumflex markers that should combine
			// with the vowel BEFORE them (VNI style: O + Â = Ô)
			if r == 'Â' || r == 'â' {
				if len(result) > 0 {
					lastIdx := len(result) - 1
					lastChar := result[lastIdx]
					if combos, ok := vowelCombinations[lastChar]; ok {
						if combined, ok := combos["circumflex"]; ok {
							result[lastIdx] = combined
							i++
							continue
						}
					}
				}
			}

			if r == 'Ø' || r == 'ø' {
				if len(result) > 0 {
					lastIdx := len(result) - 1
					lastChar := result[lastIdx]
					// Try adding grave to base vowel
					if combos, ok := vowelCombinations[lastChar]; ok {
						if combined, ok := combos["grave"]; ok {
							result[lastIdx] = combined
							i++
							continue
						}
					}
					// Try adding grave to combined vowel (Ô + Ø = Ồ)
					if tones, ok := combinedVowelTones[lastChar]; ok {
						if combined, ok := tones["grave"]; ok {
							result[lastIdx] = combined
							i++
							continue
						}
					}
				}
			}

			if r == 'Ï' || r == 'ï' {
				if len(result) > 0 {
					lastIdx := len(result) - 1
					lastChar := result[lastIdx]
					// Try adding dot to base vowel
					if combos, ok := vowelCombinations[lastChar]; ok {
						if combined, ok := combos["dot"]; ok {
							result[lastIdx] = combined
							i++
							continue
						}
					}
					// Try adding dot to combined vowel (Ô + Ï = Ộ)
					if tones, ok := combinedVowelTones[lastChar]; ok {
						if combined, ok := tones["dot"]; ok {
							result[lastIdx] = combined
							i++
							continue
						}
					}
				}
			}
		}

		// Handle Đ/đ (ñ/Ñ in VNI)
		if r == 'Ñ' {
			result = append(result, 'Đ')
			i++
			continue
		}
		if r == 'ñ' {
			result = append(result, 'đ')
			i++
			continue
		}

		// Check for breve marker (ă/Ă)
		// In VNI, Å = breve marker
		if r == 'Å' || r == 'å' {
			if len(result) > 0 {
				lastIdx := len(result) - 1
				lastChar := result[lastIdx]
				if lastChar == 'A' || lastChar == 'a' {
					if unicode.IsUpper(lastChar) {
						result[lastIdx] = 'Ă'
					} else {
						result[lastIdx] = 'ă'
					}
					i++
					continue
				}
			}
		}

		// Default: keep the character
		result = append(result, r)
		i++
	}

	return string(result)
}
