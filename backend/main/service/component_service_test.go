package service_test

// import (
// 	"slices"
// 	"testing"

// 	internalTest "github.com/scanoss/scanoss.lui/backend/main/internal"
// 	scanossBomEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
// 	scanossBomRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
// 	ssMockRepo "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository/mocks"
// 	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
// 	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
// 	mockRepo "github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository/mocks"
// 	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
// 	resultEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
// 	resultRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/result/repository"
// 	resultMockRepo "github.com/scanoss/scanoss.lui/backend/main/pkg/result/repository/mocks"
// 	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestInsertComponentFilterActions(t *testing.T) {
// 	cleanup := internalTest.InitializeTestEnvironment(t)
// 	defer cleanup()

// 	fr := utils.NewDefaultFileReader()
// 	repo := repository.NewJSONComponentRepository(fr)
// 	resultRepo := resultRepository.NewResultRepositoryJsonImpl(fr)
// 	scanossSettingsRepo := scanossBomRepository.NewScanossSettingsJsonRepository(fr)
// 	settingsFile, _ := scanossSettingsRepo.Read()
// 	scanossBomEntities.ScanossSettingsJson = &scanossBomEntities.ScanossSettings{
// 		SettingsFile: &settingsFile,
// 	}
// 	service := service.NewComponentService(repo, scanossSettingsRepo, resultRepo)

// 	tests := []struct {
// 		name string
// 		dto  entities.ComponentFilterDTO
// 	}{
// 		{
// 			name: "Include action",
// 			dto: entities.ComponentFilterDTO{
// 				Path:   "test/path1",
// 				Purl:   "pkg:purl1",
// 				Action: entities.Include,
// 			},
// 		},
// 		{
// 			name: "Remove action",
// 			dto: entities.ComponentFilterDTO{
// 				Path:   "test/path2",
// 				Purl:   "pkg:purl2",
// 				Action: entities.Remove,
// 			},
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			err := service.FilterComponents([]entities.ComponentFilterDTO{tc.dto})
// 			require.NoError(t, err)

// 			var filters = settingsFile.Bom.Include
// 			if tc.dto.Action == entities.Remove {
// 				filters = settingsFile.Bom.Remove
// 			}

// 			i := slices.IndexFunc(filters, func(cf scanossBomEntities.ComponentFilter) bool {
// 				return cf.Path == tc.dto.Path
// 			})

// 			filter := filters[i]

// 			require.Equal(t, tc.dto.Path, filter.Path)
// 			require.Equal(t, tc.dto.Purl, filter.Purl)
// 		})
// 	}
// }
// func TestGetDeclaredComponents(t *testing.T) {
// 	t.Run("WHEN Scanoss Settings is empty -> SHOULD return detected components from results file", func(t *testing.T) {
// 		repo := mockRepo.NewMockComponentRepository(t)
// 		resultRepo := resultMockRepo.NewMockResultRepository(t)
// 		scanossSettingsRepo := ssMockRepo.NewMockScanossSettingsRepository(t)
// 		results := []resultEntities.Result{
// 			{
// 				ComponentName: "component1",
// 				Path:          "path/to/file",
// 				MatchType:     "file",
// 				Purl:          &[]string{"pkg:npm/purl1"},
// 			},
// 			{
// 				ComponentName: "component2",
// 				Path:          "path/to/file2",
// 				MatchType:     "snippet",
// 				Purl:          &[]string{"pkg:github/purl2"},
// 			},
// 		}

// 		resultRepo.EXPECT().GetResults(nil).Return(results, nil)
// 		scanossSettingsRepo.EXPECT().GetDeclaredPurls().Return(nil)

// 		service := service.NewComponentService(repo, scanossSettingsRepo, resultRepo)

// 		declaredComponents, err := service.GetDeclaredComponents()
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, declaredComponents)
// 		assert.Len(t, declaredComponents, 2)
// 		assert.Equal(t, declaredComponents, []entities.DeclaredComponent{
// 			{
// 				Purl: "pkg:npm/purl1",
// 				Name: "component1",
// 			}, {
// 				Purl: "pkg:github/purl2",
// 				Name: "component2",
// 			},
// 		})

// 		resultRepo.AssertExpectations(t)
// 		scanossSettingsRepo.AssertExpectations(t)
// 	})

// 	t.Run("WHEN Scanoss Settings is NOT empty -> SHOULD return detected components from results file and declared components in scanoss settings", func(t *testing.T) {
// 		repo := mockRepo.NewMockComponentRepository(t)
// 		resultRepo := resultMockRepo.NewMockResultRepository(t)
// 		scanossSettingsRepo := ssMockRepo.NewMockScanossSettingsRepository(t)

// 		results := []resultEntities.Result{
// 			{
// 				ComponentName: "component1",
// 				Path:          "path/to/file",
// 				MatchType:     "file",
// 				Purl:          &[]string{"pkg:npm/purl1"},
// 			},
// 			{
// 				ComponentName: "component2",
// 				Path:          "path/to/file2",
// 				MatchType:     "snippet",
// 				Purl:          &[]string{"pkg:github/purl2"},
// 			},
// 		}

// 		resultRepo.EXPECT().GetResults(nil).Return(results, nil)
// 		scanossSettingsRepo.EXPECT().GetDeclaredPurls().Return([]string{"pkg:github/purl3"})

// 		service := service.NewComponentService(repo, scanossSettingsRepo, resultRepo)

// 		declaredComponents, err := service.GetDeclaredComponents()
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, declaredComponents)
// 		assert.Len(t, declaredComponents, 3)
// 		assert.Equal(t, declaredComponents, []entities.DeclaredComponent{
// 			{
// 				Purl: "pkg:npm/purl1",
// 				Name: "component1",
// 			},
// 			{
// 				Purl: "pkg:github/purl2",
// 				Name: "component2",
// 			}, {
// 				Purl: "pkg:github/purl3",
// 				Name: "",
// 			}})

// 		resultRepo.AssertExpectations(t)
// 		scanossSettingsRepo.AssertExpectations(t)
// 	})

// 	t.Run("WHEN there are repeated purls SHOULD NOT return duplicated values", func(t *testing.T) {
// 		repo := mockRepo.NewMockComponentRepository(t)
// 		resultRepo := resultMockRepo.NewMockResultRepository(t)
// 		scanossSettingsRepo := ssMockRepo.NewMockScanossSettingsRepository(t)

// 		results := []resultEntities.Result{
// 			{
// 				ComponentName: "component1",
// 				Path:          "path/to/file",
// 				MatchType:     "file",
// 				Purl:          &[]string{"pkg:npm/purl1"},
// 			},
// 			{
// 				ComponentName: "component2",
// 				Path:          "path/to/file2",
// 				MatchType:     "snippet",
// 				Purl:          &[]string{"pkg:github/purl2"},
// 			},
// 		}

// 		resultRepo.EXPECT().GetResults(nil).Return(results, nil)
// 		scanossSettingsRepo.EXPECT().GetDeclaredPurls().Return([]string{"pkg:github/purl2"})

// 		service := service.NewComponentService(repo, scanossSettingsRepo, resultRepo)

// 		declaredComponents, err := service.GetDeclaredComponents()
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, declaredComponents)
// 		assert.Len(t, declaredComponents, 2)
// 		assert.Equal(t, declaredComponents, []entities.DeclaredComponent{
// 			{
// 				Purl: "pkg:npm/purl1",
// 				Name: "component1",
// 			},
// 			{
// 				Purl: "pkg:github/purl2",
// 				Name: "component2",
// 			}})

// 		resultRepo.AssertExpectations(t)
// 		scanossSettingsRepo.AssertExpectations(t)
// 	})
// }
