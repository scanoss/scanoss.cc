package controllers

import "github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"

type ComponentController interface {
	GetComponentByPath(filePath string) (entities.ComponentDTO, error)
	FilterComponents(dto []entities.ComponentFilterDTO) error
	Undo() error
	Redo() error
	CanUndo() bool
	CanRedo() bool
	GetDeclaredComponents() ([]entities.DeclaredComponent, error)
}
