package use_cases

import (
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	"integration-git/main/pkg/common/scanoss_bom/application/services"
)

type CreateScanossJsonUseCase struct {
	s *services.ScanossBomService
}

func NewCreateScanossJsonUseCase(s *services.ScanossBomService) *CreateScanossJsonUseCase {
	return &CreateScanossJsonUseCase{
		s: s,
	}
}

func (us *CreateScanossJsonUseCase) Create() (entities.BomFile, error) {
	return us.s.Create()
}
