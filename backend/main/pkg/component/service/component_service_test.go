package service_test

import (
	"slices"
	"testing"

	internalTest "github.com/scanoss/scanoss.lui/backend/main/internal"
	scanossBomEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	scanossBomRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestInsertComponentFilterActions(t *testing.T) {
	cleanup := internalTest.InitializeTestEnvironment(t)
	defer cleanup()

	fr := utils.NewDefaultFileReader()
	repo := repository.NewJSONComponentRepository(fr)
	ssRepo := scanossBomRepository.NewScanossSettingsJsonRepository(fr)
	settingsFile, _ := ssRepo.Read()
	scanossBomEntities.ScanossSettingsJson = &scanossBomEntities.ScanossSettings{
		SettingsFile: &settingsFile,
	}
	service := service.NewComponentService(repo, ssRepo)

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
			err := service.FilterComponents([]entities.ComponentFilterDTO{tc.dto})
			require.NoError(t, err)

			var filters = settingsFile.Bom.Include
			if tc.dto.Action == entities.Remove {
				filters = settingsFile.Bom.Remove
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
