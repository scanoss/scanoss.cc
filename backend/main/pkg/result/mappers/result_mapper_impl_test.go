package mappers

import (
	"testing"

	internal_test "github.com/scanoss/scanoss.lui/backend/main/internal"
	scanossBomEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
	"github.com/stretchr/testify/assert"
)

func TestMapToResultDTO(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	testCases := []struct {
		result   entities.Result
		expected entities.ResultDTO
	}{
		{
			result: entities.Result{
				MatchType: "file",
				Path:      "path/to/file",
				Purl:      []string{"pkg:purl-test"},
			},
			expected: entities.ResultDTO{
				Path:          "path/to/file",
				MatchType:     "file",
				WorkflowState: entities.Completed,
				FilterConfig: entities.FilterConfig{
					Action: entities.Include,
					Type:   entities.ByFile,
				},
			},
		},
		{
			result: entities.Result{
				Path:      "another/path/to/file",
				MatchType: "file",
				Purl:      []string{"pkg:another-purl-test"},
			},
			expected: entities.ResultDTO{
				Path:          "another/path/to/file",
				MatchType:     "file",
				WorkflowState: entities.Pending,
				FilterConfig:  entities.FilterConfig{},
			},
		},
		{
			result: entities.Result{
				Path:      "path/to/removed/file",
				MatchType: "snippet",
				Purl:      []string{"pkg:removed-file"},
			},
			expected: entities.ResultDTO{
				Path:          "path/to/removed/file",
				MatchType:     "snippet",
				WorkflowState: entities.Completed,
				FilterConfig: entities.FilterConfig{
					Action: entities.Remove,
					Type:   entities.ByFile,
				},
			},
		},
	}

	settings := &scanossBomEntities.ScanossSettings{
		SettingsFile: &scanossBomEntities.SettingsFile{
			Bom: scanossBomEntities.Bom{
				Include: []scanossBomEntities.ComponentFilter{
					{
						Path: "path/to/file",
						Purl: "pkg:purl-test",
					},
				},
				Remove: []scanossBomEntities.ComponentFilter{{
					Purl: "pkg:purl-removed",
				}},
			},
		},
	}

	mapper := NewResultMapper(settings)

	for _, tc := range testCases {
		dto := mapper.MapToResultDTO(tc.result)
		assert.Equal(t, tc.expected.MatchType, dto.MatchType)
		assert.Equal(t, tc.expected.Path, dto.Path)
	}
}
