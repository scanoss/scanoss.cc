package main

import (
	"context"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/file"
	"integration-git/main/pkg/file/adapter"
	"integration-git/main/pkg/git"
	gitAdapter "integration-git/main/pkg/git/adapter"
	"integration-git/main/pkg/result"
	resultAdapter "integration-git/main/pkg/result/common"
	"log"
)

// App struct
type App struct {
	ctx          context.Context
	fileModule   *file.Module
	gitModule    *git.Module
	resultModule *result.Module
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		fileModule:   file.NewModule(),
		gitModule:    git.NewModule(),
		resultModule: result.NewModule(),
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
	cfg := config.Get()
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
func (a *App) FileGetLocalContent(path string) adapter.FileDTO {
	f, _ := a.fileModule.Controller.GetLocalFile(path)
	return f
}

func (a *App) FileGetRemoteContent(path string) adapter.FileDTO {
	f, _ := a.fileModule.Controller.GetRemoteFile(path)
	return f
}

// *****  RESULT MODULE ***** //
func (a *App) ResultGetAll(dto *resultAdapter.RequestResultDTO) []resultAdapter.ResultDTO {
	results, _ := a.resultModule.Controller.GetAll(dto)
	return results
}
