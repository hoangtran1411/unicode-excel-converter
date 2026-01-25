// Package main generates sample Excel files for testing.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("warning: failed to close file: %v", err)
		}
	}()

	sheet := "Sheet1"
	if _, err := f.NewSheet(sheet); err != nil {
		log.Fatal(err)
	}

	// Set Headers
	headers := []string{"VNI Content", "TCVN3 Content", "Mixed/Plain"}
	for i, h := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			log.Fatal(err)
		}
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			log.Fatal(err)
		}
		style, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
		if err != nil {
			log.Fatal(err)
		}
		if err := f.SetCellStyle(sheet, cell, cell, style); err != nil {
			log.Fatal(err)
		}
	}

	// 1. VNI Cell
	// "Việt Nam" -> "Vi\u00D6t Nam"
	if err := f.SetCellValue(sheet, "A2", "Vi\u00D6t Nam"); err != nil {
		log.Fatal(err)
	}
	styleVNI, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: "VNI-Times", Size: 12},
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := f.SetCellStyle(sheet, "A2", "A2", styleVNI); err != nil {
		log.Fatal(err)
	}

	// 2. TCVN3 Cell
	// "Công ty" -> "C\u00F6ng ty"
	if err := f.SetCellValue(sheet, "B2", "C\u00F6ng ty"); err != nil {
		log.Fatal(err)
	}
	styleTCVN3, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: ".VnTime", Size: 12},
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := f.SetCellStyle(sheet, "B2", "B2", styleTCVN3); err != nil {
		log.Fatal(err)
	}

	// 3. Rich Text (Mixed)
	runs := []excelize.RichTextRun{
		{Text: "Sample ", Font: &excelize.Font{Family: "Arial"}},
		{Text: "Vi\u00D6t", Font: &excelize.Font{Family: "VNI-Times", Bold: true, Color: "#FF0000"}},
	}
	if err := f.SetCellRichText(sheet, "C2", runs); err != nil {
		log.Fatal(err)
	}

	// Save
	if err := os.MkdirAll("samples", 0750); err != nil {
		log.Fatal(err)
	}
	output := "samples/sample_data.xlsx"
	if err := f.SaveAs(output); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sample file generated at: %s\n", output)
}
