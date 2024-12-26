//go:build unit

package service_test

import (
	"testing"

	"github.com/scanoss/scanoss.lui/backend/entities"
	mapperMocks "github.com/scanoss/scanoss.lui/backend/mappers/mocks"
	repoMocks "github.com/scanoss/scanoss.lui/backend/repository/mocks"
	"github.com/scanoss/scanoss.lui/backend/service"
	internal_test "github.com/scanoss/scanoss.lui/internal"
	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetResults(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	mu := internal_test.NewMockUtils()
	mu.On("ReadFile", config.GetInstance().ResultFilePath).Return([]byte(`{"path/to/file": [{"ID": "file", "Purl": ["pkg:example/package"]}]}`), nil)

	mockRepo := repoMocks.NewMockResultRepository(t)
	resultMapper := mapperMocks.NewMockResultMapper(t)
	mockRepo.EXPECT().GetResults(mock.AnythingOfType("*entities.ResultFilterAND")).Return([]entities.Result{
		{
			Path:      "path/to/file",
			MatchType: "file",
			Purl:      &[]string{"pkg:example/package"},
		},
	}, nil)

	resultMapper.EXPECT().MapToResultDTOList([]entities.Result{
		{
			Path:      "path/to/file",
			MatchType: "file",
			Purl:      &[]string{"pkg:example/package"},
		},
	}).Return([]entities.ResultDTO{{
		Path:         "path/to/file",
		MatchType:    "file",
		DetectedPurl: "pkg:example/package",
	}})

	service := service.NewResultServiceImpl(mockRepo, resultMapper)

	dto := &entities.RequestResultDTO{}

	results, err := service.GetAll(dto)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, results[0], entities.ResultDTO{
		Path:             "path/to/file",
		MatchType:        "file",
		ConcludedPurl:    "",
		WorkflowState:    "",
		FilterConfig:     entities.FilterConfig{},
		Comment:          "",
		DetectedPurl:     "pkg:example/package",
		ConcludedPurlUrl: "",
		ConcludedName:    "",
	})

	mockRepo.AssertExpectations(t)
	resultMapper.AssertExpectations(t)
}
