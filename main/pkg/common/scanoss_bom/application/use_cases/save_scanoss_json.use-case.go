package use_cases

import (
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	"integration-git/main/pkg/common/scanoss_bom/application/services"
)

type SaveScanossJsonUseCase struct {
	s *services.ScanossBomService
}

func NewSaveScanossJsonUseCase(s *services.ScanossBomService) *SaveScanossJsonUseCase {
	return &SaveScanossJsonUseCase{
		s: s,
	}
}

func (us *SaveScanossJsonUseCase) Save(bomFile entities.BomFile) (entities.BomFile, error) {
	return us.s.Save(bomFile)
}
