package service

import (
	"fmt"

	scanossSettingsEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	scanossSettingsRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	resultRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/result/repository"
)

type ComponentServiceImpl struct {
	repo                repository.ComponentRepository
	scanossSettingsRepo scanossSettingsRepository.ScanossSettingsRepository
	resultRepo          resultRepository.ResultRepository
}

func NewComponentService(repo repository.ComponentRepository, scanossSettingsRepo scanossSettingsRepository.ScanossSettingsRepository, resultRepo resultRepository.ResultRepository) ComponentService {
	return &ComponentServiceImpl{
		repo:                repo,
		scanossSettingsRepo: scanossSettingsRepo,
		resultRepo:          resultRepo,
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

func (s *ComponentServiceImpl) GetDeclaredComponents() ([]entities.DeclaredComponent, error) {
	results, err := s.resultRepo.GetResults(nil)
	if err != nil {
		return []entities.DeclaredComponent{}, err
	}

	scanossSettingsDeclaredPurls := s.scanossSettingsRepo.GetDeclaredPurls()

	purlToComponent := make(map[string]string)

	for _, result := range results {
		for _, purl := range result.Purl {
			purlToComponent[purl] = result.ComponentName
		}
	}

	addedPurls := make(map[string]struct{})
	declaredComponents := make([]entities.DeclaredComponent, 0)

	for _, result := range results {
		purl := result.Purl[0]
		if _, found := addedPurls[purl]; !found {
			declaredComponents = append(declaredComponents, entities.DeclaredComponent{
				Purl: purl,
				Name: result.ComponentName,
			})
			addedPurls[purl] = struct{}{}
		}
	}

	for _, purl := range scanossSettingsDeclaredPurls {
		if _, found := addedPurls[purl]; !found {
			if component, found := purlToComponent[purl]; found {
				declaredComponents = append(declaredComponents, entities.DeclaredComponent{
					Purl: purl,
					Name: component,
				})
			} else {
				declaredComponents = append(declaredComponents, entities.DeclaredComponent{
					Purl: purl,
				})
			}
			addedPurls[purl] = struct{}{}
		}
	}

	return declaredComponents, nil
}
