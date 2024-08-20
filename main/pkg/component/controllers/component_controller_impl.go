package controllers

import (
	"integration-git/main/pkg/component/entities"
	"integration-git/main/pkg/component/usecases"
)

type ComponentControllerImpl struct {
	componentUseCase usecases.ComponentUsecase
}

func NewComponentController(componentUsecase usecases.ComponentUsecase) ComponentController {
	return &ComponentControllerImpl{
		componentUseCase: componentUsecase,
	}
}

func (c *ComponentControllerImpl) GetComponentByPath(filePath string) (entities.ComponentDTO, error) {
	component, _ := c.componentUseCase.GetComponentByFilePath(filePath)
	return adaptToComponentDTO(component), nil
}
