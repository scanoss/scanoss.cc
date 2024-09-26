package internal_test

import (
	"os"
	"testing"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"

	configEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

func SetupTestFiles(t *testing.T) string {
	t.Helper()

	scanSettingsFile, err := os.CreateTemp("", "scanoss-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	scanSettingsFilePath := scanSettingsFile.Name()

	initialScanSettings := entities.SettingsFile{
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

func InitializeTestEnvironment(t *testing.T) func() {
	t.Helper()

	scanSettingsFilePath := SetupTestFiles(t)

	InitializeConfig(scanSettingsFilePath)

	cleanup := func() {
		os.Remove(scanSettingsFilePath)
	}

	cfg := config.NewConfigModule(scanSettingsFilePath)
	cfg.Init()
	cfg.LoadConfig()

	r := repository.NewScanossSettingsJsonRepository()
	settingsFile, err := r.Read()
	if err != nil {
		t.Fatal(err.Error())
	}
	// Init Scanoss bom module. Set current bom file to singleton
	entities.ScanossSettingsJson = &entities.ScanossSettings{
		SettingsFile: &settingsFile,
	}

	return cleanup
}
