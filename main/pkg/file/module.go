package file

import (
	"integration-git/main/pkg/common/config"
	service "integration-git/main/pkg/component/application/services"
	"integration-git/main/pkg/component/infraestructure"
	"integration-git/main/pkg/file/adapter"
	"integration-git/main/pkg/file/application/services"
	"integration-git/main/pkg/file/infrastructure"
)

type Module struct {
	Controller *adapter.Controller
}

func NewModule() *Module {

	componentService := service.NewComponentService(infraestructure.NewComponentRepository(config.Get().ResultFilePath))

	return &Module{
		Controller: adapter.NewFileController(services.NewFileService(infrastructure.NewFileRepository(), componentService)),
	}
}
