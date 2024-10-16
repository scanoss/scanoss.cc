package component

import (
	scanossSettingsRepo "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/mappers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	resultRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/result/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type Module struct {
	Controller controllers.ComponentController
}

func NewModule() *Module {
	fr := utils.NewDefaultFileReader()
	componentRepository := repository.NewJSONComponentRepository(fr)
	scanossSettingsRepository := scanossSettingsRepo.NewScanossSettingsJsonRepository(fr)
	resultRepo := resultRepository.NewResultRepositoryJsonImpl(fr)
	componentService := service.NewComponentService(componentRepository, scanossSettingsRepository, resultRepo)
	mapper := mappers.NewComponentMapper()

	return &Module{
		Controller: controllers.NewComponentController(componentService, mapper),
	}
}
