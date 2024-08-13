package scan

import (
	"integration-git/main/pkg/scan/adapter"
	"integration-git/main/pkg/scan/application"
)

type Module struct {
	Controller *adapter.Controller
}

func NewModule() *Module {
	return &Module{
		Controller: adapter.NewController(application.NewScanService()),
	}
}
