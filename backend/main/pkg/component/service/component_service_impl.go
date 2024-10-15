package service

import (
	"fmt"

	scanossSettingsEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	scanossSettingsRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
)

type ComponentServiceImpl struct {
	repo                repository.ComponentRepository
	scanossSettingsRepo scanossSettingsRepository.ScanossSettingsRepository
}

func NewComponentService(repo repository.ComponentRepository, scanossSettingsRepo scanossSettingsRepository.ScanossSettingsRepository) ComponentService {
	return &ComponentServiceImpl{
		repo:                repo,
		scanossSettingsRepo: scanossSettingsRepo,
	}
}

func (s *ComponentServiceImpl) GetComponentByFilePath(filePath string) (entities.Component, error) {
	return s.repo.FindByFilePath(filePath)
}

func (s *ComponentServiceImpl) FilterComponents(dto []entities.ComponentFilterDTO) error {
	for _, item := range dto {
		newFilter := &scanossSettingsEntities.ComponentFilter{
			Path:    item.Path,
			Purl:    item.Purl,
			Usage:   scanossSettingsEntities.ComponentFilterUsage(item.Usage),
			Comment: item.Comment,
		}
		if err := s.scanossSettingsRepo.AddBomEntry(*newFilter, string(item.Action)); err != nil {
			fmt.Printf("error adding bom entry: %s", err)
			return err
		}
	}

	return nil
}

func (s *ComponentServiceImpl) ClearAllFilters() error {
	return s.scanossSettingsRepo.ClearAllFilters()
}

func (s *ComponentServiceImpl) GetInitialFilters() ([]scanossSettingsEntities.ComponentFilter, []scanossSettingsEntities.ComponentFilter) {
	sf := s.scanossSettingsRepo.GetSettingsFileContent()
	include, remove := sf.Bom.Include, sf.Bom.Remove

	return include, remove
}
