package engine

import (
	"context"
	"convert-vni-to-unicode/internal/converter"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

// Job represents a single cell to be processed.
// Why: Standard unit of work for the worker pool.
type Job struct {
	SheetName string
	Axis      string
	Text      string
	RichText  []excelize.RichTextRun
	IsRich    bool
}

// Result represents the outcome of a job.
type Result struct {
	Job       Job
	Converted string
	NewRuns   []excelize.RichTextRun
	Error     error
}

// Processor manages the conversion process.
// Why: Central controller for the batch operation.
type Processor struct {
	InputPath string
	SheetName string
	// State
	f            *excelize.File
	jobs         chan Job
	results      chan Result
	progressChan chan float64
	processed    int

	// Format Preservers for different encodings
	vniPreserver   *FormatPreserver
	tcvn3Preserver *FormatPreserver
}

// NewProcessor creates a new processor instance.
func NewProcessor(inputPath, sheetName string) *Processor {
	return &Processor{
		InputPath:      inputPath,
		SheetName:      sheetName,
		jobs:           make(chan Job, 100),
		results:        make(chan Result, 100),
		vniPreserver:   NewFormatPreserver(converter.NewVNIConverter()),
		tcvn3Preserver: NewFormatPreserver(converter.NewTCVN3Converter()),
	}
}

// SetProgressChan sets the channel for progress updates.
func (p *Processor) SetProgressChan(ch chan float64) {
	p.progressChan = ch
}

// Run executes the conversion process.
func (p *Processor) Run(ctx context.Context) (string, error) {
	var err error
	p.f, err = excelize.OpenFile(p.InputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open excel: %w", err)
	}
	defer p.f.Close()

	// 1. Determine sheets to process
	sheets := p.f.GetSheetList()
	if p.SheetName != "" {
		// Validating if sheet exists
		found := false
		for _, s := range sheets {
			if s == p.SheetName {
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("sheet %s not found", p.SheetName)
		}
		sheets = []string{p.SheetName}
	}

	// 2. Count total cells (estimate) for progress
	// Note: Accurate count is hard without full scan. We will count rows.
	// For simplicity, we update progress based on processed count vs estimated.
	// Let's iterate first to dispatch jobs.

	// Start Workers
	var wg sync.WaitGroup
	workerCount := 10 // Default
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go p.worker(&wg)
	}

	// Dispatcher
	go func() {
		defer close(p.jobs)
		for _, sheet := range sheets {
			rows, err := p.f.Rows(sheet)
			if err != nil {
				continue
			}

			// We need column names (A, B, C...) to construct Axis (A1, B2...)
			// Rows iterator returns []string.
			// But to update specific cells including RichText, we need coordinates.
			// rows.Next() -> rows.Columns() returns values.
			// To get Axis, we track row index.

			rowIdx := 0
			for rows.Next() {
				rowIdx++
				cols, err := rows.Columns()
				if err != nil {
					fmt.Printf("Error getting columns for row %d: %v\n", rowIdx, err)
					continue
				}
				for colIdx, text := range cols {
					// 0-indexed colIdx -> "A", "B"
					axis, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx)
					if err != nil {
						fmt.Printf("Error converting coordinates for row %d col %d: %v\n", rowIdx, colIdx+1, err)
						continue
					}

					if strings.TrimSpace(text) == "" {
						continue
					}

					// Strategy: Unify everything to RichText for consistent processing.
					// 1. Try to get existing RichText
					runs, err := p.f.GetCellRichText(sheet, axis)
					isRich := false
					if err == nil && len(runs) > 0 {
						isRich = true
					}

					// 2. If no RichText, create synthetic RichText from Plain Text + Style Font
					if !isRich {
						fontName := ""
						styleID, err := p.f.GetCellStyle(sheet, axis)
						if err == nil {
							style, err := p.f.GetStyle(styleID)
							if err == nil && style.Font != nil {
								fontName = style.Font.Family
								fmt.Printf("DEBUG: Cell %s has Font: %s\n", axis, fontName)
							} else {
								fmt.Printf("DEBUG: Cell %s has NO Font (Style Error: %v)\n", axis, err)
							}
						} else {
							fmt.Printf("DEBUG: Cell %s GetCellStyle Error: %v\n", axis, err)
						}
						// Create synthetic run
						runs = []excelize.RichTextRun{
							{
								Text: text,
								Font: &excelize.Font{Family: fontName, Size: 11},
							},
						}
					}

					// Send Job
					p.jobs <- Job{
						SheetName: sheet,
						Axis:      axis,
						Text:      text, // Optional fallback
						RichText:  runs,
						IsRich:    isRich, // Track if it originated as Rich to optionally optimize write back? No, just always write Rich for consistency.
					}
				}
			}
			rows.Close()
		}
	}()

	// Collector (Writer)
	go func() {
		wg.Wait()
		close(p.results)
	}()

	p.processed = 0

	for res := range p.results {
		if res.Error != nil {
			fmt.Printf("Error processing %s: %v\n", res.Job.Axis, res.Error)
			continue
		}

		// Always write Rich Text to enforce font/format
		if err := p.f.SetCellRichText(res.Job.SheetName, res.Job.Axis, res.NewRuns); err != nil {
			fmt.Printf("Error writing rich text to %s: %v\n", res.Job.Axis, err)
		}

		p.processed++
		if p.progressChan != nil {
			p.progressChan <- float64(p.processed)
		}
	}

	// Save
	timestamp := time.Now().Format("2006_01_02_15_04_05") // yyyy_MM_dd_ss format as requested
	// User req: output_yyyy_MM_dd_ss
	// Actually: "sufix là output_yyyy_MM_dd_ss"
	// Example: contract.xlsx -> contract_output_2026_01_21_45.xlsx

	ext := filepath.Ext(p.InputPath)
	base := strings.TrimSuffix(p.InputPath, ext)
	outputPath := fmt.Sprintf("%s_output_%s%s", base, timestamp, ext)

	if err := p.f.SaveAs(outputPath); err != nil {
		return "", err
	}

	return outputPath, nil
}

func (p *Processor) worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range p.jobs {
		// Read Rich Text to get Fonts
		// Note: `f` is technically not thread-safe for Writes, but Reads?
		// excelize documentation says "File is not thread safe".
		// We CANNOT access `p.f` inside workers if they run in parallel.
		// MAJOR ARCHITECTURE FIX NEEDED:
		// We cannot read `p.f.GetCellRichText` inside workers concurrently.
		// We must read EVERYTHING in the Dispatcher (Single Thread) or use a Mutex.
		// Given we want speed, Mutex on `f` makes it serial.
		// Strategy:
		// Dispatcher (Serial) reads the Cell content (Text OR RichText) and creates the Job.
		// Workers (Parallel) process the String conversion (Pure CPU).
		// Collector (Serial) writes back.

		// Wait, `Scan` in Dispatcher:
		// `rows.Next()` gives Text. It does NOT give RichText.
		// We have to call `f.GetCellRichText` for every cell? That's slow.
		// BUT `rows` iterator only gives string values.
		// If we want format preservation, implementation using `rows` iterator is insufficient if we rely on it for content.
		// We MUST assume cells might be RichText.

		// Revised Flow:
		// Dispatcher:
		// Iterate rows. Get Axis.
		// Call `f.GetCellRichText`.
		// If runs > 0 or error == nil -> It's RichText or String.
		// Retrieve runs.
		// Create Job with runs.
		// Send to Worker.

		// Since `p.f` access must be serialized, Dispatcher does the heavy lifting of reading.
		// Worker does converting (CPU).
		// Writer does writing.

		// This puts load on Dispatcher.
		// Is `GetCellRichText` fast? It reads XML.
		// It's the only way to get font info per run.

		// Update logic in `Run` (Dispatcher part) to read RichText.
		// Here in Worker, we just process.

		res := Result{Job: job}

		// Detect encoding from Fonts in RichText Runs
		// Heuristic: Check first run's font or majority.
		// Or process run-by-run.

		// We use `vniPreserver` or `tcvn3Preserver`?
		// We need to AUTO detect which preserver to use for the cell.
		// If font is "VNI-Times" -> VNI.
		// If font is ".VnTime" -> TCVN3.

		// What if mixed? (Impossible usually).
		// We iterate runs and check font for EACH run.

		var newRuns []excelize.RichTextRun

		if len(job.RichText) > 0 {
			// Rich Text Handling
			// We need a generic Processor logic that can mix converters?
			// FormatPreserver IS the logic that iterates runs.
			// But FormatPreserver is tied to ONE converter.
			// We should make FormatPreserver smart or pass both?

			// Let's create a dynamic helper here.
			newRuns = make([]excelize.RichTextRun, 0, len(job.RichText))

			for _, run := range job.RichText {
				var text string
				fontName := ""
				if run.Font != nil {
					fontName = run.Font.Family
				}

				encoding := DetectEncoding(fontName, run.Text)

				// Apply conversion
				switch encoding {
				case converter.EncodingVNI:
					text = p.vniPreserver.converter.ToUnicode(run.Text)
					// Map Font
					if mapped, ok := FontMap[fontName]; ok {
						if run.Font == nil {
							run.Font = &excelize.Font{}
						}
						run.Font.Family = mapped
					} else {
						if run.Font == nil {
							run.Font = &excelize.Font{}
						}
						run.Font.Family = "Arial"
					}
				case converter.EncodingTCVN3:
					text = p.tcvn3Preserver.converter.ToUnicode(run.Text)
					if mapped, ok := FontMap[fontName]; ok {
						if run.Font == nil {
							run.Font = &excelize.Font{}
						}
						run.Font.Family = mapped
					} else {
						if run.Font == nil {
							run.Font = &excelize.Font{}
						}
						run.Font.Family = "Arial"
					}
				default:
					text = run.Text // No change
					// Even if no change, should we enforce Arial if it looks like garbage?
					// If unknown, leave it.
				}

				run.Text = text
				newRuns = append(newRuns, run)
			}
			res.NewRuns = newRuns
			res.Job.IsRich = true

		} else {
			// Plain text case (Dispatcher sent Text string, empty RichText)
			// But wait, if Dispatcher calls `GetCellRichText`, it gets runs even for plain text (usually 1 run with nil font?).
			// Check excelize behavior: "return error if no rich text".

			// So if Dispatcher failed to get RichText, it's a plain cell.
			// We have `job.Text`.
			// We don't know the font (it's in Cell Style).
			// We can try to detect by Content.

			// If we can't detect by Font, we detect by Content.
			// If detected -> Convert -> Force Arial.

			// Heuristic: Try VNI first, then TCVN3? Or check byte patterns.
			// TCVN3 uses specific bytes. VNI uses others.
			// For now, let's try to Detect by Content (not implemented yet, returns Unknown).

			// If Unknown, we return original.
			res.Converted = job.Text
			res.Job.IsRich = false
		}

		p.results <- res
	}
}
