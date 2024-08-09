package application

import (
	"integration-git/main/pkg/scan/domain"
)

type ScanUseCase struct {
	scanService *ScanService
}

func NewScanUseCase(service *ScanService) *ScanUseCase {
	return &ScanUseCase{
		scanService: service,
	}
}

func (uc *ScanUseCase) Execute(scanInput domain.ScanInputPaths) (domain.ScanResult, error) {
	return uc.scanService.Scan(scanInput)
}
