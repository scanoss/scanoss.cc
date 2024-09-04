package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"integration-git/main/pkg/common/config"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.WindowSetTitle(ctx, fmt.Sprintf("Scanoss Lui %s", version))
	runtime.LogInfo(ctx, "Config file path: "+config.GetInstance().GetConfigPath())
	runtime.LogInfo(ctx, "Results file path: "+config.Get().ResultFilePath)
	runtime.LogInfo(ctx, "Scan Root file path: "+config.Get().ScanRoot)
}
