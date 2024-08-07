package use_cases

import (
	"integration-git/main/pkg/file/application/services"
	"integration-git/main/pkg/file/domain"
)

type GetLocalFileUseCase struct {
	fileService *services.FileService
}

func NewGetLocalFileUseCase(s *services.FileService) *GetLocalFileUseCase {
	return &GetLocalFileUseCase{fileService: s}
}

func (uc *GetLocalFileUseCase) ReadFile(path string) (domain.File, error) {
	return uc.fileService.GetLocalFileContent(path)
}
