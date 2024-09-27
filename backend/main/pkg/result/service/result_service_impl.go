package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/repository"
)

type ResultServiceImpl struct {
	repo repository.ResultRepository
}

func NewResultServiceImpl(repo repository.ResultRepository) ResultService {
	return &ResultServiceImpl{
		repo: repo,
	}
}

func (s *ResultServiceImpl) GetResults(filters entities.ResultFilter) ([]entities.Result, error) {
	return s.repo.GetResults(filters)
}
