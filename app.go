package main

import (
	"context"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/file/adapter"
	"log"
)

// App struct
type App struct {
	ctx            context.Context
	fileController adapter.Controller
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		fileController: adapter.Controller{},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	//Read config
	_, err := config.LoadConfig("config.json")

	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	cfg := config.GetConfig()
	fmt.Println(cfg.Scanoss.ApiToken)
	fmt.Println(cfg.Scanoss.ApiUrl)

}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// File
func (a *App) FileLocalFileContent(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) FileRemoteFileContent(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Git
