package use_cases

import (
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	"integration-git/main/pkg/common/scanoss_bom/application/services"
)

type ReadScanossJsonUseCase struct {
	s *services.ScanossBomService
}

func NewReadScanossJsonUseCase(s *services.ScanossBomService) *ReadScanossJsonUseCase {
	return &ReadScanossJsonUseCase{
		s: s,
	}
}

func (us *ReadScanossJsonUseCase) Read() (entities.BomFile, error) {
	return us.s.Read()
}
