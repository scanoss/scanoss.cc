package component

import (
	scanossSettingsRepo "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
)

type Module struct {
	Controller controllers.ComponentController
}

func NewModule() *Module {
	componentRepository := repository.NewJSONComponentRepository()
	scanossSettingsRepository := scanossSettingsRepo.NewScanossSettingsJsonRepository()
	componentService := service.NewComponentService(componentRepository, scanossSettingsRepository)

	return &Module{
		Controller: controllers.NewComponentController(componentService),
	}
}
