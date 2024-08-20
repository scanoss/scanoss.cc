package usecases_test

import (
	"encoding/json"
	"integration-git/main/pkg/common/config/domain"
	"integration-git/main/pkg/component/entities"
	"integration-git/main/pkg/component/repositories"
	"integration-git/main/pkg/utils"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupTestFiles(t *testing.T) (string, func()) {
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

	cleanup := func() {
		os.Remove(scanSettingsFilePath)
	}

	return scanSettingsFilePath, cleanup
}

func initializeConfig(scanSettingsFilePath string) *domain.Config {
	config := &domain.Config{
		ScanSettingsFilePath: scanSettingsFilePath,
	}

	return config
}

func TestInsertComponentFilterActions(t *testing.T) {
	scanSettingsFilePath, cleanup := setupTestFiles(t)
	defer cleanup()

	config := initializeConfig(scanSettingsFilePath)

	repo := repositories.NewComponentRepository(config)

	tests := []struct {
		name string
		dto  entities.ComponentFilterDTO
	}{
		{
			name: "Include action",
			dto: entities.ComponentFilterDTO{
				Path:   "test/path1",
				Purl:   "pkg:purl1",
				Action: entities.Include,
			},
		},
		{
			name: "Remove action",
			dto: entities.ComponentFilterDTO{
				Path:   "test/path2",
				Purl:   "pkg:purl2",
				Action: entities.Remove,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.InsertComponentFilter(&tc.dto)
			require.NoError(t, err)

			var scanSettings entities.ScanSettingsFile
			fileBytes, err := os.ReadFile(scanSettingsFilePath)
			if err != nil {
				t.Fatalf("Failed to read scan settings file: %v", err)
			}

			if err := json.Unmarshal(fileBytes, &scanSettings); err != nil {
				t.Fatalf("Failed to unmarshal scan settings file: %v", err)
			}

			var filters = scanSettings.Bom.Include
			if tc.dto.Action == entities.Remove {
				filters = scanSettings.Bom.Remove
			}

			i := slices.IndexFunc(filters, func(cf entities.ComponentFilter) bool {
				return cf.Path == tc.dto.Path
			})

			filter := filters[i]

			require.Equal(t, tc.dto.Path, filter.Path)
			require.Equal(t, tc.dto.Purl, filter.Purl)
		})
	}
}
