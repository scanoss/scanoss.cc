package controllers

import "integration-git/main/pkg/component/entities"

type ComponentController interface {
	GetComponentByPath(filePath string) (entities.ComponentDTO, error)
	FilterComponent(dto entities.ComponentFilterDTO) error
}
