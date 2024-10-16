package service

import (
	scanossSettingsEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
)

type ComponentService interface {
	GetComponentByFilePath(filePath string) (entities.Component, error)
	FilterComponents(dto []entities.ComponentFilterDTO) error
	ClearAllFilters() error
	GetInitialFilters() ([]scanossSettingsEntities.ComponentFilter, []scanossSettingsEntities.ComponentFilter)
	GetDeclaredComponents() ([]entities.DeclaredComponent, error)
}
