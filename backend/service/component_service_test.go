//go:build unit

package service_test

import (
	"slices"
	"testing"

	"github.com/scanoss/scanoss.lui/backend/entities"
	"github.com/scanoss/scanoss.lui/backend/mappers"
	"github.com/scanoss/scanoss.lui/backend/mappers/mocks"
	"github.com/scanoss/scanoss.lui/backend/repository"
	mocksRepo "github.com/scanoss/scanoss.lui/backend/repository/mocks"
	"github.com/scanoss/scanoss.lui/backend/service"
	internal_test "github.com/scanoss/scanoss.lui/internal"
	"github.com/scanoss/scanoss.lui/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertComponentFilterActions(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	fr := utils.NewDefaultFileReader()
	repo := repository.NewJSONComponentRepository(fr)
	resultRepo := repository.NewResultRepositoryJsonImpl(fr)
	scanossSettingsRepo := repository.NewScanossSettingsJsonRepository(fr)
	componentMapper := mappers.NewComponentMapper()

	scanossSettingsRepo.Init()
	settingsFile, _ := scanossSettingsRepo.Read()
	entities.ScanossSettingsJson = &entities.ScanossSettings{
		SettingsFile: &settingsFile,
	}
	service := service.NewComponentServiceImpl(repo, scanossSettingsRepo, resultRepo, componentMapper)

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

			i := slices.IndexFunc(filters, func(cf entities.ComponentFilter) bool {
				return cf.Path == tc.dto.Path
			})

			filter := filters[i]

			require.Equal(t, tc.dto.Path, filter.Path)
			require.Equal(t, tc.dto.Purl, filter.Purl)
		})
	}
}
func TestGetDeclaredComponents(t *testing.T) {
	t.Run("WHEN Scanoss Settings is empty -> SHOULD return detected components from results file", func(t *testing.T) {
		repo := mocksRepo.NewMockComponentRepository(t)
		resultRepo := mocksRepo.NewMockResultRepository(t)
		scanossSettingsRepo := mocksRepo.NewMockScanossSettingsRepository(t)
		componentMapper := mocks.NewMockComponentMapper(t)
		results := []entities.Result{
			{
				ComponentName: "component1",
				Path:          "path/to/file",
				MatchType:     "file",
				Purl:          &[]string{"pkg:npm/purl1"},
			},
			{
				ComponentName: "component2",
				Path:          "path/to/file2",
				MatchType:     "snippet",
				Purl:          &[]string{"pkg:github/purl2"},
			},
		}

		resultRepo.EXPECT().GetResults(nil).Return(results, nil)
		scanossSettingsRepo.EXPECT().GetSettings().Return(&entities.SettingsFile{})
		scanossSettingsRepo.EXPECT().GetDeclaredPurls().Return(nil)

		service := service.NewComponentServiceImpl(repo, scanossSettingsRepo, resultRepo, componentMapper)

		declaredComponents, err := service.GetDeclaredComponents()
		assert.NoError(t, err)
		assert.NotEmpty(t, declaredComponents)
		assert.Len(t, declaredComponents, 2)
		assert.Equal(t, declaredComponents, []entities.DeclaredComponent{
			{
				Purl: "pkg:npm/purl1",
				Name: "component1",
			}, {
				Purl: "pkg:github/purl2",
				Name: "component2",
			},
		})

		resultRepo.AssertExpectations(t)
		scanossSettingsRepo.AssertExpectations(t)
	})

	t.Run("WHEN Scanoss Settings is NOT empty -> SHOULD return detected components from results file and declared components in scanoss settings", func(t *testing.T) {
		repo := mocksRepo.NewMockComponentRepository(t)
		resultRepo := mocksRepo.NewMockResultRepository(t)
		scanossSettingsRepo := mocksRepo.NewMockScanossSettingsRepository(t)
		componentMapper := mocks.NewMockComponentMapper(t)

		results := []entities.Result{
			{
				ComponentName: "component1",
				Path:          "path/to/file",
				MatchType:     "file",
				Purl:          &[]string{"pkg:npm/purl1"},
			},
			{
				ComponentName: "component2",
				Path:          "path/to/file2",
				MatchType:     "snippet",
				Purl:          &[]string{"pkg:github/purl2"},
			},
		}

		resultRepo.EXPECT().GetResults(nil).Return(results, nil)
		scanossSettingsRepo.EXPECT().GetSettings().Return(&entities.SettingsFile{
			Bom: entities.Bom{
				Include: []entities.ComponentFilter{
					{
						Purl: "pkg:npm/purl3",
					},
				},
			},
		})
		scanossSettingsRepo.EXPECT().GetDeclaredPurls().Return([]string{"pkg:github/purl3"})

		service := service.NewComponentServiceImpl(repo, scanossSettingsRepo, resultRepo, componentMapper)

		declaredComponents, err := service.GetDeclaredComponents()
		assert.NoError(t, err)
		assert.NotEmpty(t, declaredComponents)
		assert.Len(t, declaredComponents, 3)
		assert.Equal(t, declaredComponents, []entities.DeclaredComponent{
			{
				Purl: "pkg:npm/purl1",
				Name: "component1",
			},
			{
				Purl: "pkg:github/purl2",
				Name: "component2",
			},
			{
				Purl: "pkg:github/purl3",
				Name: "purl3",
			},
		})

		resultRepo.AssertExpectations(t)
		scanossSettingsRepo.AssertExpectations(t)
		componentMapper.AssertExpectations(t)
	})

	t.Run("WHEN there are repeated purls SHOULD NOT return duplicated values", func(t *testing.T) {
		repo := mocksRepo.NewMockComponentRepository(t)
		resultRepo := mocksRepo.NewMockResultRepository(t)
		scanossSettingsRepo := mocksRepo.NewMockScanossSettingsRepository(t)
		componentMapper := mocks.NewMockComponentMapper(t)

		results := []entities.Result{
			{
				ComponentName: "component1",
				Path:          "path/to/file",
				MatchType:     "file",
				Purl:          &[]string{"pkg:npm/purl1"},
			},
			{
				ComponentName: "component2",
				Path:          "path/to/file2",
				MatchType:     "snippet",
				Purl:          &[]string{"pkg:github/purl2"},
			},
		}

		resultRepo.EXPECT().GetResults(nil).Return(results, nil)
		scanossSettingsRepo.EXPECT().GetSettings().Return(&entities.SettingsFile{
			Bom: entities.Bom{
				Include: []entities.ComponentFilter{
					{
						Purl: "pkg:npm/purl2",
					},
				},
			},
		})
		scanossSettingsRepo.EXPECT().GetDeclaredPurls().Return([]string{"pkg:github/purl2"})

		service := service.NewComponentServiceImpl(repo, scanossSettingsRepo, resultRepo, componentMapper)

		declaredComponents, err := service.GetDeclaredComponents()
		assert.NoError(t, err)
		assert.NotEmpty(t, declaredComponents)
		assert.Len(t, declaredComponents, 2)
		assert.Equal(t, declaredComponents, []entities.DeclaredComponent{
			{
				Purl: "pkg:npm/purl1",
				Name: "component1",
			},
			{
				Purl: "pkg:github/purl2",
				Name: "component2",
			}})

		resultRepo.AssertExpectations(t)
		scanossSettingsRepo.AssertExpectations(t)
		componentMapper.AssertExpectations(t)
	})
}
