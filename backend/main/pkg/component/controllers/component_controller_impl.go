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
	service        service.ComponentService
	mapper         mappers.ComponentMapper
	initialFilters []entities.ComponentFilterDTO
	undoStack      [][]entities.ComponentFilterDTO
	redoStack      [][]entities.ComponentFilterDTO
}

func NewComponentController(service service.ComponentService, mapper mappers.ComponentMapper) ComponentController {
	controller := &ComponentControllerImpl{
		service:        service,
		initialFilters: []entities.ComponentFilterDTO{},
		undoStack:      [][]entities.ComponentFilterDTO{},
		redoStack:      [][]entities.ComponentFilterDTO{},
		mapper:         mapper,
	}

	controller.setInitialFilters()

	return controller
}

func (c *ComponentControllerImpl) setInitialFilters() {
	initialFilters := c.service.GetInitialFilters()

	for _, include := range initialFilters.Include {
		c.initialFilters = append(c.initialFilters, entities.ComponentFilterDTO{
			Path:   include.Path,
			Purl:   include.Purl,
			Usage:  string(include.Usage),
			Action: entities.Include,
		})
	}
	for _, remove := range initialFilters.Remove {
		c.initialFilters = append(c.initialFilters, entities.ComponentFilterDTO{
			Path:   remove.Path,
			Purl:   remove.Purl,
			Usage:  string(remove.Usage),
			Action: entities.Remove,
		})
	}
	for _, replace := range initialFilters.Replace {
		c.initialFilters = append(c.initialFilters, entities.ComponentFilterDTO{
			Path:   replace.Path,
			Purl:   replace.Purl,
			Usage:  string(replace.Usage),
			Action: entities.Replace,
		})
	}
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

	if err := c.service.FilterComponents(dto); err != nil {
		return err
	}

	c.undoStack = append(c.undoStack, dto)
	c.redoStack = nil

	return nil
}

func (c *ComponentControllerImpl) Undo() error {
	if !c.CanUndo() {
		return nil
	}

	lastAction := c.undoStack[len(c.undoStack)-1]
	c.undoStack = c.undoStack[:len(c.undoStack)-1]
	c.redoStack = append(c.redoStack, lastAction)

	return c.reapplyActions()
}

func (c *ComponentControllerImpl) Redo() error {
	if !c.CanRedo() {
		return nil
	}

	nextAction := c.redoStack[len(c.redoStack)-1]
	c.redoStack = c.redoStack[:len(c.redoStack)-1]
	c.undoStack = append(c.undoStack, nextAction)

	return c.service.FilterComponents(nextAction)
}

func (c *ComponentControllerImpl) CanUndo() bool {
	return len(c.undoStack) > 0
}

func (c *ComponentControllerImpl) CanRedo() bool {
	return len(c.redoStack) > 0
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

func (c *ComponentControllerImpl) reapplyActions() error {
	// We need to clear bom filters so workflow state is properly resetted
	err := c.service.ClearAllFilters()
	if err != nil {
		return err
	}

	for _, initialFilterDto := range c.initialFilters {
		err := c.service.FilterComponents([]entities.ComponentFilterDTO{initialFilterDto})
		if err != nil {
			return err
		}
	}

	for _, filterDto := range c.undoStack {
		err := c.service.FilterComponents(filterDto)
		if err != nil {
			return err
		}
	}

	return nil
}
