package file

import (
	scanossSettingsRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	componentRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	componentService "github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type Module struct {
	Controller controllers.FileController
}

func NewModule() *Module {
	fr := utils.NewDefaultFileReader()
	repo := repository.NewFileRepositoryImpl()
	service := service.NewFileService(repo)

	cRepo := componentRepository.NewJSONComponentRepository(fr)
	ssRepo := scanossSettingsRepository.NewScanossSettingsJsonRepository(fr)
	cService := componentService.NewComponentService(cRepo, ssRepo)

	return &Module{
		Controller: controllers.NewFileController(service, cService),
	}
}
