package scanoss_settings

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/service"
)

type Module struct {
	Controller controllers.ScanossSettingsController
	Service    service.ScanossSettingsService
}

func NewModule() *Module {
	repository := repository.NewScanossSettingsJsonRepository()
	settingsFile, _ := repository.Read()

	entities.ScanossSettingsJson = &entities.ScanossSettings{
		SettingsFile: &settingsFile,
	}

	service := service.NewScanossSettingsService(repository)

	return &Module{
		Controller: controllers.NewScanossSettingsController(service),
		Service:    service,
	}
}
