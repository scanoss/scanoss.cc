package services

import (
	"integration-git/main/pkg/component/application/services"
	"integration-git/main/pkg/file/domain"
)

type FileService struct {
	componentService *services.ComponentService
	repository       domain.FileRepository
}

func NewFileService(r domain.FileRepository, componentService *services.ComponentService) *FileService {
	return &FileService{
		repository:       r,
		componentService: componentService,
	}
}

func (s *FileService) GetLocalFileContent(path string) (domain.File, error) {
	return s.repository.ReadLocalFile(path)
}

func (s *FileService) GetRemoteFileContent(path string) (domain.File, error) {
	component, err := s.componentService.GetComponentByPath(path)

	if err != nil {
		return domain.File{}, err
	}

	f, err := s.repository.ReadRemoteFileByMD5(path, component.FileHash)

	return f, err
}
