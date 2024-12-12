package service

import (
	"context"
	"sort"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.lui/backend/entities"
	"github.com/scanoss/scanoss.lui/backend/mappers"
	"github.com/scanoss/scanoss.lui/backend/repository"
	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/scanoss/scanoss.lui/internal/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func (s *ResultServiceImpl) WatchResults() {
	if s.watcher == nil {
		return
	}

	err := s.watcher.Add(config.GetInstance().ResultFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Error watching results file")
		return
	}

	log.Debug().Msgf("Watching %s for changes...", config.GetInstance().ResultFilePath)

	go func() {
		for {
			select {
			case event, ok := <-s.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					runtime.WindowReloadApp(s.ctx)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					runtime.WindowReloadApp(s.ctx)
				}

			case err, ok := <-s.watcher.Errors:
				if !ok {
					return
				}
				log.Error().Err(err).Msg("Error watching results file")
			case <-s.ctx.Done():
				s.watcher.Close()
				return
			}
		}
	}()
}
