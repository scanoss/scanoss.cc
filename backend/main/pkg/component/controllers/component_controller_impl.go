package controllers

import (
	"github.com/go-playground/validator"
	"github.com/labstack/gommon/log"
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
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		log.Errorf("Validation error: %v", err)
		return err
	}

	return c.service.FilterComponent(dto)
}
