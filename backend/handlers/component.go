package handlers

import (
	"github.com/labstack/gommon/log"

	"github.com/go-playground/validator"
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
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		log.Errorf("Validation error: %v", err)
		return err
	}

	return h.componentModule.Controller.FilterComponent(dto)
}
