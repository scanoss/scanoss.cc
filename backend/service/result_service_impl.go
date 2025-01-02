package service

import (
	"context"
	"sort"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/mappers"
	"github.com/scanoss/scanoss.cc/backend/repository"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type ResultServiceImpl struct {
	ctx     context.Context
	watcher *fsnotify.Watcher
	repo    repository.ResultRepository
	mapper  mappers.ResultMapper
}

func NewResultServiceImpl(repo repository.ResultRepository, mapper mappers.ResultMapper) ResultService {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("Error creating watcher: %v")
	}

	return &ResultServiceImpl{
		repo:    repo,
		mapper:  mapper,
		watcher: watcher,
	}
}

func (s *ResultServiceImpl) GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error) {
	err := utils.GetValidator().Struct(dto)
	if err != nil {
		log.Error().Err(err).Msg("Validation error: %v")
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

func (s *ResultServiceImpl) SetContext(ctx context.Context) {
	s.ctx = ctx
}
