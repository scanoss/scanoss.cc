package handler

import (
	"integration-git/main/pkg/component"
	"integration-git/main/pkg/component/entities"
)

type ComponentHandler struct {
	componentModule *component.Module
}

func NewComponentHandler() *ComponentHandler {
	return &ComponentHandler{
		componentModule: component.NewModule(),
	}
}

func (ch *ComponentHandler) ComponentGet(filePath string) entities.ComponentDTO {
	comp, _ := ch.componentModule.Controller.GetComponentByPath(filePath)
	return comp
}

func (ch *ComponentHandler) ComponentFilter(dto entities.ComponentFilterDTO) error {
	return ch.componentModule.Controller.FilterComponent(dto)
}
