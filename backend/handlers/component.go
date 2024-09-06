package handlers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
)

type ComponentHandler struct {
	componentModule *component.Module
}

func NewComponentHandler() *ComponentHandler {
	return &ComponentHandler{
		componentModule: component.NewModule(),
	}
}

func (h *ComponentHandler) ComponentGet(filePath string) entities.ComponentDTO {
	comp, _ := h.componentModule.Controller.GetComponentByPath(filePath)
	return comp
}

func (h *ComponentHandler) ComponentFilter(dto entities.ComponentFilterDTO) error {
	return h.componentModule.Controller.FilterComponent(dto)
}
