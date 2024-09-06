package internal_test

import (
	"os"
	"testing"

	configEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

func SetupTestFiles(t *testing.T) string {
	t.Helper()

	scanSettingsFile, err := os.CreateTemp("", "scanoss-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	scanSettingsFilePath := scanSettingsFile.Name()

	initialScanSettings := entities.BomFile{
		Bom: entities.Bom{
			Include: []entities.ComponentFilter{},
			Remove:  []entities.ComponentFilter{},
		},
	}
	scanSettingsFile.Close()
	if err := utils.WriteJsonFile(scanSettingsFilePath, initialScanSettings); err != nil {
		t.Fatalf("Failed to write initial scan settings file: %v", err)
	}

	return scanSettingsFilePath
}

func InitializeConfig(scanSettingsFilePath string) *configEntities.Config {
	config := &configEntities.Config{
		ScanSettingsFilePath: scanSettingsFilePath,
	}

	return config
}

func InitializeTestEnvironment(t *testing.T) (*configEntities.Config, func()) {
	t.Helper()

	scanSettingsFilePath := SetupTestFiles(t)

	config := InitializeConfig(scanSettingsFilePath)

	cleanup := func() {
		os.Remove(scanSettingsFilePath)
	}

	return config, cleanup
}
