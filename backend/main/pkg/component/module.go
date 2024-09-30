package component

import (
	scanossSettingsRepo "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type Module struct {
	Controller controllers.ComponentController
}

func NewModule() *Module {
	fr := utils.NewDefaultFileReader()
	componentRepository := repository.NewJSONComponentRepository(fr)
	scanossSettingsRepository := scanossSettingsRepo.NewScanossSettingsJsonRepository(fr)
	componentService := service.NewComponentService(componentRepository, scanossSettingsRepository)

	return &Module{
		Controller: controllers.NewComponentController(componentService),
	}
}
