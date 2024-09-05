package repository

import (
	"fmt"
	"os"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type ScanossBomJsonRepository struct {
}

func NewScanossBomJsonRepository() *ScanossBomJsonRepository {
	return &ScanossBomJsonRepository{}
}

func (r *ScanossBomJsonRepository) Save() error {
	if err := utils.WriteJsonFile(config.Get().ScanSettingsFilePath, entities.BomJson); err != nil {
		return err
	}
	return nil
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
