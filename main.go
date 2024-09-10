package main

import (
	"embed"
	"fmt"

	"github.com/scanoss/scanoss.lui/backend/main/cmd"

	"github.com/scanoss/scanoss.lui/backend/handlers"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// Get's updated build time using -ldflags
var version = ""

func main() {

	cmd.Init()

	fmt.Println("App Version: ", version)
	// Create an instance of the app structure

	app := NewApp()

	scanossBomHandler := handlers.NewScanossBomHandler()
	//Create application with options
	err := wails.Run(&options.App{
		Title:  "scanoss-lui",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			handlers.NewFileHandler(),
			handlers.NewResultHandler(),
			handlers.NewComponentHandler(),
			scanossBomHandler,
		},
		OnShutdown: scanossBomHandler.OnShutDown,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
