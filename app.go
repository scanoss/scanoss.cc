package main

import (
	"context"
	"fmt"
	"integration-git/main/pkg/component"
	componentAdapter "integration-git/main/pkg/component/adapter"
	"integration-git/main/pkg/file"
	"integration-git/main/pkg/file/adapter"
	"integration-git/main/pkg/git"
	git_module_adapter "integration-git/main/pkg/git/adapter"
	"integration-git/main/pkg/result"
	module_result_adapter "integration-git/main/pkg/result/common"
	"integration-git/main/pkg/scan"
)

// App struct
type App struct {
	ctx             context.Context
	fileModule      *file.Module
	gitModule       *git.Module
	resultModule    *result.Module
	componentModule *component.Module
	scanModule      *scan.Module
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		fileModule:      file.NewModule(),
		gitModule:       git.NewModule(),
		resultModule:    result.NewModule(),
		componentModule: component.NewModule(),
		scanModule:      scan.NewModule(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// *****  GIT MODULE ***** //
func (a *App) GetFilesToBeCommited() []git_module_adapter.GitFileDTO {
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
func (a *App) ResultGetAll(dto *module_result_adapter.RequestResultDTO) []module_result_adapter.ResultDTO {
	results, _ := a.resultModule.Controller.GetAll(dto)
	return results
}

// *****  COMPONENT MODULE ***** //
func (a *App) ComponentGet(filePath string) componentAdapter.ComponentDTO {
	comp, _ := a.componentModule.Controller.GetComponentByPath(filePath)
	return comp
}
