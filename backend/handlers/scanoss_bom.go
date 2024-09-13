package handlers

import (
	"context"
	"log"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom"
)

type ScanossBomHandler struct {
	ctx              context.Context
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

func (sh *ScanossBomHandler) OnShutDown(ctx context.Context) {
	sh.ctx = ctx
	err := sh.scanossBomModule.Controller.Save()
	if err != nil {
		log.Fatalf("Error saving scanoss bom file: %s", err)
	}
}
