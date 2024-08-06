package use_cases

import "integration-git/main/pkg/file/application/services"

type GetRemoteFileUseCase struct {
	fileService *services.FileService
}

func NewGetRemoteFileUseCase(s *services.FileService) *GetRemoteFileUseCase {
	return &GetRemoteFileUseCase{fileService: s}
}

func (s *GetRemoteFileUseCase) ReadFile(path string) {

	s.fileService.GetRemoteFileContent(path)
}
