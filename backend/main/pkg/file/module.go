package file

import (
	scanossSettingsRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	componentRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	componentService "github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/service"
	resultRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/result/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type Module struct {
	Controller controllers.FileController
}

func NewModule() *Module {
	fr := utils.NewDefaultFileReader()
	repo := repository.NewFileRepositoryImpl()
	service := service.NewFileService(repo)

	componentRepo := componentRepository.NewJSONComponentRepository(fr)
	scanossSettingsRepo := scanossSettingsRepository.NewScanossSettingsJsonRepository(fr)
	resultRepo := resultRepository.NewResultRepositoryJsonImpl(fr)
	componentService := componentService.NewComponentService(componentRepo, scanossSettingsRepo, resultRepo)

	return &Module{
		Controller: controllers.NewFileController(service, componentService),
	}
}
