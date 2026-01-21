package engine

import (
	"convert-vni-to-unicode/internal/converter"

	"github.com/xuri/excelize/v2"
)

// FontMap maps legacy font names to Unicode standard fonts.
// Why: Provides a lookup table to automatically switch fonts after conversion.
var FontMap = map[string]string{
	// VNI Fonts
	"VNI-Times": "Times New Roman",
	"VNI-Arial": "Arial",
	"VNI-Helve": "Helvetica",
	"VNI-Hobo":  "Hobo Std", // Example
	// TCVN3 Fonts
	".VnTime":  "Times New Roman",
	".VnTimeH": "Times New Roman",
	".VnArial": "Arial",
	".VnHelve": "Helvetica",
}

// DefaultFont is the fallback font for converted text.
const DefaultFont = "Arial"

// FormatPreserver handles the preservation of styles while changing text.
// Why: Separates formatting logic from the main processor.
type FormatPreserver struct {
	converter converter.Converter
}

// NewFormatPreserver creates a new instance.
func NewFormatPreserver(c converter.Converter) *FormatPreserver {
	return &FormatPreserver{converter: c}
}

// ProcessRichText converts the text in runs and maps the fonts.
// Why: Rich Text allows multiple styles in one cell. We must iterate runs to preserve mixed styles.
func (fp *FormatPreserver) ProcessRichText(runs []excelize.RichTextRun) []excelize.RichTextRun {
	newRuns := make([]excelize.RichTextRun, len(runs))
	for i, run := range runs {
		// Convert text
		convertedText := fp.converter.ToUnicode(run.Text)

		// Create copy
		newRun := run
		newRun.Text = convertedText

		// Handle Font mapping
		if newRun.Font != nil {
			if mapping, ok := FontMap[newRun.Font.Family]; ok {
				newRun.Font.Family = mapping
			} else {
				// If no mapping found but text was converted, default to "Arial"
				// to ensure legible Unicode display if the old font doesn't support it.
				// However, user requested "Default font is Arial".
				// So we enforce Arial if it was a legacy font or if we act generically?
				// "font mặc định sau khi chuyển là Arial" -> Implies strictly Arial if not overridden?
				// Let's respect the user's specific request: Default = Arial.
				// But keeping original bold/italic is key.
				newRun.Font.Family = DefaultFont
			}
		} else {
			// If no font struct exists, create one with DefaultFont
			newRun.Font = &excelize.Font{
				Family: DefaultFont,
				Size:   11, // Default Excel size usually
			}
		}

		newRuns[i] = newRun
	}
	return newRuns
}

// GetConvertedFontFamily determines the new font family based on input.
func (fp *FormatPreserver) GetConvertedFontFamily(originalFont string) string {
	if mapped, ok := FontMap[originalFont]; ok {
		return mapped
	}
	return DefaultFont
}
