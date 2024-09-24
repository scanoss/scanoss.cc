package repository

import (
	"fmt"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type ScanossSettingsJsonRepository struct {
}

func NewScanossSettingsJsonRepository() *ScanossSettingsJsonRepository {
	return &ScanossSettingsJsonRepository{}
}

func (r *ScanossSettingsJsonRepository) Save() error {
	if err := utils.WriteJsonFile(config.Get().ScanSettingsFilePath, entities.ScanossSettingsJson.SettingsFile); err != nil {
		return err
	}
	return nil
}

func (r *ScanossSettingsJsonRepository) Read() (entities.SettingsFile, error) {

	if config.Get() == nil {
		return entities.SettingsFile{}, fmt.Errorf("config is nil")
	}
	scanSettingsFileBytes, err := utils.ReadFile(config.Get().ScanSettingsFilePath)
	if err != nil {
		return entities.SettingsFile{}, err
	}

	scanossSettings, err := utils.JSONParse[entities.SettingsFile](scanSettingsFileBytes)
	if err != nil {
		return entities.SettingsFile{}, err
	}

	return scanossSettings, nil
}

func (r *ScanossSettingsJsonRepository) HasUnsavedChanges() (bool, error) {
	originalBom, err := r.Read()
	if err != nil {
		return false, err
	}

	return !originalBom.Equal(entities.ScanossSettingsJson.SettingsFile), nil
}
