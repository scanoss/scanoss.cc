package service_test

// import (
// 	"testing"

// 	"github.com/scanoss/scanoss.lui/backend/main/config"
// 	"github.com/scanoss/scanoss.lui/backend/main/entities"
// 	internal_test "github.com/scanoss/scanoss.lui/backend/main/internal"
// 	"github.com/scanoss/scanoss.lui/backend/main/service"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestGetResults(t *testing.T) {
// 	cleanup := internal_test.InitializeTestEnvironment(t)
// 	defer cleanup()

// 	mu := internal_test.NewMockUtils()
// 	mu.On("ReadFile", config.Get().ResultFilePath).Return([]byte(`{"path/to/file": [{"ID": "file", "Purl": ["pkg:example/package"]}]}`), nil)

// 	filter := entitiesMocks.MockResultFilter{}
// 	filter.EXPECT().IsValid(mock.Anything).Return(true)

// 	mockRepo := mocks.NewMockResultService(t)
// 	mockRepo.EXPECT().GetResults(&filter).Return([]entities.Result{
// 		{
// 			Path:      "path/to/file",
// 			MatchType: "file",
// 			Purl:      &[]string{"pkg:example/package"},
// 		},
// 	}, nil)

// 	service := service.NewResultServiceImpl(mockRepo)

// 	results, err := service.GetResults(&filter)

// 	assert.NoError(t, err)
// 	assert.Len(t, results, 1)
// 	assert.Equal(t, results[0], entities.Result{
// 		Path:      "path/to/file",
// 		MatchType: "file",
// 		Purl:      &[]string{"pkg:example/package"},
// 	})
// }
