package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// main is the entry point of the application.
// Why: It initializes the Wails application, configures the window properties, and binds the backend logic (App) to the frontend.
func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	// Why: Defines the window dimensions, title, and theme to match the minimal aesthetic requested.
	err := wails.Run(&options.App{
		Title:         "VNI to Unicode Converter",
		Width:         900,
		Height:        835,
		DisableResize: false, // Allow resizing for better UX on different screens
		MinWidth:      700,
		MinHeight:     600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1}, // Matches the dark theme background
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
