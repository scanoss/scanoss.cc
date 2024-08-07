package use_cases

import (
	"integration-git/main/pkg/common/config/application/services"
	"integration-git/main/pkg/common/config/domain"
)

type ReadJsonConfigFileUSeCase struct {
	s *services.ConfigService
}

func NewReadJsonConfigFileUseCase(s *services.ConfigService) *ReadJsonConfigFileUSeCase {
	return &ReadJsonConfigFileUSeCase{
		s: s,
	}
}

func (us ReadJsonConfigFileUSeCase) ReadConfig(path string) (domain.Config, error) {
	return us.s.ReadFromJson(path)
}
