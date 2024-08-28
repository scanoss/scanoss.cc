package use_cases

import (
	"integration-git/main/pkg/common/scanoss_bom/infraestructure"
	modules "integration-git/main/pkg/common/scanoss_bom/module"
)

type ScanossBomUseCaseImp struct {
	repository infraestructure.ScanossBomRepository
}

func NewScanossBomUseCase(r infraestructure.ScanossBomRepository) *ScanossBomUseCaseImp {
	return &ScanossBomUseCaseImp{
		repository: r,
	}
}

func (us *ScanossBomUseCaseImp) Save() error {
	return us.repository.Save(*modules.Bom.BomFile)
}
