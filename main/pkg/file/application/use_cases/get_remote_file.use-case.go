package use_cases

import (
	"integration-git/main/pkg/file/application/services"
	"integration-git/main/pkg/file/domain"
)

type GetRemoteFileUseCase struct {
	fileService *services.FileService
}

func NewGetRemoteFileUseCase(s *services.FileService) *GetRemoteFileUseCase {
	return &GetRemoteFileUseCase{fileService: s}
}

func (uc *GetRemoteFileUseCase) ReadFile(path string) (domain.File, error) {
	return uc.fileService.GetRemoteFileContent(path)
}
