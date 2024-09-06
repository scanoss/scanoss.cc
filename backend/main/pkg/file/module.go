package file

import (
	componentRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	componentService "github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/service"
)

type Module struct {
	Controller controllers.FileController
}

func NewModule() *Module {
	repo := repository.NewFileRepositoryImpl()
	service := service.NewFileService(repo)

	cRepo := componentRepository.NewJSONComponentRepository()
	cService := componentService.NewComponentService(cRepo)

	return &Module{
		Controller: controllers.NewFileController(service, cService),
	}
}
