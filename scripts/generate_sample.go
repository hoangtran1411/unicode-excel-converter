package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Sheet1"
	f.NewSheet(sheet)

	// Set Headers
	headers := []string{"VNI Content", "TCVN3 Content", "Mixed/Plain"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
		f.SetCellStyle(sheet, cell, cell, style)
	}

	// 1. VNI Cell
	// "Việt Nam" -> "Vi\u00D6t Nam"
	f.SetCellValue(sheet, "A2", "Vi\u00D6t Nam")
	styleVNI, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: "VNI-Times", Size: 12},
	})
	f.SetCellStyle(sheet, "A2", "A2", styleVNI)

	// 2. TCVN3 Cell
	// "Công ty" -> "C\u00F6ng ty"
	f.SetCellValue(sheet, "B2", "C\u00F6ng ty")
	styleTCVN3, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: ".VnTime", Size: 12},
	})
	f.SetCellStyle(sheet, "B2", "B2", styleTCVN3)

	// 3. Rich Text (Mixed)
	runs := []excelize.RichTextRun{
		{Text: "Sample ", Font: &excelize.Font{Family: "Arial"}},
		{Text: "Vi\u00D6t", Font: &excelize.Font{Family: "VNI-Times", Bold: true, Color: "#FF0000"}},
	}
	f.SetCellRichText(sheet, "C2", runs)

	// Save
	if err := os.MkdirAll("samples", 0755); err != nil {
		log.Fatal(err)
	}
	output := "samples/sample_data.xlsx"
	if err := f.SaveAs(output); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sample file generated at: %s\n", output)
}
