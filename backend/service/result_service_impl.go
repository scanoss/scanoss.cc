package service

import (
	"sort"

	"github.com/labstack/gommon/log"
	"github.com/scanoss/scanoss.lui/backend/entities"
	"github.com/scanoss/scanoss.lui/backend/mappers"
	"github.com/scanoss/scanoss.lui/backend/repository"
	"github.com/scanoss/scanoss.lui/internal/utils"
)

type ResultServiceImpl struct {
	repo   repository.ResultRepository
	mapper mappers.ResultMapper
}

func NewResultServiceImpl(repo repository.ResultRepository, mapper mappers.ResultMapper) ResultService {
	return &ResultServiceImpl{
		repo:   repo,
		mapper: mapper,
	}
}

func (s *ResultServiceImpl) GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error) {
	err := utils.GetValidator().Struct(dto)
	if err != nil {
		log.Errorf("Validation error: %v", err)
		return []entities.ResultDTO{}, err
	}

	filter := entities.NewResultFilterFactory().Create(dto)
	results, err := s.repo.GetResults(filter)
	if err != nil {
		return []entities.ResultDTO{}, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})

	return s.mapper.MapToResultDTOList(results), nil
}
