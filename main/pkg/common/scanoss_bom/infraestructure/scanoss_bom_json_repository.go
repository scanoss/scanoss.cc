package infraestructure

import (
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	"integration-git/main/pkg/utils"
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
