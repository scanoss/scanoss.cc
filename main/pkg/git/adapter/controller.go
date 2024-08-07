package adapter

import (
	"integration-git/main/pkg/git/application"
)

type Controller struct {
	getFilesToBeCommited *application.GetFilesToBeCommitedUseCase
}

func NewGitController(service *application.GitService) *Controller {
	return &Controller{
		getFilesToBeCommited: application.NewGetFilesToBeCommitedUseCase(service),
	}
}
