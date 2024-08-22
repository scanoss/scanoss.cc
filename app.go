package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/scanoss/scanoss.lui/backend/handlers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
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
		runtime.WindowSetTitle(a.ctx, fmt.Sprintf("Scanoss Lui %s", version))
	}
}

func (a *App) Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to read user home directory")
		os.Exit(1)
	}
	var defaultScanossFolder = homeDir + string(os.PathSeparator) + "." + config.GetDefaultGlobalFolder() + string(os.PathSeparator) + config.GetDefaultConfigFileName()

	a.createScanossFolder(defaultScanossFolder)
	a.createConfigFile(defaultScanossFolder)
	a.initScanSettingsFile()

}

func (r *App) createScanossFolder(path string) {
	// Get the directory path
	dirPath := filepath.Dir(path)
	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create .scanoss directory: %v\n", err)
			os.Exit(1)
		}
	}
}

func (a *App) createConfigFile(path string) {

	_, err := os.Stat(path)
	if err == nil {
		return // File exists
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Failed to create configuration file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	content := fmt.Sprintf(`{ "apiToken": "" , "apiUrl": "%s"}`, config.GetDefaultApiURL())
	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Failed to write configuration file: %v\n", err)
		os.Exit(1)
	}
}

func (a *App) initScanSettingsFile() {
	scanSettingsFilePath := config.GetScanSettingDefaultLocation()
	dirPath := filepath.Dir(scanSettingsFilePath)
	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create .scanoss directory: %v\n", err)
			os.Exit(1)
		}
		file, err := os.Create(scanSettingsFilePath)
		defer file.Close()
		defaultScanSettings := componentEntities.ScanSettingsFile{
			Bom: componentEntities.Bom{
				Include: []componentEntities.ComponentFilter{},
				Remove:  []componentEntities.ComponentFilter{},
			},
		}

		jsonData, err := json.Marshal(defaultScanSettings)
		if err != nil {
			fmt.Printf("Failed to write configuration file: %v\n", err)
		}
		_, err = file.WriteString(string(jsonData))
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

func (a *App) beforeClose(ctx context.Context, sbh *handlers.ScanossBomHandler) (prevent bool) {

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
		err := sbh.SaveScanossBomFile()
		if err != nil {
			runtime.LogError(ctx, "Error saving scanoss bom file: "+err.Error())
		}
	}

	return false
}
