package converter

import (
	"testing"
)

func TestVNIConverter_ToUnicode(t *testing.T) {
	c := NewVNIConverter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Lowercase a with tones",
			input:    "\u00E1 \u00E0 \u001E \u00E3 \u00B5", // á à ả ã ạ (VNI codes)
			expected: "á à ả ã ạ",
		},
		{
			name:     "Lowercase a circumflex with tones",
			input:    "\u00E2 \u00CA \u00C7 \u00C8 \u00C9 \u00CB", // â ấ ầ ẩ ẫ ậ
			expected: "â ấ ầ ẩ ẫ ậ",
		},
		{
			name: "Mixed sentence",
			// "Việt Nam" in VNI: V i \u00D6 t N a m
			input:    "Vi\u00D6t Nam",
			expected: "Việt Nam",
		},

		{
			name:     "Plain text",
			input:    "Hello World",
			expected: "Hello World",
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
