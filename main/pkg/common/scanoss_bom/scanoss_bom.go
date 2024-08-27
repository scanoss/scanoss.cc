package scanoss_bom

import (
	"fmt"
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	"integration-git/main/pkg/common/scanoss_bom/application/services"
	"integration-git/main/pkg/common/scanoss_bom/application/use_cases"
	"sync"
)

type ScanossBom struct {
	readScanossJsonUseCase   *use_cases.ReadScanossJsonUseCase
	createScanossJsonUseCase *use_cases.CreateScanossJsonUseCase
	saveScanossJsonUseCase   *use_cases.SaveScanossJsonUseCase
}

func NewScanossBom() *ScanossBom {
	return &ScanossBom{
		readScanossJsonUseCase:   use_cases.NewReadScanossJsonUseCase(services.NewScanossBomService()),
		createScanossJsonUseCase: use_cases.NewCreateScanossJsonUseCase(services.NewScanossBomService()),
		saveScanossJsonUseCase:   use_cases.NewSaveScanossJsonUseCase(services.NewScanossBomService()),
	}
}

var (
	bom        *entities.BomFile
	once       sync.Once
	configLock sync.Mutex
)

// LoadConfig reads the configuration and sets it as the singleton instance
func (sc *ScanossBom) Read() (*entities.BomFile, error) {
	fmt.Println("Initializing in memory bom")
	configLock.Lock()
	defer configLock.Unlock()
	var scanossBom, err = sc.readScanossJsonUseCase.Read()
	if err != nil {
		scanossBom, _ = sc.createScanossJsonUseCase.Create()
	}

	once.Do(func() {
		bom = &scanossBom
	})
	return bom, nil
}

func (sc *ScanossBom) Init() {
	sc.createScanossJsonUseCase.Create()
}

func Get() *entities.BomFile {
	return bom
}

func (sc *ScanossBom) Save() (entities.BomFile, error) {
	return sc.saveScanossJsonUseCase.Save(*bom)
}
