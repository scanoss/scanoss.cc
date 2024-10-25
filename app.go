package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/scanoss/scanoss.lui/backend/handlers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/keyboard"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/version"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Only set window title on non-darwin platforms
func (a *App) maybeSetWindowTitle() {
	if env := runtime.Environment(a.ctx); env.Platform != "darwin" {
		runtime.WindowSetTitle(a.ctx, fmt.Sprintf("Scanoss Lui %s", version.AppVersion))
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.maybeSetWindowTitle()
	runtime.LogInfo(ctx, "Config file path: "+config.GetInstance().GetConfigPath())
	runtime.LogInfo(ctx, "Results file path: "+config.Get().ResultFilePath)
	runtime.LogInfo(ctx, "Scan Root file path: "+config.Get().ScanRoot)
}

func (a *App) beforeClose(ctx context.Context, sbh *handlers.ScanossSettingsHandler) (prevent bool) {

	hasUnsavedChanges, err := sbh.HasUnsavedChanges()

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
		err := sbh.SaveScanossSettingsFile()
		if err != nil {
			runtime.LogError(ctx, "Error saving scanoss bom file: "+err.Error())
		}
	}

	return false
}

func (a *App) initializeMenu(groupedShortcuts map[keyboard.Group][]keyboard.Shortcut) {
	AppMenu := menu.NewMenu()

	if env := runtime.Environment(a.ctx); env.Platform == "darwin" {
		AppMenu.Append(menu.AppMenu())
	}

	actionShortcuts := groupedShortcuts[keyboard.GroupActions]
	globalShortcuts := groupedShortcuts[keyboard.GroupGlobal]

	// Global edit menu
	EditMenu := AppMenu.AddSubmenu("Edit")
	for action, shortcut := range globalShortcuts {
		EditMenu.AddText(shortcut.Name, shortcut.Accelerator, func(cd *menu.CallbackData) {
			runtime.EventsEmit(a.ctx, string(action))
		})
	}

	// Actions menu
	ActionsMenu := AppMenu.AddSubmenu("Actions")
	for action, shortcut := range actionShortcuts {
		ActionsMenu.AddText(shortcut.Name, shortcut.Accelerator, func(cd *menu.CallbackData) {
			runtime.EventsEmit(a.ctx, string(action))
		})
	}

	// Help menu
	HelpMenu := AppMenu.AddSubmenu("Help")
	HelpMenu.AddText("Report Issue", nil, func(cd *menu.CallbackData) {
		utils.OpenMailClient(utils.SCANOSS_SUPPORT_MAILBOX, "Report an issue", utils.GetIssueReportBody(a.ctx))
	})
	HelpMenu.AddText("Keyboard Shortcuts", keys.Combo("k", keys.ShiftKey, keys.CmdOrCtrlKey), func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "showKeyboardShortcuts")
	})

	runtime.MenuSetApplicationMenu(a.ctx, AppMenu)
}
