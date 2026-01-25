// Package engine handles the Excel processing logic.
package engine

import (
	"context"
	"convert-vni-to-unicode/internal/converter"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

// Constants for processor configuration
const (
	// DefaultWorkerCount is the number of concurrent workers for cell processing.
	// This is CPU-bound work, so we use a reasonable fixed value.
	DefaultWorkerCount = 10

	// JobChannelBuffer is the buffer size for job and result channels.
	JobChannelBuffer = 100
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
// Thread-safety: The `f` (*excelize.File) field is NOT thread-safe.
// Only the dispatcher goroutine should read from `f`, and only the
// collector goroutine should write to `f`. Workers only process data.
type Processor struct {
	InputPath string
	SheetName string
	// State - NOT thread-safe, access must be serialized
	f            *excelize.File
	jobs         chan Job
	results      chan Result
	progressChan chan float64
	processed    int

	// Format Preservers for different encodings (thread-safe for reads)
	vniPreserver   *FormatPreserver
	tcvn3Preserver *FormatPreserver
}

// NewProcessor creates a new processor instance.
func NewProcessor(inputPath, sheetName string) *Processor {
	return &Processor{
		InputPath:      inputPath,
		SheetName:      sheetName,
		jobs:           make(chan Job, JobChannelBuffer),
		results:        make(chan Result, JobChannelBuffer),
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
	defer func() {
		if closeErr := p.f.Close(); closeErr != nil {
			slog.Error("failed to close excel file", "error", closeErr)
		}
	}()

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
			return "", fmt.Errorf("sheet %q not found", p.SheetName)
		}
		sheets = []string{p.SheetName}
	}

	// Start Workers
	var wg sync.WaitGroup
	for i := 0; i < DefaultWorkerCount; i++ {
		wg.Add(1)
		go p.worker(&wg)
	}

	// Dispatcher - runs in a separate goroutine
	go p.processSheets(ctx, sheets)

	// Collector (Writer) - waits for workers to finish, then closes results
	go func() {
		wg.Wait()
		close(p.results)
	}()

	p.processed = 0

	for res := range p.results {
		if res.Error != nil {
			slog.Error("failed to process cell", "cell", res.Job.Axis, "error", res.Error)
			continue
		}

		// Always write Rich Text to enforce font/format
		if err := p.f.SetCellRichText(res.Job.SheetName, res.Job.Axis, res.NewRuns); err != nil {
			slog.Error("failed to write rich text", "cell", res.Job.Axis, "error", err)
		}

		p.processed++
		if p.progressChan != nil {
			p.progressChan <- float64(p.processed)
		}
	}

	// Save with timestamp suffix
	timestamp := time.Now().Format("2006_01_02_15_04_05")
	ext := filepath.Ext(p.InputPath)
	base := strings.TrimSuffix(p.InputPath, ext)
	outputPath := fmt.Sprintf("%s_output_%s%s", base, timestamp, ext)

	if err := p.f.SaveAs(outputPath); err != nil {
		return "", fmt.Errorf("failed to save output file: %w", err)
	}

	return outputPath, nil
}

// processSheets iterates through sheets to dispatch jobs
func (p *Processor) processSheets(ctx context.Context, sheets []string) {
	defer close(p.jobs)
	for _, sheet := range sheets {
		p.processSheet(ctx, sheet)
	}
}

func (p *Processor) processSheet(ctx context.Context, sheet string) {
	rows, err := p.f.Rows(sheet)
	if err != nil {
		slog.Error("failed to get rows", "sheet", sheet, "error", err)
		return
	}

	rowIdx := 0
	for rows.Next() {
		rowIdx++
		cols, err := rows.Columns()
		if err != nil {
			slog.Error("failed to get columns", "sheet", sheet, "row", rowIdx, "error", err)
			continue
		}
		for colIdx, text := range cols {
			// Check for cancellation
			select {
			case <-ctx.Done():
				return
			default:
			}

			axis, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx)
			if err != nil {
				slog.Error("failed to convert coordinates", "row", rowIdx, "col", colIdx+1, "error", err)
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
						slog.Debug("cell font detected", "cell", axis, "font", fontName)
					}
				}
				// Create synthetic run with capacity hint
				runs = make([]excelize.RichTextRun, 0, 1)
				runs = append(runs, excelize.RichTextRun{
					Text: text,
					Font: &excelize.Font{Family: fontName, Size: 11},
				})
			}

			// Send Job
			p.jobs <- Job{
				SheetName: sheet,
				Axis:      axis,
				Text:      text,
				RichText:  runs,
				IsRich:    isRich,
			}
		}
	}
	if err := rows.Close(); err != nil {
		slog.Error("failed to close rows iterator", "sheet", sheet, "error", err)
	}
}

func (p *Processor) worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range p.jobs {
		// Worker only processes data, does NOT access p.f (not thread-safe)
		res := Result{Job: job}

		// Pre-allocate with capacity hint
		newRuns := make([]excelize.RichTextRun, 0, len(job.RichText))

		if len(job.RichText) > 0 {
			// Rich Text Handling - process each run independently
			for _, run := range job.RichText {
				var text string
				fontName := ""
				if run.Font != nil {
					fontName = run.Font.Family
				}

				encoding := DetectEncoding(fontName, run.Text)

				// Apply conversion based on detected encoding
				switch encoding {
				case converter.EncodingVNI:
					text = p.vniPreserver.converter.ToUnicode(run.Text)
					// Map Font to Unicode equivalent
					if mapped, ok := FontMap[fontName]; ok {
						if run.Font == nil {
							run.Font = &excelize.Font{}
						}
						run.Font.Family = mapped
					} else {
						if run.Font == nil {
							run.Font = &excelize.Font{}
						}
						run.Font.Family = DefaultFont
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
						run.Font.Family = DefaultFont
					}
				default:
					text = run.Text // No change for unknown encoding
				}

				run.Text = text
				newRuns = append(newRuns, run)
			}
			res.NewRuns = newRuns
			res.Job.IsRich = true

		} else {
			// Plain text fallback (should rarely happen with new dispatcher logic)
			res.Converted = job.Text
			res.Job.IsRich = false
		}

		p.results <- res
	}
}
