package controllers

import "github.com/scanoss/scanoss.lui/backend/main/internal/common/scanoss_bom/service"

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
