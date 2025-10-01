// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/mappers"
	"github.com/scanoss/scanoss.cc/backend/repository"
	"github.com/scanoss/scanoss.cc/backend/service"
	"github.com/scanoss/scanoss.cc/cmd"
	"github.com/scanoss/scanoss.cc/internal/utils"

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
	scanossSettingsRepository := repository.NewScanossSettingsJsonRepository(fr)
	scanossSettingsRepository.Init()
	resultRepository, err := repository.NewResultRepositoryJsonImplWithWatcher(fr)
	if err != nil {
		return fmt.Errorf("error initializing results repository")
	}
	defer resultRepository.Close()
	componentRepository := repository.NewJSONComponentRepository(fr, resultRepository)
	fileRepository := repository.NewFileRepositoryImpl()
	licenseRepository := repository.NewLicenseJsonRepository(fr)

	// Mappers
	resultMapper := mappers.NewResultMapper(entities.ScanossSettingsJson)
	componentMapper := mappers.NewComponentMapper()

	// Services
	scanossApiService, err := service.NewScanossApiServiceHttpImpl()
	if err != nil {
		return fmt.Errorf("error initializing scanoss api service: %v", err)
	}
	componentService := service.NewComponentServiceImpl(componentRepository, scanossSettingsRepository, resultRepository, scanossApiService, componentMapper)
	fileService := service.NewFileService(fileRepository, componentRepository)
	keyboardService := service.NewKeyboardServiceInMemoryImpl()
	resultService := service.NewResultServiceImpl(resultRepository, resultMapper)
	scanossSettingsService := service.NewScanossSettingsServiceImpl(scanossSettingsRepository)
	licenseService := service.NewLicenseServiceImpl(licenseRepository, scanossApiService)
	scanService := service.NewScanServicePythonImpl()
	treeService := service.NewTreeServiceImpl(resultService, scanossSettingsRepository)

	//Create application with options
	err = wails.Run(&options.App{
		Title: "Scanoss Code Compare",
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		WindowStartState: options.Maximised,
		OnStartup: func(ctx context.Context) {
			app.Init(ctx, scanossSettingsService, keyboardService)
			scanService.SetContext(ctx)
			resultService.SetContext(ctx)
			scanossApiService.SetContext(ctx)
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
			treeService,
		},
		EnumBind: []interface{}{
			entities.AllShortcutActions,
		},
		Linux: &linux.Options{
			Icon:        icon,
			ProgramName: "Scanoss Code Compare",
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "Scanoss Code Compare",
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
