package controllers

import (
	"github.com/go-playground/validator"
	"github.com/labstack/gommon/log"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
)

type ComponentControllerImpl struct {
	service       service.ComponentService
	actionHistory []entities.ComponentFilterDTO
	currentAction int
}

func NewComponentController(service service.ComponentService) ComponentController {
	return &ComponentControllerImpl{
		service:       service,
		actionHistory: []entities.ComponentFilterDTO{},
		currentAction: -1,
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

	err = c.service.FilterComponent(dto)
	if err != nil {
		return err
	}

	c.currentAction++
	c.actionHistory = append(c.actionHistory[:c.currentAction], dto)

	return nil
}

func (c *ComponentControllerImpl) Undo() error {
	if c.currentAction < 0 {
		return nil
	}

	action := c.actionHistory[c.currentAction]
	inverseAction := c.getInverseAction(action)

	err := c.service.FilterComponent(inverseAction)
	if err != nil {
		return err
	}

	c.currentAction--
	return nil
}

func (c *ComponentControllerImpl) Redo() error {
	if c.currentAction >= len(c.actionHistory)-1 {
		return nil
	}

	c.currentAction++
	action := c.actionHistory[c.currentAction]

	return c.service.FilterComponent(action)
}

func (c *ComponentControllerImpl) getInverseAction(action entities.ComponentFilterDTO) entities.ComponentFilterDTO {
	inverse := action
	if action.Action == entities.Include {
		inverse.Action = entities.Remove
	} else {
		inverse.Action = entities.Include
	}
	return inverse
}

func (c *ComponentControllerImpl) CanUndo() (bool, error) {
	return c.currentAction >= 0, nil
}

func (c *ComponentControllerImpl) CanRedo() (bool, error) {
	return c.currentAction < len(c.actionHistory)-1, nil
}
