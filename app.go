package main

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/scanoss/scanoss.lui/backend/entities"
	"github.com/scanoss/scanoss.lui/backend/service"
	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/scanoss/scanoss.lui/internal/utils"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx                    context.Context
	scanossSettingsService service.ScanossSettingsService
	keyboardService        service.KeyboardService
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(ctx context.Context, scanossSettingsService service.ScanossSettingsService, keyboardService service.KeyboardService) {
	a.ctx = ctx
	a.scanossSettingsService = scanossSettingsService
	a.keyboardService = keyboardService
	a.startup()
}

func (a *App) startup() {
	a.maybeSetWindowTitle()
	a.initializeMenu()
	runtime.LogInfo(a.ctx, fmt.Sprintf("Scan Settings file path: %s", config.GetInstance().ScanSettingsFilePath))
	runtime.LogInfo(a.ctx, fmt.Sprintf("Results file path: %s", config.GetInstance().ResultFilePath))
	runtime.LogInfo(a.ctx, fmt.Sprintf("Scan Root file path: %s", config.GetInstance().ScanRoot))
}

func (a *App) maybeSetWindowTitle() {
	if env := runtime.Environment(a.ctx); env.Platform != "darwin" {
		runtime.WindowSetTitle(a.ctx, fmt.Sprintf("Scanoss Lui %s", entities.AppVersion))
	}
}

func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	hasUnsavedChanges, err := a.scanossSettingsService.HasUnsavedChanges()

	if err != nil {
		runtime.LogError(ctx, "Error checking for unsaved changes: "+err.Error())
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
		runtime.LogError(ctx, "Error showing dialog: "+err.Error())
	}

	confirmOptions := []string{"Yes", "Ok"}

	if slices.Contains(confirmOptions, result) {
		err := a.scanossSettingsService.Save()
		if err != nil {
			runtime.LogError(ctx, "Error saving scanoss bom file: "+err.Error())
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
	ScanMenu.AddText("Scan Current Directory", keys.Combo("b", keys.ShiftKey, keys.CmdOrCtrlKey), func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "scanCurrentDirectory")
	})
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
	directory, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return "", fmt.Errorf("error selecting directory: %v", err)
	}

	if directory == "" {
		return "", nil
	}

	relativePath, err := utils.GetRelativePath(directory)
	if err != nil {
		return "", fmt.Errorf("error selecting directory: %v", err)
	}
	return relativePath, nil
}

func (a *App) GetWorkingDir() string {
	workingDir, _ := os.Getwd()
	return workingDir
}

func (a *App) SetScanRoot(path string) {
	config.GetInstance().ScanRoot = path
}
