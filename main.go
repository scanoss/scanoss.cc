package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"github.com/scanoss/scanoss.lui/backend/entities"
	"github.com/scanoss/scanoss.lui/backend/mappers"
	"github.com/scanoss/scanoss.lui/backend/repository"
	"github.com/scanoss/scanoss.lui/backend/service"
	"github.com/scanoss/scanoss.lui/cmd"
	"github.com/scanoss/scanoss.lui/internal/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/assets/appicon.png
var icon []byte

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	validate := validator.New()
	validate.RegisterValidation("valid-purl", utils.ValidatePurl)
	utils.SetValidator(validate)

	err := cmd.Execute()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	app := NewApp()

	fr := utils.NewDefaultFileReader()

	// Repositories
	componentRepository := repository.NewJSONComponentRepository(fr)
	scanossSettingsRepository := repository.NewScanossSettingsJsonRepository(fr)
	scanossSettingsRepository.Init()
	resultRepository := repository.NewResultRepositoryJsonImpl(fr)
	fileRepository := repository.NewFileRepositoryImpl()
	licenseRepository := repository.NewLicenseJsonRepository(fr)

	// Mappers
	resultMapper := mappers.NewResultMapper(entities.ScanossSettingsJson)
	componentMapper := mappers.NewComponentMapper()

	// Services
	componentService := service.NewComponentServiceImpl(componentRepository, scanossSettingsRepository, resultRepository, componentMapper)
	fileService := service.NewFileService(fileRepository, componentRepository)
	keyboardService := service.NewKeyboardServiceInMemoryImpl()
	resultService := service.NewResultServiceImpl(resultRepository, resultMapper)
	scanossSettingsService := service.NewScanossSettingsServiceImpl(scanossSettingsRepository)
	licenseService := service.NewLicenseServiceImpl(licenseRepository)
	scanService := service.NewScanServicePythonImpl()

	//Create application with options
	err = wails.Run(&options.App{
		Title: "Scanoss Lui",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		WindowStartState: options.Maximised,
		OnStartup: func(ctx context.Context) {
			app.Init(ctx, scanossSettingsService, keyboardService)
			scanService.SetContext(ctx)
			resultService.SetContext(ctx)

			// Add watchers
			go resultService.WatchResults()
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
			licenseService,
			scanService,
		},
		EnumBind: []interface{}{
			entities.AllShortcutActions,
		},
		Linux: &linux.Options{
			Icon:        icon,
			ProgramName: "Scanoss Lui",
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "Scanoss Lightweight User Interface",
				Message: "Version: " + entities.AppVersion,
				Icon:    icon,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	return nil
}
