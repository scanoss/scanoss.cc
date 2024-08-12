package use_cases

import (
	"integration-git/main/pkg/component/application/services"
	"integration-git/main/pkg/component/domain"
)

type GetComponentByFilePathUseCase struct {
	s *services.ComponentService
}

func NewGetComponentByFilePathUseCase(s *services.ComponentService) *GetComponentByFilePathUseCase {
	return &GetComponentByFilePathUseCase{
		s: s,
	}
}

func (uc *GetComponentByFilePathUseCase) GetComponentByPath(filePath string) (domain.Component, error) {
	return uc.s.GetComponentByPath(filePath)
}
