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
	"fmt"
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/service"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx                    context.Context
	scanossSettingsService service.ScanossSettingsService
	keyboardService        service.KeyboardService
	cfg                    *config.Config
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(ctx context.Context, scanossSettingsService service.ScanossSettingsService, keyboardService service.KeyboardService) {
	a.ctx = ctx
	a.scanossSettingsService = scanossSettingsService
	a.keyboardService = keyboardService
	a.cfg = config.GetInstance()
	a.startup()
}

func (a *App) startup() {
	a.maybeSetWindowTitle()
	a.initializeMenu()
	log.Debug().Msgf("Scan Settings file path: %s", a.cfg.GetScanSettingsFilePath())
	log.Debug().Msgf("Results file path: %s", a.cfg.GetResultFilePath())
	log.Debug().Msgf("Scan Root file path: %s", a.cfg.GetScanRoot())
	log.Info().Msgf("App Version: %s", entities.AppVersion)
}

func (a *App) maybeSetWindowTitle() {
	if env := runtime.Environment(a.ctx); env.Platform != "darwin" {
		runtime.WindowSetTitle(a.ctx, fmt.Sprintf("Scanoss Code Compare %s", entities.AppVersion))
	}
}

func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	hasUnsavedChanges, err := a.scanossSettingsService.HasUnsavedChanges()

	if err != nil {
		log.Error().Msg("Error checking for unsaved changes: " + err.Error())
		return false
	}
	if !hasUnsavedChanges {
		return false
	}

	result, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Unsaved Changes",
		Message:       "Do you want to save changes before closing the app?",
		CancelButton:  "No",
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "Yes",
	})
	if err != nil {
		log.Error().Msg("Error showing dialog: " + err.Error())
	}

	confirmOptions := []string{"Yes", "Ok"}

	if slices.Contains(confirmOptions, result) {
		err := a.scanossSettingsService.Save()
		if err != nil {
			log.Error().Msg("Error saving scanoss bom file: " + err.Error())
		}
	}

	return false
}

func (a *App) initializeMenu() {
	AppMenu := menu.NewMenu()

	if env := runtime.Environment(a.ctx); env.Platform == "darwin" {
		AppMenu.Append(menu.AppMenu())
	}

	groupedShortcuts := a.keyboardService.GetGroupedShortcuts()
	actionShortcuts := groupedShortcuts[entities.GroupActions]
	globalShortcuts := groupedShortcuts[entities.GroupGlobal]

	// Global edit menu
	EditMenu := AppMenu.AddSubmenu("Edit")
	for _, shortcut := range globalShortcuts {
		EditMenu.AddText(shortcut.Name, shortcut.Accelerator, func(cd *menu.CallbackData) {})
	}

	// Actions menu
	ActionsMenu := AppMenu.AddSubmenu("Actions")
	for _, shortcut := range actionShortcuts {
		ActionsMenu.AddText(shortcut.Name, shortcut.Accelerator, func(cd *menu.CallbackData) {})
	}

	// View menu
	ViewMenu := AppMenu.AddSubmenu("View")
	ViewMenu.AddText("Sync Scroll Position", keys.Combo("e", keys.ShiftKey, keys.CmdOrCtrlKey), func(cd *menu.CallbackData) {})

	// Scan menu
	ScanMenu := AppMenu.AddSubmenu("Scan")
	ScanMenu.AddText("Scan With Options", keys.Combo("c", keys.ShiftKey, keys.CmdOrCtrlKey), func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "scanWithOptions")
	})

	// Help menu
	HelpMenu := AppMenu.AddSubmenu("Help")
	HelpMenu.AddText("Report Issue", nil, func(cd *menu.CallbackData) {
		utils.OpenMailClient(utils.SCANOSS_SUPPORT_MAILBOX, "Report an issue", utils.GetIssueReportBody(a.ctx))
	})
	HelpMenu.AddText("Keyboard Shortcuts", keys.Combo("k", keys.ShiftKey, keys.CmdOrCtrlKey), func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionShowKeyboardShortcutsModal))
	})

	runtime.MenuSetApplicationMenu(a.ctx, AppMenu)
}

func (a *App) SelectDirectory() (string, error) {
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Select Directory",
		CanCreateDirectories: true,
	})
	if err != nil {
		return "", err
	}

	if dirPath == "" {
		return "", nil
	}

	return dirPath, nil
}

func (a *App) SelectFile(defaultDir string) (string, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Select File",
		DefaultDirectory:     defaultDir,
		ShowHiddenFiles:      true,
		CanCreateDirectories: true,
	})
	if err != nil {
		log.Error().Err(err).Msgf("error selecting file %v", err.Error())
		return "", err
	}

	if filePath == "" {
		return "", nil
	}

	return filePath, nil
}

func (a *App) GetScanRoot() (string, error) {
	return a.cfg.GetScanRoot(), nil
}

func (a *App) GetResultFilePath() (string, error) {
	return a.cfg.GetResultFilePath(), nil
}

func (a *App) GetScanSettingsFilePath() (string, error) {
	return a.cfg.GetScanSettingsFilePath(), nil
}

func (a *App) SetScanRoot(path string) {
	a.cfg.SetScanRoot(path)
}

func (a *App) SetResultFilePath(path string) {
	a.cfg.SetResultFilePath(path)
}

func (a *App) SetScanSettingsFilePath(path string) {
	a.cfg.SetScanSettingsFilePath(path)
}
