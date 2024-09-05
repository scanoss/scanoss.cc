package transport

import (
	"integration-git/main/pkg/common/scanoss_bom/application/use_cases"
	"integration-git/main/pkg/common/scanoss_bom/infraestructure"
	modules "integration-git/main/pkg/common/scanoss_bom/module"
)

type ScanossBomController struct {
	scanossBomUseCase use_cases.ScanossBomUseCase
}

func NewScanossBomController() *ScanossBomController {
	// Init Scanoss bom repository.
	r := infraestructure.NewScanossBomJonRepository()
	bomFile, _ := r.Read()
	// Init Scanoss bom module. Set current bom file to singleton
	modules.NewScanossBomModule().Init(&bomFile)
	return &ScanossBomController{
		scanossBomUseCase: use_cases.NewScanossBomUseCase(r),
	}
}

func (c *ScanossBomController) Save() error {
	return c.scanossBomUseCase.Save()
}
