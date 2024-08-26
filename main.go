package main

import (
	"embed"
	"integration-git/main/cmd"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	// Create an instance of the app structure
	app := NewApp()

	app.Init()

	// Inputs provided will override the values in configuration if they differ.
	if err := cmd.Execute(); err != nil {
		panic(err)
	}

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
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
