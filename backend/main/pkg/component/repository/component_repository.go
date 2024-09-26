package repository

import "github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"

type ComponentRepository interface {
	FindByFilePath(path string) (entities.Component, error)
}
