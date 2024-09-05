package handlers

import (
	"context"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom"
)

type ScanossBomHandler struct {
	ctx               context.Context
	scannossBomModule *scanoss_bom.Module
}

func NewScanossBomHandler() *ScanossBomHandler {
	return &ScanossBomHandler{
		scannossBomModule: scanoss_bom.NewModule(),
	}
}

// *****  SCANOSS BOM MODULE ***** //
func (sh *ScanossBomHandler) SaveScanossBomFile() error {
	return sh.scannossBomModule.Controller.Save()
}

func (sh *ScanossBomHandler) OnShutDown(ctx context.Context) {
	sh.ctx = ctx
	sh.scannossBomModule.Controller.Save()
}
