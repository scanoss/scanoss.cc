package adapter

import (
	"integration-git/main/pkg/file/application/services"
	"integration-git/main/pkg/file/application/use_cases"
)

type Controller struct {
	getRemoteFileUseCase *use_cases.GetRemoteFileUseCase
}

func NewFileController(service *services.FileService) *Controller {
	return &Controller{
		getRemoteFileUseCase: use_cases.NewGetRemoteFileUseCase(service),
	}
}

func (c *Controller) GetRemoteFile(path string) {
	c.getRemoteFileUseCase.ReadFile(path)
}

func (c *Controller) GetLocalFile() {

}
