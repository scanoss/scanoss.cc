package service

import "github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"

type ComponentService interface {
	GetComponentByFilePath(filePath string) (entities.Component, error)
	FilterComponent(dto entities.ComponentFilterDTO) error
}
