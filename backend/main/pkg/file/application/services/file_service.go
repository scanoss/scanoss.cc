package services

import (
	"integration-git/main/pkg/component/usecases"
	"integration-git/main/pkg/file/domain"
)

type FileService struct {
	componentUsecase usecases.ComponentUsecase
	repository       domain.FileRepository
}

func NewFileService(r domain.FileRepository, componentUsecase usecases.ComponentUsecase) *FileService {
	return &FileService{
		repository:       r,
		componentUsecase: componentUsecase,
	}
}

func (s *FileService) GetLocalFileContent(path string) (domain.File, error) {
	return s.repository.ReadLocalFile(path)
}

func (s *FileService) GetRemoteFileContent(path string) (domain.File, error) {
	component, err := s.componentUsecase.GetComponentByFilePath(path)

	if err != nil {
		return domain.File{}, err
	}

	f, err := s.repository.ReadRemoteFileByMD5(path, component.FileHash)

	return f, err
}
