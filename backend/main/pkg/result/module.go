package result

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/mappers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type Module struct {
	Controller controllers.ResultController
}

func NewModule(bomModule *scanoss_settings.Module) *Module {
	fr := utils.NewDefaultFileReader()
	repo := repository.NewResultRepositoryJsonImpl(fr)
	serv := service.NewResultServiceImpl(repo)
	mapper := mappers.NewResultMapper(bomModule.Service)

	return &Module{
		Controller: controllers.NewResultController(serv, mapper),
	}
}
