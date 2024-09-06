package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/repository"
)

type FileServiceImpl struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return &FileServiceImpl{
		repo: repo,
	}
}

func (s *FileServiceImpl) GetLocalFileContent(path string) (entities.File, error) {
	return s.repo.ReadLocalFile(path)
}

func (s *FileServiceImpl) GetRemoteFileContent(path string, md5 string) (entities.File, error) {
	f, err := s.repo.ReadRemoteFileByMD5(path, md5)
	return f, err
}
