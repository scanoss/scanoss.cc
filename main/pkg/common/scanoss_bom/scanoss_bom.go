package scanoss_bom

import (
	"fmt"
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	"integration-git/main/pkg/common/scanoss_bom/application/services"
	"integration-git/main/pkg/common/scanoss_bom/application/use_cases"
	"sync"
)

type ScanossBom struct {
	readScanossJsonUseCase *use_cases.ReadScanossJsonUseCase
	saveScanossJsonUseCase *use_cases.SaveScanossJsonUseCase
	BomFile                *entities.BomFile
}

func NewScanossBom() *ScanossBom {
	return &ScanossBom{
		readScanossJsonUseCase: use_cases.NewReadScanossJsonUseCase(services.NewScanossBomService()),
		saveScanossJsonUseCase: use_cases.NewSaveScanossJsonUseCase(services.NewScanossBomService()),
	}
}

var (
	Bom        *ScanossBom
	once       sync.Once
	configLock sync.Mutex
)

func Read() {
	configLock.Lock()
	defer configLock.Unlock()
	scanossBom := NewScanossBom()
	_, err := scanossBom.read()
	if err != nil {
		fmt.Println(err)
	}
	once.Do(func() {
		Bom = scanossBom
	})
}

// LoadConfig reads the configuration and sets it as the singleton instance
func (sc *ScanossBom) read() (*ScanossBom, error) {
	fmt.Println("Initializing in memory bom")
	var scanossBomFile, _ = sc.readScanossJsonUseCase.Read()
	sc.BomFile = &scanossBomFile
	return sc, nil
}

func Init() {
	createScanossJsonUseCase := use_cases.NewCreateScanossJsonUseCase(services.NewScanossBomService())
	createScanossJsonUseCase.Create()
}

func (sc *ScanossBom) Save() (entities.BomFile, error) {
	return sc.saveScanossJsonUseCase.Save(*sc.BomFile)
}
