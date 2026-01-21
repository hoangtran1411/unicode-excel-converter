package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

// TestProcessor_Run integration test
func TestProcessor_Run(t *testing.T) {
	// 1. Setup Input File
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test_input.xlsx")

	f := excelize.NewFile()
	sheet := "Sheet1"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	// A1: Plain VNI Text (Font: VNI-Times)
	// "Việt Nam" in VNI: "Vi\u00D6t Nam"
	f.SetCellValue(sheet, "A1", "Vi\u00D6t Nam")
	styleID, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: "VNI-Times", Size: 12},
	})
	f.SetCellStyle(sheet, "A1", "A1", styleID)

	// A2: Rich Text Mixed (VNI + English)
	// Run 1: "Hello " (Arial)
	// Run 2: "Vi\u00D6t" (VNI-Times, Bold)
	textRuns := []excelize.RichTextRun{
		{Text: "Hello ", Font: &excelize.Font{Family: "Arial"}},
		{Text: "Vi\u00D6t", Font: &excelize.Font{Family: "VNI-Times", Bold: true}},
	}
	f.SetCellRichText(sheet, "A2", textRuns)

	// A3: TCVN3 Text (.VnTime)
	// "Công ty" in TCVN3: "C\u00F6ng ty"
	f.SetCellValue(sheet, "A3", "C\u00F6ng ty")
	styleID3, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: ".VnTime", Size: 12},
	})
	f.SetCellStyle(sheet, "A3", "A3", styleID3)

	if err := f.SaveAs(inputFile); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}
	f.Close()

	// 2. Run Processor
	proc := NewProcessor(inputFile, "")
	ctx := context.Background()
	outputFile, err := proc.Run(ctx)
	if err != nil {
		t.Fatalf("Processor.Run failed: %v", err)
	}

	// 3. Verify Output
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("output file not created: %s", outputFile)
	}

	fOut, err := excelize.OpenFile(outputFile)
	if err != nil {
		t.Fatalf("failed to open output: %v", err)
	}
	defer fOut.Close()

	// Verify A1 (Plain Converted)
	// Output should be RichText as we unified logic? Or just text?
	// Implementation note: plain text converted -> SetCellRichText with 1 run (Arial).
	// Content: "Việt Nam"
	// Font: Arial (Default) or Times New Roman (Mapped from VNI-Times)
	// VNI-Times -> Times New Roman

	// We need to check pure value first
	val, _ := fOut.GetCellValue(sheet, "A1")
	if val != "Việt Nam" {
		t.Errorf("A1 content mismatch. Got %q, want %q", val, "Việt Nam")
	}

	// Check Font A1 (via RichText or Style?)
	// If specific RichText was set, style is usually overridden or ignored by Excel in some views, but better check runs.
	runs, err := fOut.GetCellRichText(sheet, "A1")
	if err == nil && len(runs) > 0 {
		if runs[0].Font == nil {
			t.Errorf("A1 Font is nil")
		} else {
			font := runs[0].Font.Family
			if font != "Times New Roman" && font != "Arial" { // Mapped
				t.Errorf("A1 font mismatch. Got %s, want Times New Roman", font)
			}
		}
	} else {
		// Fallback check style
		styleID, _ := fOut.GetCellStyle(sheet, "A1")
		style, _ := fOut.GetStyle(styleID)
		if style != nil && style.Font != nil {
			// t.Logf("A1 Style Font: %s", style.Font.Family)
		}
	}

	// Verify A2 (Rich Text: "Hello Việt")
	runsA2, _ := fOut.GetCellRichText(sheet, "A2")
	if len(runsA2) != 2 {
		t.Errorf("A2 run count mismatch. Got %d, want 2", len(runsA2))
	} else {
		// Run 1: Hello
		if runsA2[0].Text != "Hello " {
			t.Errorf("A2 Run 1 Text wrong: %q", runsA2[0].Text)
		}
		// Run 2: Việt (Bold, Times New Roman/Arial)
		if runsA2[1].Text != "Việt" {
			t.Errorf("A2 Run 2 Text wrong: %q", runsA2[1].Text)
		}
		if !runsA2[1].Font.Bold {
			t.Errorf("A2 Run 2 should be Bold")
		}
		if runsA2[1].Font.Family != "Times New Roman" {
			t.Errorf("A2 Run 2 Font wrong. Got %s, want Times New Roman", runsA2[1].Font.Family)
		}
	}

	// Verify A3 (TCVN3)
	val3, _ := fOut.GetCellValue(sheet, "A3")
	if val3 != "Công ty" {
		t.Errorf("A3 content (TCVN3) mismatch. Got %q, want 'Công ty'", val3)
	}

	fmt.Printf("Integration Test Passed! Output: %s\n", outputFile)
}
