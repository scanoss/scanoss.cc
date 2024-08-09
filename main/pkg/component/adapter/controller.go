package adapter

import (
	"integration-git/main/pkg/component/application/services"
	"integration-git/main/pkg/component/application/use_cases"
)

type ComponentController struct {
	getComponentByFilePathUseCase *use_cases.GetComponentByFilePathUseCase
}

func NewComponentController(s *services.ComponentService) *ComponentController {
	return &ComponentController{
		getComponentByFilePathUseCase: use_cases.NewGetComponentByFilePathUseCase(s),
	}
}

func (c *ComponentController) GetComponentByPath(filePath string) (ComponentDTO, error) {
	component, _ := c.getComponentByFilePathUseCase.GetComponentByPath(filePath)
	return adaptToComponentDTO(component), nil
}
