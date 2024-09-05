package controllers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
)

type ComponentControllerImpl struct {
	service service.ComponentService
}

func NewComponentController(service service.ComponentService) ComponentController {
	return &ComponentControllerImpl{
		service: service,
	}
}

func (c *ComponentControllerImpl) GetComponentByPath(filePath string) (entities.ComponentDTO, error) {
	component, _ := c.service.GetComponentByFilePath(filePath)
	return adaptToComponentDTO(component), nil
}

func (c *ComponentControllerImpl) FilterComponent(dto entities.ComponentFilterDTO) error {
	return c.service.FilterComponent(dto)
}
