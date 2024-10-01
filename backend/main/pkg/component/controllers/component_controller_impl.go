package controllers

import (
	"github.com/go-playground/validator"
	"github.com/labstack/gommon/log"
	scanossSettingsEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
)

type ComponentControllerImpl struct {
	service       service.ComponentService
	actionHistory []entities.ComponentFilterDTO
	currentAction int
}

func NewComponentController(service service.ComponentService) ComponentController {
	controller := &ComponentControllerImpl{
		service:       service,
		actionHistory: []entities.ComponentFilterDTO{},
		currentAction: -1,
	}

	controller.initializeActionHistory()
	return controller
}

func (c *ComponentControllerImpl) initializeActionHistory() {
	settings := scanossSettingsEntities.ScanossSettingsJson
	for _, include := range settings.SettingsFile.Bom.Include {
		c.actionHistory = append(c.actionHistory, entities.ComponentFilterDTO{
			Path:   include.Path,
			Purl:   include.Purl,
			Usage:  string(include.Usage),
			Action: entities.Include,
		})
	}
	for _, remove := range settings.SettingsFile.Bom.Remove {
		c.actionHistory = append(c.actionHistory, entities.ComponentFilterDTO{
			Path:   remove.Path,
			Purl:   remove.Purl,
			Usage:  string(remove.Usage),
			Action: entities.Remove,
		})
	}
	c.currentAction = len(c.actionHistory) - 1
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

	c.currentAction--
	return c.reapplyActions()
}

func (c *ComponentControllerImpl) Redo() error {
	if c.currentAction >= len(c.actionHistory)-1 {
		return nil
	}

	c.currentAction++
	return c.reapplyActions()
}

func (c *ComponentControllerImpl) reapplyActions() error {
	// We need to clear bom filters so workflow state is properly resetted
	err := c.service.ClearAllFilters()
	if err != nil {
		return err
	}

	for i := 0; i <= c.currentAction; i++ {
		err := c.service.FilterComponent(c.actionHistory[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ComponentControllerImpl) CanUndo() (bool, error) {
	return c.currentAction >= 0, nil
}

func (c *ComponentControllerImpl) CanRedo() (bool, error) {
	return c.currentAction < len(c.actionHistory)-1, nil
}
