package component

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
)

type Module struct {
	Controller controllers.ComponentController
}

func NewModule() *Module {
	componentRepository := repository.NewJSONComponentRepository()
	componentService := service.NewComponentService(componentRepository)

	return &Module{
		Controller: controllers.NewComponentController(componentService),
	}
}
