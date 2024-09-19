package controllers

import "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom/service"

type ScanossBomControllerImpl struct {
	service service.ScanossBomService
}

func NewScanossBomController(service service.ScanossBomService) ScanossBomController {
	return &ScanossBomControllerImpl{
		service: service,
	}
}

func (c *ScanossBomControllerImpl) Save() error {
	return c.service.Save()
}

func (c *ScanossBomControllerImpl) HasUnsavedChanges() (bool, error) {
	return c.service.HasUnsavedChanges()
}
