package result

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/service"
)

type Module struct {
	Controller controllers.ResultController
}

func NewModule() *Module {
	repo := repository.NewResultRepositoryJsonImpl()
	serv := service.NewResultServiceImpl(repo)

	return &Module{
		Controller: controllers.NewResultController(serv),
	}
}
