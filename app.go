// Package main handles Wails runtime binding and application logic.
package main

import (
	"context"
	"convert-vni-to-unicode/internal/engine"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Config holds the processing configuration from Frontend
// Why: Standard DTO for passing parameters.
type Config struct {
	InputPath string `json:"inputPath"`
	SheetName string `json:"sheetName"` // Optional
}

// ProcessResult holds the result to send back to Frontend
type ProcessResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	OutputPath string `json:"outputPath"`
}

// SelectFile opens a file dialog to select the Excel file
// Why: Native dialog for better UX.
func (a *App) SelectFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Excel File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Excel Files", Pattern: "*.xlsx"},
		},
	})
}

// Process runs the conversion
// Why: Main entry point for the frontend to trigger logic.
func (a *App) Process(cfg Config) ProcessResult {
	if cfg.InputPath == "" {
		return ProcessResult{Success: false, Message: "Please select an input file"}
	}

	// Create processor
	p := engine.NewProcessor(cfg.InputPath, cfg.SheetName)

	// Setup progress tracing
	progressChan := make(chan float64, 100)
	p.SetProgressChan(progressChan)

	// Stream progress to frontend
	go func() {
		for prog := range progressChan {
			runtime.EventsEmit(a.ctx, "progress", prog)
		}
	}()

	// Run conversion
	// Note: Run blocks until completion.
	outputPath, err := p.Run(a.ctx)
	if err != nil {
		return ProcessResult{Success: false, Message: err.Error()}
	}

	return ProcessResult{
		Success:    true,
		Message:    "Conversion completed successfully!",
		OutputPath: outputPath,
	}
}

// ShowInFolder opens the file explorer and selects the file.
// Why: Native Windows integration for better UX.
func (a *App) ShowInFolder(path string) {
	if path == "" {
		return
	}
	// Use Windows-native "explorer /select" to open folder and highlight file
	// Using CommandContext to suppress noctx linter, though context cancellation isn't strictly needed
	// for fire-and-forget.
	cmd := exec.CommandContext(a.ctx, "explorer", "/select,", path)
	_ = cmd.Start() // Fire and forget, error is non-critical
}
