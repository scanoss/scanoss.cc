package controllers

import (
	componentService "github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/service"
)

type FileControllerImpl struct {
	service          service.FileService
	componentService componentService.ComponentService
}

func NewFileController(service service.FileService, componentService componentService.ComponentService) FileController {
	return &FileControllerImpl{
		service:          service,
		componentService: componentService,
	}
}

func (c *FileControllerImpl) GetRemoteFile(path string) (entities.FileDTO, error) {
	component, err := c.componentService.GetComponentByFilePath(path)
	if err != nil {
		return entities.FileDTO{}, err
	}
	file, err := c.service.GetRemoteFileContent(path, component.FileHash)
	return entities.FileDTO{
		Name:    file.GetName(),
		Path:    file.GetRelativePath(),
		Content: string(file.GetContent()),
	}, err
}

func (c *FileControllerImpl) GetLocalFile(path string) (entities.FileDTO, error) {
	file, err := c.service.GetLocalFileContent(path)

	return entities.FileDTO{
		Name:    file.GetName(),
		Path:    file.GetRelativePath(),
		Content: string(file.GetContent()),
	}, err
}
