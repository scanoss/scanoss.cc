package result

import (
	"integration-git/main/pkg/result/adapter"
	"integration-git/main/pkg/result/application/services"
	"integration-git/main/pkg/result/infraestructure"
)

type Module struct {
	Controller *adapter.ResultController
}

func NewModule() *Module {
	return &Module{
		Controller: adapter.NewResultController(services.NewResultService(infraestructure.NewInMemoryResultRepository())),
	}
}
