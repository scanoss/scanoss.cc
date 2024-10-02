package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/scanoss/scanoss.lui/backend/main/cmd"

	"github.com/scanoss/scanoss.lui/backend/handlers"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// Gets updated build time using -ldflags
var version = ""

//go:embed build/assets/icon.gif
var icon []byte

func main() {

	cmd.Init()

	fmt.Println("App Version: ", version)

	app := NewApp()

	fileHandler := handlers.NewFileHandler()
	scanossBomHandler := handlers.NewScanossSettingsHandler()
	resultHandler := handlers.NewResultHandler()
	componentHandler := handlers.NewComponentHandler()

	//Create application with options
	err := wails.Run(&options.App{
		Title: "Scanoss Lui",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		WindowStartState: options.Maximised,
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			return app.beforeClose(ctx, scanossBomHandler)
		},
		Bind: []interface{}{
			app,
			fileHandler,
			resultHandler,
			componentHandler,
			scanossBomHandler,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: true,
			Icon:                icon,
			ProgramName:         "Scanoss Lui",
		},
		Mac: &mac.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Scanoss Lightweight User Interface",
				Message: "Version: " + version,
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
