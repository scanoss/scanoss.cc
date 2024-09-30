package controllers_test

import (
	"testing"

	internal_test "github.com/scanoss/scanoss.lui/backend/main/internal"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
	mappersMocks "github.com/scanoss/scanoss.lui/backend/main/pkg/result/mappers/mocks"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllResults(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	mockService := mocks.NewResultService(t)
	mockMapper := mappersMocks.NewResultMapper(t)

	mockService.EXPECT().GetResults(nil).Return([]entities.Result{
		{
			Path:      "path/to/file",
			MatchType: "file",
			Purl:      []string{"pkg:npm/purl1"},
		},
		{
			Path:      "path/to/file2",
			MatchType: "snippet",
			Purl:      []string{"pkg:github/purl2"},
		},
	}, nil)

	mockMapper.EXPECT().MapToResultDTOList(mock.Anything).Return([]entities.ResultDTO{
		{
			Path:          "path/to/file",
			MatchType:     "file",
			WorkflowState: entities.Pending,
			FilterConfig: entities.FilterConfig{
				Action: entities.Include,
				Type:   entities.ByFile,
			},
		},
	})

	c := controllers.NewResultController(mockService, mockMapper)
	results, err := c.GetAll(&entities.RequestResultDTO{})

	assert.NoError(t, err)
	assert.Len(t, results, 1)
}
