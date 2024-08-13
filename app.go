package main

import (
	"context"
	"fmt"
	"integration-git/main/cmd"
	componentadapter "integration-git/main/pkg/component/adapter"
	fileadapter "integration-git/main/pkg/file/adapter"
	gitmoduleadapter "integration-git/main/pkg/git/adapter"
	moduleresultadapter "integration-git/main/pkg/result/common"
)

// App struct
type App struct {
	ctx    context.Context
	server *cmd.Server
}

// NewApp creates a new App application struct
func NewApp(container *cmd.Server) *App {
	return &App{
		server: container,
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
func (a *App) GetFilesToBeCommited() []gitmoduleadapter.GitFileDTO {
	f, _ := a.server.GitModule.Controller.GetFilesToBeCommited()
	return f
}

// *****  FILE MODULE ***** //
func (a *App) FileGetLocalContent(path string) fileadapter.FileDTO {
	f, _ := a.server.FileModule.Controller.GetLocalFile(path)
	return f
}

func (a *App) FileGetRemoteContent(path string) fileadapter.FileDTO {
	f, _ := a.server.FileModule.Controller.GetRemoteFile(path)
	return f
}

// *****  RESULT MODULE ***** //
func (a *App) ResultGetAll(dto *moduleresultadapter.RequestResultDTO) []moduleresultadapter.ResultDTO {
	results, _ := a.server.ResultModule.Controller.GetAll(dto)
	return results
}

// *****  COMPONENT MODULE ***** //
func (a *App) ComponentGet(filePath string) componentadapter.ComponentDTO {
	comp, _ := a.server.ComponentModule.Controller.GetComponentByPath(filePath)
	return comp
}
