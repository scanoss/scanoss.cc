package modules

import (
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	"sync"
)

type ScanossBom struct {
	BomFile *entities.BomFile
}

func NewScanossBomModule() *ScanossBom {
	return &ScanossBom{}
}

var (
	Bom        *ScanossBom
	once       sync.Once
	configLock sync.Mutex
)

func (sc *ScanossBom) Init(bomFile *entities.BomFile) {
	configLock.Lock()
	defer configLock.Unlock()
	once.Do(func() {
		sc.BomFile = bomFile
		Bom = sc
	})
}
