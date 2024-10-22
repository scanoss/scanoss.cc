package controllers

import (
	"sort"

	"github.com/labstack/gommon/log"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/mappers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type ComponentControllerImpl struct {
	service       service.ComponentService
	mapper        mappers.ComponentMapper
	actionHistory [][]entities.ComponentFilterDTO
	currentAction int
}

func NewComponentController(service service.ComponentService, mapper mappers.ComponentMapper) ComponentController {
	controller := &ComponentControllerImpl{
		service:       service,
		mapper:        mapper,
		actionHistory: [][]entities.ComponentFilterDTO{},
		currentAction: -1,
	}

	controller.initializeActionHistory()
	return controller
}

func (c *ComponentControllerImpl) initializeActionHistory() {
	initialFilters := c.service.GetInitialFilters()
	for _, include := range initialFilters.Include {
		c.actionHistory = append(c.actionHistory, []entities.ComponentFilterDTO{
			{
				Path:   include.Path,
				Purl:   include.Purl,
				Usage:  string(include.Usage),
				Action: entities.Include,
			},
		})
	}
	for _, remove := range initialFilters.Remove {
		c.actionHistory = append(c.actionHistory, []entities.ComponentFilterDTO{
			{
				Path:   remove.Path,
				Purl:   remove.Purl,
				Usage:  string(remove.Usage),
				Action: entities.Remove,
			},
		})
	}

	for _, replace := range initialFilters.Replace {
		c.actionHistory = append(c.actionHistory, []entities.ComponentFilterDTO{
			{
				Path:   replace.Path,
				Purl:   replace.Purl,
				Usage:  string(replace.Usage),
				Action: entities.Replace,
			},
		})
	}

	c.currentAction = len(c.actionHistory) - 1
}

func (c *ComponentControllerImpl) GetComponentByPath(filePath string) (entities.ComponentDTO, error) {
	component, _ := c.service.GetComponentByFilePath(filePath)
	return adaptToComponentDTO(component), nil
}

func (c *ComponentControllerImpl) FilterComponents(dto []entities.ComponentFilterDTO) error {
	for _, filter := range dto {
		err := utils.GetValidator().Struct(filter)
		if err != nil {
			log.Errorf("Validation error: %v", err)
			return err
		}
	}

	err := c.service.FilterComponents(dto)
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
		err := c.service.FilterComponents(c.actionHistory[i])
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

func (c *ComponentControllerImpl) GetDeclaredComponents() ([]entities.DeclaredComponent, error) {
	declaredComponents, err := c.service.GetDeclaredComponents()

	if err != nil {
		return nil, err
	}

	sort.Slice(declaredComponents, func(i, j int) bool {
		return declaredComponents[i].Name < declaredComponents[j].Name
	})

	return declaredComponents, nil
}
