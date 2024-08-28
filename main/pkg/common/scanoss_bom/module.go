package scanoss_bom

import (
	"integration-git/main/pkg/common/scanoss_bom/transport"
)

type Module struct {
	Controller *transport.ScanossBomController
}

func NewModule() *Module {
	return &Module{
		// Init Scanoss bom repository
		Controller: transport.NewScanossBomController(),
	}
}
