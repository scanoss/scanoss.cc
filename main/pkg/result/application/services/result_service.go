package services

import (
	"integration-git/main/pkg/result/common"
	"integration-git/main/pkg/result/domain"
	"integration-git/main/pkg/result/repository"
)

type ResultService struct {
	r repository.ResultRepository
}

func NewResultService(r repository.ResultRepository) *ResultService {
	return &ResultService{
		r: r,
	}
}

func (s ResultService) GetResults(filter common.ResultFilter) ([]domain.Result, error) {
	return s.r.GetResults(filter)
}
