package adapter

import (
	"integration-git/main/pkg/file/application/services"
	"integration-git/main/pkg/file/application/use_cases"
)

type Controller struct {
	getRemoteFileUseCase *use_cases.GetRemoteFileUseCase
	getLocalFileUseCase  *use_cases.GetLocalFileUseCase
}

func NewFileController(service *services.FileService) *Controller {
	return &Controller{
		getRemoteFileUseCase: use_cases.NewGetRemoteFileUseCase(service),
		getLocalFileUseCase:  use_cases.NewGetLocalFileUseCase(service),
	}
}

func (c *Controller) GetRemoteFile(path string) (FileDTO, error) {
	file, err := c.getRemoteFileUseCase.ReadFile(path)
	return FileDTO{
		Name:    file.GetName(),
		Path:    file.GetRelativePath(),
		Content: string(file.GetContent()),
	}, err
}

func (c *Controller) GetLocalFile(path string) (FileDTO, error) {
	file, err := c.getLocalFileUseCase.ReadFile(path)

	return FileDTO{
		Name:    file.GetName(),
		Path:    file.GetRelativePath(),
		Content: string(file.GetContent()),
	}, err
}
