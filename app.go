package main

import (
	"context"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/common/scanoss_bom"
	"integration-git/main/pkg/component"
	componentEntities "integration-git/main/pkg/component/entities"
	"integration-git/main/pkg/file"
	"integration-git/main/pkg/file/adapter"
	"integration-git/main/pkg/git"
	git_module_adapter "integration-git/main/pkg/git/adapter"
	"integration-git/main/pkg/result"
	module_result_adapter "integration-git/main/pkg/result/common"
	"integration-git/main/pkg/scan"
	"os"
	"path/filepath"
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

func (a *App) Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to read user home directory")
		os.Exit(1)
	}

	var defaultScanossFolder = homeDir + string(os.PathSeparator) + "." + config.GetDefaultGlobalFolder() + string(os.PathSeparator) + config.GetDefaultConfigFileName()

	a.createScanossFolder(defaultScanossFolder)
	a.createConfigFile(defaultScanossFolder)

	scanoss_bom.NewScanossBom().Init()

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

func (a *App) onShutDown(ctx context.Context) {
	a.ctx = ctx
	_, err := scanoss_bom.NewScanossBom().Save()
	if err != nil {
		fmt.Println("scanoss.json file not saved!")
	}
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
func (a *App) ComponentGet(filePath string) componentEntities.ComponentDTO {
	comp, _ := a.componentModule.Controller.GetComponentByPath(filePath)
	return comp
}

func (a *App) ComponentFilter(dto componentEntities.ComponentFilterDTO) error {
	return a.componentModule.Controller.FilterComponent(dto)
}
