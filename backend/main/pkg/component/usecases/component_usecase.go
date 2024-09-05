package usecases

import (
	"integration-git/main/pkg/component/entities"
)

type ComponentUsecase interface {
	GetComponentByFilePath(filePath string) (entities.Component, error)
	FilterComponent(dto entities.ComponentFilterDTO) error
}
