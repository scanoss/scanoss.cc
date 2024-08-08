package use_cases

import (
	"integration-git/main/pkg/result/application/services"
	"integration-git/main/pkg/result/common"
	"integration-git/main/pkg/result/domain"
)

type GetResultsUseCase struct {
	s *services.ResultService
}

func NewGetResultsUseCase(s *services.ResultService) *GetResultsUseCase {
	return &GetResultsUseCase{
		s: s,
	}
}

func (us *GetResultsUseCase) GetResults(filter common.ResultFilter) ([]domain.Result, error) {
	return us.s.GetResults(filter)
}
