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
	repo := repository.NewJSONComponentRepository()
	service := service.NewComponentService(repo)

	return &Module{
		Controller: controllers.NewComponentController(service),
	}
}
