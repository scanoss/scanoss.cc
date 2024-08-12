package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"integration-git/main/pkg/common/config"
	module_scan_adapter "integration-git/main/pkg/scan/adapter"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	//Read config
	_, err := config.LoadConfig("config.json")

	// Create an instance of the app structure
	app := NewApp()

	gitFiles, err := app.gitModule.Controller.GetFilesToBeCommited()
	if err != nil {
		panic(err)
	}

	root, _ := os.Getwd()
	var paths []string
	for _, gitFile := range gitFiles {
		paths = append(paths, gitFile.Path)
	}

	err = app.scanModule.Controller.RunScan(
		module_scan_adapter.ScanInputPathsDTO{
			Paths:    paths,
			BasePath: root,
		})

	if err != nil {
		panic(err)
	}

	//Create application with options
	err = wails.Run(&options.App{
		Title:  "integration-git",
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
