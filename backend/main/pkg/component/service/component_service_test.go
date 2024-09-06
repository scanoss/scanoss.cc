package service_test

import (
	"encoding/json"
	"os"
	"slices"
	"testing"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	scanossBomEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_bom/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	"github.com/stretchr/testify/require"
)

func TestInsertComponentFilterActions(t *testing.T) {
	//config, cleanup := internal_test.InitializeTestEnvironment(t)
	// defer cleanup()

	repo := repository.NewJSONComponentRepository()

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

			var scanSettings scanossBomEntities.BomFile
			fileBytes, err := os.ReadFile(config.Get().ScanSettingsFilePath)
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

			i := slices.IndexFunc(filters, func(cf scanossBomEntities.ComponentFilter) bool {
				return cf.Path == tc.dto.Path
			})

			filter := filters[i]

			require.Equal(t, tc.dto.Path, filter.Path)
			require.Equal(t, tc.dto.Purl, filter.Purl)
		})
	}
}
