package main

import (
	"context"
	"fmt"
	"integration-git/main/pkg/common/config"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
	/*	fileModule        *file.Module
		gitModule         *git.Module
		resultModule      *result.Module
		componentModule   *component.Module
		scanModule        *scan.Module
		scannossBomModule *scanoss_bom.Module*/
}

// NewApp creates a new App application struct
func NewApp() *App {
	/*	return &App{
		fileModule:        file.NewModule(),
		gitModule:         git.NewModule(),
		resultModule:      result.NewModule(),
		componentModule:   component.NewModule(),
		scanModule:        scan.NewModule(),
		scannossBomModule: scanoss_bom.NewModule(),
	}*/
	return &App{}
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

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
