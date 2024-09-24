package handlers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings"
)

type ScanossSettingsHandler struct {
	scanossSettingsModule *scanoss_settings.Module
}

func NewScanossSettingsHandler() *ScanossSettingsHandler {
	return &ScanossSettingsHandler{
		scanossSettingsModule: scanoss_settings.NewModule(),
	}
}

func (sh *ScanossSettingsHandler) SaveScanossSettingsFile() error {
	return sh.scanossSettingsModule.Controller.Save()
}

func (sh *ScanossSettingsHandler) HasUnsavedChanges() (bool, error) {
	return sh.scanossSettingsModule.Controller.HasUnsavedChanges()
}

func (sh *ScanossSettingsHandler) GetScanossSettings() *scanoss_settings.Module {
	return sh.scanossSettingsModule
}
