package repository

import "github.com/scanoss/scanoss.cc/backend/entities"

type ComponentRepository interface {
	FindByFilePath(path string) (entities.Component, error)
}
