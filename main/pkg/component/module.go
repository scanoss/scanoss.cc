package component

import (
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/component/controllers"
	"integration-git/main/pkg/component/repositories"
	"integration-git/main/pkg/component/usecases"
)

type Module struct {
	Controller controllers.ComponentController
}

func NewModule() *Module {
	componentRepository := repositories.NewComponentRepository(config.Get())
	componentUsecase := usecases.NewComponentUseCase(componentRepository)

	return &Module{
		Controller: controllers.NewComponentController(componentUsecase),
	}
}
