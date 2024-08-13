package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"integration-git/main/cmd"
	"integration-git/main/pkg/common/config"
	moduleresultcommon "integration-git/main/pkg/result/common"
	modulescanadapter "integration-git/main/pkg/scan/adapter"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func main() {

	//Read config
	_, err := config.LoadConfig("config.json")
	server := cmd.NewServer()

	gitFiles, err := server.GitModule.Controller.GetFilesToBeCommited()
	check(err)

	root, _ := os.Getwd()
	var paths []string
	for _, gitFile := range gitFiles {
		paths = append(paths, gitFile.Path)
	}

	err = server.ScanModule.Controller.RunScan(
		modulescanadapter.ScanInputPathsDTO{
			Paths:    paths,
			BasePath: root,
		})
	check(err)

	results, err := server.ResultModule.Controller.GetAll(nil)
	check(err)

	totalMatches := 0
	for _, result := range results {
		if result.MatchType == moduleresultcommon.File || result.MatchType == moduleresultcommon.Snippet {
			totalMatches++
		}
	}
	fmt.Printf("Total matches: %d\n", totalMatches)
	if totalMatches == 0 {
		os.Exit(0)
	}

	// Create an instance of the app structure
	app := NewApp(server)

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
