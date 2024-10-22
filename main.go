package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/scanoss/scanoss.lui/backend/main/cmd"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/version"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"

	"github.com/scanoss/scanoss.lui/backend/handlers"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/assets/icon.gif
var icon []byte

func main() {

	validate := validator.New()
	validate.RegisterValidation("valid-purl", utils.ValidatePurl)
	utils.SetValidator(validate)

	cmd.Init()

	fmt.Println("App Version: ", version.AppVersion)

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
			app.initializeMenu()
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
				Message: "Version: " + version.AppVersion,
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
