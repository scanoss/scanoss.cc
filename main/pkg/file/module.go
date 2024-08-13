package file

import (
	"integration-git/main/pkg/common/config"
	services2 "integration-git/main/pkg/component/application/services"
	infraestructure2 "integration-git/main/pkg/component/infraestructure"
	"integration-git/main/pkg/file/adapter"
	"integration-git/main/pkg/file/application/services"
	"integration-git/main/pkg/file/infrastructure"
)

type Module struct {
	Controller *adapter.Controller
}

func NewModule() *Module {

	componentService := services2.NewComponentService(infraestructure2.NewComponentRepository(config.Get().Scanoss.ResultFilePath))

	return &Module{
		Controller: adapter.NewFileController(services.NewFileService(infrastructure.NewFileRepository(), componentService)),
	}
}
