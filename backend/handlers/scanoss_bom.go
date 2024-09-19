package handlers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom"
)

type ScanossBomHandler struct {
	scanossBomModule *scanoss_bom.Module
}

func NewScanossBomHandler() *ScanossBomHandler {
	return &ScanossBomHandler{
		scanossBomModule: scanoss_bom.NewModule(),
	}
}

func (sh *ScanossBomHandler) SaveScanossBomFile() error {
	return sh.scanossBomModule.Controller.Save()
}

func (sh *ScanossBomHandler) HasUnsavedChanges() (bool, error) {
	return sh.scanossBomModule.Controller.HasUnsavedChanges()
}
