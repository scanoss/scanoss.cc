package main

import (
	"context"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/file"
	"integration-git/main/pkg/file/adapter"
	"integration-git/main/pkg/git"
	gitAdapter "integration-git/main/pkg/git/adapter"
	"log"
)

// App struct
type App struct {
	ctx        context.Context
	fileModule *file.Module
	gitModule  *git.Module
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		fileModule: file.NewModule(),
		gitModule:  git.NewModule(),
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

// *****  GIT MODULE ***** //
func (a *App) GetFilesToBeCommited() []gitAdapter.GitFileDTO {
	f, _ := a.gitModule.Controller.GetFilesToBeCommited()
	return f
}

// *****  FILE MODULE ***** //
func (a *App) GetLocalFileContent(path string) adapter.FileDTO {
	f, _ := a.fileModule.Controller.GetLocalFile()
	return f
}

func (a *App) FileRemoteFileContent(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
