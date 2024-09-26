package controllers

import "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/service"

type ScanossSettingsControllerImpl struct {
	service service.ScanossSettingsService
}

func NewScanossSettingsController(service service.ScanossSettingsService) ScanossSettingsController {
	return &ScanossSettingsControllerImpl{
		service: service,
	}
}

func (c *ScanossSettingsControllerImpl) Save() error {
	return c.service.Save()
}

func (c *ScanossSettingsControllerImpl) HasUnsavedChanges() (bool, error) {
	return c.service.HasUnsavedChanges()
}
