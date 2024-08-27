package services

import (
	"encoding/json"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	"integration-git/main/pkg/utils"
	"os"
	"path/filepath"
)

type ScanossBomService struct {
}

func NewScanossBomService() *ScanossBomService {
	return &ScanossBomService{}
}

func (s *ScanossBomService) Read() (entities.BomFile, error) {

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

func (s *ScanossBomService) Create() (entities.BomFile, error) {
	scanSettingsFilePath := config.GetScanSettingDefaultLocation()
	dirPath := filepath.Dir(scanSettingsFilePath)
	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create .scanoss directory: %v\n", err)
			os.Exit(1)
		}
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

func (s *ScanossBomService) Save(bomFile entities.BomFile) (entities.BomFile, error) {
	if err := utils.WriteJsonFile(config.Get().ScanSettingsFilePath, bomFile); err != nil {
		return entities.BomFile{}, err
	}
	return bomFile, nil
}
