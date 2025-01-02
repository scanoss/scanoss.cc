package service

import "github.com/scanoss/scanoss.cc/backend/entities"

type ComponentService interface {
	GetComponentByPath(filePath string) (entities.ComponentDTO, error)
	FilterComponents(dto []entities.ComponentFilterDTO) error
	Undo() error
	Redo() error
	CanUndo() bool
	CanRedo() bool
	GetDeclaredComponents() ([]entities.DeclaredComponent, error)
}
