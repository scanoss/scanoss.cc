package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/entities"
	"github.com/scanoss/scanoss.lui/backend/main/repository"
)

type FileServiceImpl struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return &FileServiceImpl{
		repo: repo,
	}
}

func (c *FileServiceImpl) GetRemoteFile(path string) (entities.FileDTO, error) {
	component, err := c.repo.GetComponentByFilePath(path)
	if err != nil {
		return entities.FileDTO{}, err
	}
	file, err := c.repo.ReadRemoteFileByMD5(path, component.FileHash)
	return entities.FileDTO{
		Name:    file.GetName(),
		Path:    file.GetRelativePath(),
		Content: string(file.GetContent()),
	}, err
}

func (c *FileServiceImpl) GetLocalFile(path string) (entities.FileDTO, error) {
	file, err := c.repo.ReadLocalFile(path)

	return entities.FileDTO{
		Name:     file.GetName(),
		Path:     file.GetRelativePath(),
		Content:  string(file.GetContent()),
		Language: file.GetLanguage(),
	}, err
}
