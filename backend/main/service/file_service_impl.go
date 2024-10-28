package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/entities"
	"github.com/scanoss/scanoss.lui/backend/main/repository"
)

type FileServiceImpl struct {
	repo          repository.FileRepository
	componentRepo repository.ComponentRepository
}

func NewFileService(repo repository.FileRepository, componentRepo repository.ComponentRepository) FileService {
	return &FileServiceImpl{
		repo:          repo,
		componentRepo: componentRepo,
	}
}

func (c *FileServiceImpl) GetRemoteFile(path string) (entities.FileDTO, error) {
	component, err := c.componentRepo.FindByFilePath(path)
	if err != nil {
		return entities.FileDTO{}, err
	}

	file, err := c.repo.ReadRemoteFileByMD5(path, component.FileHash)
	return entities.FileDTO{
		Name:     file.GetName(),
		Path:     file.GetRelativePath(),
		Content:  string(file.GetContent()),
		Language: file.GetLanguage(),
	}, err
}

func (c *FileServiceImpl) GetLocalFile(path string) (entities.FileDTO, error) {
	file, err := c.repo.ReadLocalFile(path)
	if err != nil {
		return entities.FileDTO{}, err
	}

	return entities.FileDTO{
		Name:     file.GetName(),
		Path:     file.GetRelativePath(),
		Content:  string(file.GetContent()),
		Language: file.GetLanguage(),
	}, err
}
