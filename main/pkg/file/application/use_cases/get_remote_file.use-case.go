package use_cases

import "integration-git/main/pkg/file/application/services"

type GetRemoteFileUseCase struct {
	fileService *services.FileService
}

func NewGetRemoteFileUseCase(s *services.FileService) *GetRemoteFileUseCase {
	return &GetRemoteFileUseCase{fileService: s}
}

func (uc *GetRemoteFileUseCase) ReadFile(path string) {

	uc.fileService.GetRemoteFileContent(path)
}
