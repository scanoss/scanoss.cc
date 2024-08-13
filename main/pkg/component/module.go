package component

import (
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/component/adapter"
	"integration-git/main/pkg/component/application/services"
	"integration-git/main/pkg/component/infraestructure"
)

type Module struct {
	Controller *adapter.ComponentController
}

func NewModule() *Module {
	return &Module{
		Controller: adapter.NewComponentController(services.NewComponentService(infraestructure.NewComponentRepository(config.Get().Scanoss.ResultFilePath))),
	}
}
