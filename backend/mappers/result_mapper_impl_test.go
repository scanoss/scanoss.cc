//go:build unit

package mappers_test

import (
	"testing"

	"github.com/scanoss/scanoss.lui/backend/entities"
	mappers "github.com/scanoss/scanoss.lui/backend/mappers"
	internal_test "github.com/scanoss/scanoss.lui/internal"
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
				Purl:      &[]string{"pkg:purl-test"},
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
				Purl:      &[]string{"pkg:another-purl-test"},
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
				Purl:      &[]string{"pkg:removed-file"},
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

	settings := &entities.ScanossSettings{
		SettingsFile: &entities.SettingsFile{
			Bom: entities.Bom{
				Include: []entities.ComponentFilter{
					{
						Path: "path/to/file",
						Purl: "pkg:purl-test",
					},
				},
				Remove: []entities.ComponentFilter{{
					Purl: "pkg:purl-removed",
				}},
			},
		},
	}

	mapper := mappers.NewResultMapper(settings)

	for _, tc := range testCases {
		dto := mapper.MapToResultDTO(tc.result)
		assert.Equal(t, tc.expected.MatchType, dto.MatchType)
		assert.Equal(t, tc.expected.Path, dto.Path)
	}
}
