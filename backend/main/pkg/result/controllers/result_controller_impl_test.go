package controllers_test

import (
	"testing"

	internal_test "github.com/scanoss/scanoss.lui/backend/main/internal"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
	mapperMocks "github.com/scanoss/scanoss.lui/backend/main/pkg/result/mappers/mocks"
	serviceMocks "github.com/scanoss/scanoss.lui/backend/main/pkg/result/service/mocks"
	"github.com/stretchr/testify/assert"
)

func init() {
	internal_test.InitValidatorForTests()
}

func TestGetAll(t *testing.T) {
	mockService := serviceMocks.NewMockResultService(t)
	mockMapper := mapperMocks.NewMockResultMapper(t)
	controller := controllers.NewResultController(mockService, mockMapper)

	t.Run("Valid request", func(t *testing.T) {
		dto := &entities.RequestResultDTO{
			MatchType: entities.File,
		}
		filter := entities.NewResultFilterFactory().Create(dto)
		results := []entities.Result{
			{
				Path:      "path/to/file",
				MatchType: "file",
				Purl:      &[]string{"pkg:npm/purl1"},
			},
			{
				Path:      "path/to/file2",
				MatchType: "snippet",
				Purl:      &[]string{"pkg:github/purl2"},
			},
		}
		expectedDTOs := []entities.ResultDTO{
			{
				Path:          "path/to/file",
				MatchType:     "file",
				WorkflowState: entities.Pending,
				FilterConfig: entities.FilterConfig{
					Action: entities.Include,
					Type:   entities.ByFile,
				},
			},
		}

		mockService.EXPECT().GetResults(filter).Return(results, nil)
		mockMapper.EXPECT().MapToResultDTOList(results).Return(expectedDTOs)

		actualDTOs, err := controller.GetAll(dto)

		assert.NoError(t, err)
		assert.Equal(t, expectedDTOs, actualDTOs)
		assert.Len(t, actualDTOs, 1)
		assert.Equal(t, "path/to/file", actualDTOs[0].Path)
		assert.Equal(t, entities.MatchType("file"), actualDTOs[0].MatchType)
		assert.Equal(t, entities.Pending, actualDTOs[0].WorkflowState)

		mockService.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

}
