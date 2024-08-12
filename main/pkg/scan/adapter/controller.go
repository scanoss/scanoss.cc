package adapter

import (
	"integration-git/main/pkg/scan/application"
	"integration-git/main/pkg/scan/domain"
)

type Controller struct {
	scanUseCase *application.ScanUseCase
}

func NewController(service *application.ScanService) *Controller {
	return &Controller{
		scanUseCase: application.NewScanUseCase(service),
	}
}

func (c *Controller) RunScan(request ScanInputPathsDTO) error {
	scanPaths, err := domain.NewScanInputPaths(request.BasePath, request.Paths)
	if err != nil {
		return err
	}

	_, err = c.scanUseCase.Execute(scanPaths)
	return err
}
