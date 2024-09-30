package scanoss_settings

import (
	"log"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type Module struct {
	Controller controllers.ScanossSettingsController
	Service    service.ScanossSettingsService
}

func NewModule() *Module {
	fr := utils.NewDefaultFileReader()
	repository := repository.NewScanossSettingsJsonRepository(fr)
	settingsFile, err := repository.Read()
	if err != nil {
		log.Panicf("error reading scanoss.json file: %s", err)
	}

	entities.ScanossSettingsJson = &entities.ScanossSettings{
		SettingsFile: &settingsFile,
	}

	service := service.NewScanossSettingsService(repository)

	return &Module{
		Controller: controllers.NewScanossSettingsController(service),
		Service:    service,
	}
}
