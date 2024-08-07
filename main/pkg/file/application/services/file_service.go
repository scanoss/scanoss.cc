package services

import (
	"integration-git/main/pkg/file/domain"
)

type FileService struct {
	repository domain.FileRepository
}

func NewFileService(r domain.FileRepository) *FileService {
	return &FileService{repository: r}
}

func (fs *FileService) GetLocalFileContent(path string) (domain.File, error) {
	return fs.repository.ReadLocalFile(path)
}

func (fs *FileService) GetRemoteFileContent(path string) (domain.File, error) {
	return fs.repository.ReadRemoteFile(path)
}
