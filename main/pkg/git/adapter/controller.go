package adapter

import (
	"integration-git/main/pkg/git/application"
)

type Controller struct {
	getFilesToBeCommitedUseCase *application.GetFilesToBeCommitedUseCase
}

func NewGitController(service *application.GitService) *Controller {
	return &Controller{
		getFilesToBeCommitedUseCase: application.NewGetFilesToBeCommitedUseCase(service),
	}
}

func (c *Controller) GetFilesToBeCommited() ([]GitFileDTO, error) {
	var output []GitFileDTO

	files, err := c.getFilesToBeCommitedUseCase.Execute()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		output = append(output, GitFileDTO{Path: file.Path})
	}

	return output, nil
}
