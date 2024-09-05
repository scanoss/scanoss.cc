package repositories

import (
	"integration-git/main/pkg/component/entities"
)

type ComponentRepository interface {
	FindByFilePath(path string) (entities.Component, error)
	InsertComponentFilter(dto *entities.ComponentFilterDTO) error
}
