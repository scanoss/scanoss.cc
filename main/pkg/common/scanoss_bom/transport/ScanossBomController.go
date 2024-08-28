package transport

import (
	"integration-git/main/pkg/common/scanoss_bom/application/use_cases"
	"integration-git/main/pkg/common/scanoss_bom/infraestructure"
)

type ScanossBomController struct {
	scanossBomUseCase use_cases.ScanossBomUseCase
}

func NewScanossBomController() *ScanossBomController {
	// Init Scanoss bom repository.
	r := infraestructure.NewScanossBomJonRepository()
	r.Init()
	return &ScanossBomController{
		scanossBomUseCase: use_cases.NewScanossBomUseCase(r),
	}
}

func (c *ScanossBomController) Save() {
	c.scanossBomUseCase.Save()
}
