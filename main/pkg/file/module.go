package file

import (
	"integration-git/main/pkg/file/adapter"
	"integration-git/main/pkg/file/application/services"
	"integration-git/main/pkg/file/infrastructure"
)

type Module struct {
	Controller *adapter.Controller
}

func NewModule() *Module {
	return &Module{
		Controller: adapter.NewFileController(services.NewFileService(infrastructure.NewFileRepository())),
	}
}
