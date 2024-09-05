package component

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repositories"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/usecases"
)

type Module struct {
	Controller controllers.ComponentController
}

func NewModule() *Module {
	componentRepository := repositories.NewComponentRepository()
	componentUsecase := usecases.NewComponentUseCase(componentRepository)

	return &Module{
		Controller: controllers.NewComponentController(componentUsecase),
	}
}
