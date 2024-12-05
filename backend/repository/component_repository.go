package repository

import "github.com/scanoss/scanoss.lui/backend/entities"

type ComponentRepository interface {
	FindByFilePath(path string) (entities.Component, error)
}
