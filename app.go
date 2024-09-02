package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"integration-git/main/pkg/common/config"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) Init() {
	// Creates .scanoss hidden folder in home dir
	a.createScanossFolder(config.GetDefaultConfigLocation())
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

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogInfo(ctx, "Config file path: "+config.GetConfigPath())
	runtime.LogInfo(ctx, "Results file path: "+config.Get().ResultFilePath)
	runtime.LogInfo(ctx, "Scan Root file path: "+config.Get().ScanRoot)
}
