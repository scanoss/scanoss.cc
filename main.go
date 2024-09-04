package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"integration-git/handler"
	"integration-git/main/cmd"
)

//go:embed all:frontend/dist
var assets embed.FS

var version = "0.1.0"

func main() {

	fmt.Println("App Version: ", version)
	// Create an instance of the app structure
	app := NewApp()

	// Inputs provided will override the values in configuration if they differ.
	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	scanossBomHandler := handler.NewScanossBomHandler()
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
			handler.NewFileHandler(),
			handler.NewResultHandler(),
			handler.NewComponentHandler(),
			scanossBomHandler,
		},
		OnShutdown: scanossBomHandler.OnShutDown,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
