package infraestructure

import (
	"encoding/json"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	"integration-git/main/pkg/utils"
	"os"
)

type ScanossBomJsonRepository struct {
}

func NewScanossBomJonRepository() *ScanossBomJsonRepository {
	return &ScanossBomJsonRepository{}
}

func (r *ScanossBomJsonRepository) Save(bomFile entities.BomFile) error {
	if err := utils.WriteJsonFile(config.Get().ScanSettingsFilePath, bomFile); err != nil {
		return err
	}
	return nil
}

func (r *ScanossBomJsonRepository) Init() (entities.BomFile, error) {
	scanSettingsFilePath := config.Get().ScanSettingsFilePath

	// Check if the directory exists
	if _, err := os.Stat(scanSettingsFilePath); os.IsNotExist(err) {
		file, err := os.Create(scanSettingsFilePath)
		defer file.Close()
		defaultScanSettings := entities.BomFile{
			Bom: entities.Bom{
				Include: []entities.ComponentFilter{},
				Remove:  []entities.ComponentFilter{},
			},
		}
		jsonData, err := json.Marshal(defaultScanSettings)
		if err != nil {
			fmt.Printf("Failed to write configuration file: %v\n", err)
		}
		_, err = file.WriteString(string(jsonData))

		return defaultScanSettings, nil
	}

	return entities.BomFile{}, nil
}

func (r *ScanossBomJsonRepository) Read() (entities.BomFile, error) {

	if config.Get() == nil {
		return entities.BomFile{}, fmt.Errorf("config is nil")
	}

	scanSettingsFileBytes, err := utils.ReadFile(config.Get().ScanSettingsFilePath)
	if err != nil {
		return entities.BomFile{}, err
	}

	scanossBom, err := utils.JSONParse[entities.BomFile](scanSettingsFileBytes)
	if err != nil {
		return entities.BomFile{}, err
	}

	return scanossBom, nil
}
