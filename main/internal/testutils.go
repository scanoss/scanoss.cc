package internal_test

import (
	"integration-git/main/pkg/common/config/domain"
	"integration-git/main/pkg/component/entities"
	"integration-git/main/pkg/utils"
	"os"
	"testing"
)

func SetupTestFiles(t *testing.T) string {
	t.Helper()

	scanSettingsFile, err := os.CreateTemp("", "scanoss-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	scanSettingsFilePath := scanSettingsFile.Name()

	initialScanSettings := entities.ScanSettingsFile{
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

func InitializeConfig(scanSettingsFilePath string) *domain.Config {
	config := &domain.Config{
		ScanSettingsFilePath: scanSettingsFilePath,
	}

	return config
}

func InitializeTestEnvironment(t *testing.T) (*domain.Config, func()) {
	t.Helper()

	scanSettingsFilePath := SetupTestFiles(t)

	config := InitializeConfig(scanSettingsFilePath)

	cleanup := func() {
		os.Remove(scanSettingsFilePath)
	}

	return config, cleanup
}
