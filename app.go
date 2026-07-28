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
	"path/filepath"
	goRuntime "runtime"
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
	a.restoreWindowGeometry()
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

// restoreWindowGeometry restores the window size/position/maximized state from the
// previous session, if any, clamped to fit the screen the app is currently opening on.
// If nothing was saved, the window keeps the WindowStartState configured in main.go.
func (a *App) restoreWindowGeometry() {
	bounds, hasSaved := a.cfg.GetWindowBounds()
	if !hasSaved {
		return
	}

	screens, err := runtime.ScreenGetAll(a.ctx)
	if err != nil || len(screens) == 0 {
		log.Warn().Err(err).Msg("unable to determine screen bounds; skipping window geometry restore")
		return
	}
	screen := currentScreen(screens)

	width, height := bounds.Width, bounds.Height
	if width > screen.Size.Width {
		width = screen.Size.Width
	}
	if height > screen.Size.Height {
		height = screen.Size.Height
	}

	runtime.WindowUnmaximise(a.ctx)
	runtime.WindowSetSize(a.ctx, width, height)

	if fitsOnScreen(bounds.X, bounds.Y, width, height, screen) {
		runtime.WindowSetPosition(a.ctx, bounds.X, bounds.Y)
	} else {
		runtime.WindowCenter(a.ctx)
	}

	if bounds.Fullscreen {
		runtime.WindowFullscreen(a.ctx)
	} else if bounds.Maximized {
		runtime.WindowMaximise(a.ctx)
	}
}

func currentScreen(screens []runtime.Screen) runtime.Screen {
	for _, s := range screens {
		if s.IsCurrent {
			return s
		}
	}
	for _, s := range screens {
		if s.IsPrimary {
			return s
		}
	}
	return screens[0]
}

func fitsOnScreen(x, y, width, height int, screen runtime.Screen) bool {
	return x >= 0 && y >= 0 && x+width <= screen.Size.Width && y+height <= screen.Size.Height
}

// saveWindowGeometry persists the current window size/position/maximized state so it
// can be restored on the next launch.
func (a *App) saveWindowGeometry(ctx context.Context) {
	maximized := runtime.WindowIsMaximised(ctx)
	fullscreen := runtime.WindowIsFullscreen(ctx)

	bounds := config.WindowBounds{
		Maximized:  maximized,
		Fullscreen: fullscreen,
	}

	if !maximized && !fullscreen {
		width, height := runtime.WindowGetSize(ctx)
		x, y := runtime.WindowGetPosition(ctx)

		bounds.Width = width
		bounds.Height = height
		bounds.X = x
		bounds.Y = y
	} else {
		currentConfig, _ := a.cfg.GetWindowBounds()
		bounds.Width = currentConfig.Width
		bounds.Height = currentConfig.Height
		bounds.X = currentConfig.X
		bounds.Y = currentConfig.Y
	}

	if err := a.cfg.SetWindowBounds(bounds); err != nil {
		log.Error().Err(err).Msg("error saving window bounds")
	}
}

func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	a.saveWindowGeometry(ctx)

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

func (a *App) BuildMenu(keyboardService service.KeyboardService) *menu.Menu {
	AppMenu := menu.NewMenu()

	if goRuntime.GOOS == "darwin" {
		AppMenu.Append(menu.AppMenu())
	}

	groupedShortcuts := keyboardService.GetGroupedShortcuts()
	globalShortcuts := groupedShortcuts[entities.GroupGlobal]
	viewShortcuts := groupedShortcuts[entities.GroupView]

	// Add custom global shortcuts to File menu
	FileMenu := AppMenu.AddSubmenu("File")
	for _, shortcut := range globalShortcuts {
		if shortcut.Action == entities.ActionUndo ||
			shortcut.Action == entities.ActionRedo ||
			shortcut.Action == entities.ActionSelectAll {
			continue
		}
		sc := shortcut
		FileMenu.AddText(sc.Name, sc.Accelerator, func(cd *menu.CallbackData) {
			runtime.EventsEmit(a.ctx, string(sc.Action))
		})
	}

	// Actions menu with submenus
	ActionsMenu := AppMenu.AddSubmenu("Actions")

	// Include submenu
	IncludeMenu := ActionsMenu.AddSubmenu("Include")
	IncludeMenu.AddText("Include file", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionIncludeFile))
	})
	IncludeMenu.AddText("Include folder", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionIncludeFolder))
	})
	IncludeMenu.AddText("Include component", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionIncludeComponent))
	})

	// Dismiss submenu
	DismissMenu := ActionsMenu.AddSubmenu("Dismiss")
	DismissMenu.AddText("Dismiss file", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionDismissFile))
	})
	DismissMenu.AddText("Dismiss folder", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionDismissFolder))
	})
	DismissMenu.AddText("Dismiss component", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionDismissComponent))
	})

	// Replace submenu
	ReplaceMenu := ActionsMenu.AddSubmenu("Replace")
	ReplaceMenu.AddText("Replace file", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionReplaceFile))
	})
	ReplaceMenu.AddText("Replace folder", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionReplaceFolder))
	})
	ReplaceMenu.AddText("Replace component", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionReplaceComponent))
	})

	// Skip submenu
	SkipMenu := ActionsMenu.AddSubmenu("Skip")
	SkipMenu.AddText("Skip file", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionSkipFile))
	})
	SkipMenu.AddText("Skip folder", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionSkipFolder))
	})
	SkipMenu.AddText("Skip extension", nil, func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, string(entities.ActionSkipExtension))
	})

	// View menu
	ViewMenu := AppMenu.AddSubmenu("View")
	for _, shortcut := range viewShortcuts {
		sc := shortcut
		ViewMenu.AddText(sc.Name, sc.Accelerator, func(cd *menu.CallbackData) {
			runtime.EventsEmit(a.ctx, string(sc.Action))
		})
	}

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

	if goRuntime.GOOS == "darwin" {
		AppMenu.Append(menu.EditMenu())
	}

	return AppMenu
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

func (a *App) GetRecentScanRoots() ([]string, error) {
	return a.cfg.GetRecentScanRoots(), nil
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

func (a *App) JoinPaths(elements []string) string {
	return filepath.Join(elements...)
}
