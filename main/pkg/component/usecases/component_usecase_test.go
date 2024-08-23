package usecases_test

import (
	"encoding/json"
	"fmt"
	internal_test "integration-git/main/internal"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/common/scanoss_bom/application/entities"
	modules "integration-git/main/pkg/common/scanoss_bom/module"
	componentEntities "integration-git/main/pkg/component/entities"
	"integration-git/main/pkg/component/repositories"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertComponentFilterActions(t *testing.T) {
	mockConfig, cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()
	config.SetConfig(mockConfig)
	modules.NewScanossBomModule().Init(&entities.BomFile{})
	repo := repositories.NewComponentRepository()

	tests := []struct {
		name string
		dto  componentEntities.ComponentFilterDTO
	}{
		{
			name: "Include action",
			dto: componentEntities.ComponentFilterDTO{
				Path:   "test/path1",
				Purl:   "pkg:purl1",
				Action: componentEntities.Include,
			},
		},
		{
			name: "Remove action",
			dto: componentEntities.ComponentFilterDTO{
				Path:   "test/path2",
				Purl:   "pkg:purl2",
				Action: componentEntities.Remove,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.InsertComponentFilter(&tc.dto)
			require.NoError(t, err)

			var scanSettings entities.BomFile
			fileBytes, err := os.ReadFile(config.Get().ScanSettingsFilePath)
			if err != nil {
				t.Fatalf("Failed to read scan settings file: %v", err)
			}

			if err := json.Unmarshal(fileBytes, &scanSettings); err != nil {
				t.Fatalf("Failed to unmarshal scan settings file: %v", err)
			}

			var filters = scanSettings.Bom.Include
			if tc.dto.Action == componentEntities.Remove {
				filters = scanSettings.Bom.Remove
			}

			i := slices.IndexFunc(filters, func(cf entities.ComponentFilter) bool {
				return cf.Path == tc.dto.Path
			})

			fmt.Println("FILTERS", i)
			fmt.Println("FILTERS", filters)

			//		require.Equal(t, tc.dto.Path, filter.Path)
			//		require.Equal(t, tc.dto.Purl, filter.Purl)
		})
	}
}
