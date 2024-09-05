package scanoss_bom

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom/service"
)

type Module struct {
	Controller controllers.ScanossBomController
}

func NewModule() *Module {
	repository := repository.NewScanossBomJsonRepository()
	bomFile, _ := repository.Read()

	entities.BomJson = &entities.ScanossBom{
		BomFile: &bomFile,
	}

	service := service.NewScanossBomService(repository)

	return &Module{
		Controller: controllers.NewScanossBomController(service),
	}
}
