package converter

import (
	"testing"
)

func TestTCVN3Converter_ToUnicode(t *testing.T) {
	c := NewTCVN3Converter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Lowercase a with tones",
			input:    "\u00B8 \u00B5 \u00B6 \u00B7 \u00B9", // á à ả ã ạ
			expected: "á à ả ã ạ",
		},
		{
			name:     "TCVN3 Sample Word",
			input:    "C\u00F6ng ty", // "Cöng ty" in TCVN3 font displays as "Công ty"
			expected: "Công ty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.ToUnicode(tt.input)
			if got != tt.expected {
				t.Errorf("ToUnicode() = %q, want %q", got, tt.expected)
			}
		})
	}
}
