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
	"github.com/scanoss/scanoss.lui/backend/main/entities"
	"github.com/scanoss/scanoss.lui/backend/main/mappers"
	"github.com/scanoss/scanoss.lui/backend/main/repository"
	"github.com/scanoss/scanoss.lui/backend/main/service"
	"github.com/scanoss/scanoss.lui/backend/main/utils"

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

	fmt.Println("App Version: ", entities.AppVersion)

	app := NewApp()

	fr := utils.NewDefaultFileReader()

	// Repositories
	componentRepository := repository.NewJSONComponentRepository(fr)
	scanossSettingsRepository := repository.NewScanossSettingsJsonRepository(fr)
	scanossSettingsRepository.Init()
	resultRepository := repository.NewResultRepositoryJsonImpl(fr)
	fileRepository := repository.NewFileRepositoryImpl()

	// Mappers
	resultMapper := mappers.NewResultMapper(entities.ScanossSettingsJson)
	componentMapper := mappers.NewComponentMapper()

	// Services
	componentService := service.NewComponentServiceImpl(componentRepository, scanossSettingsRepository, resultRepository, componentMapper)
	fileService := service.NewFileService(fileRepository, componentRepository)
	keyboardService := service.NewKeyboardServiceInMemoryImpl()
	resultService := service.NewResultServiceImpl(resultRepository, resultMapper)
	scanossSettingsService := service.NewScanossSettingsServiceImpl(scanossSettingsRepository)

	//Create application with options
	err := wails.Run(&options.App{
		Title: "Scanoss Lui",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		WindowStartState: options.Maximised,
		OnStartup: func(ctx context.Context) {
			app.Init(ctx, scanossSettingsService, keyboardService)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			return app.BeforeClose(ctx)
		},
		Bind: []interface{}{
			app,
			componentService,
			fileService,
			keyboardService,
			resultService,
			scanossSettingsService,
		},
		EnumBind: []interface{}{
			entities.AllShortcutActions,
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
				Message: "Version: " + entities.AppVersion,
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
