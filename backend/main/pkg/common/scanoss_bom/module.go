package scanoss_bom

import (
	"github.com/scanoss/scanoss.lui/backend/main/internal/common/scanoss_bom/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/internal/common/scanoss_bom/entities"
	"github.com/scanoss/scanoss.lui/backend/main/internal/common/scanoss_bom/repository"
	"github.com/scanoss/scanoss.lui/backend/main/internal/common/scanoss_bom/service"
)

type Module struct {
	Controller controllers.ScanossBomController
}

func NewModule() *Module {
	// Init repository to read singleton bom file
	repository := repository.NewScanossBomJsonRepository()
	repository.Init()
	bomFile, _ := repository.Read()
	entities.BomJson.Init(&bomFile)

	service := service.NewScanossBomService(repository)

	return &Module{
		Controller: controllers.NewScanossBomController(service),
	}
}
