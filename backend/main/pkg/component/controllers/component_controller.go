package controllers

import "github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"

type ComponentController interface {
	GetComponentByPath(filePath string) (entities.ComponentDTO, error)
	FilterComponent(dto entities.ComponentFilterDTO) error
}
