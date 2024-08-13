package git

import (
	"integration-git/main/pkg/git/adapter"
	"integration-git/main/pkg/git/application"
	"integration-git/main/pkg/git/infrastructure"
)

type Module struct {
	Controller *adapter.Controller
}

func NewModule() *Module {
	return &Module{
		Controller: adapter.NewGitController(application.NewGitService(infrastructure.NewGitRepository())),
	}
}
