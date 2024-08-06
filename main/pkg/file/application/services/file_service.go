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

func (fs *FileService) GetLocalFileContent(path string) {
	fs.repository.ReadLocalFile(path)
}

func (fs *FileService) GetRemoteFileContent(path string) {
	fs.repository.ReadRemoteFile(path)
}
