//go:build unit

package repository_test

import (
	"testing"

	"github.com/scanoss/scanoss.lui/backend/entities"
	"github.com/scanoss/scanoss.lui/backend/entities/mocks"
	"github.com/scanoss/scanoss.lui/backend/repository"
	internal_test "github.com/scanoss/scanoss.lui/internal"
	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetResults(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	t.Run("No filter", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().ResultFilePath).Return([]byte(`{"path/to/file": [{"ID": "file", "Purl": ["pkg:example/package"]}]}`), nil)

		repo := repository.NewResultRepositoryJsonImpl(mu)
		results, err := repo.GetResults(nil)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "path/to/file", results[0].Path)
		assert.Equal(t, "file", results[0].MatchType)
		assert.Equal(t, &[]string{"pkg:example/package"}, results[0].Purl)
	})

	t.Run("With filter", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().ResultFilePath).Return([]byte(`{"path/to/file": [{"ID": "file", "Purl": ["pkg:example/package"]}]}`), nil)

		filter := mocks.MockResultFilter{}
		filter.EXPECT().IsValid(mock.Anything).Return(true)

		repo := repository.NewResultRepositoryJsonImpl(mu)
		results, err := repo.GetResults(&filter)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "path/to/file", results[0].Path)
		assert.Equal(t, "file", results[0].MatchType)
		assert.Equal(t, []string{"pkg:example/package"}, (*results[0].Purl))
	})

	t.Run("Read file error", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().ResultFilePath).Return([]byte{}, entities.ErrReadingResultFile)

		repo := repository.NewResultRepositoryJsonImpl(mu)
		results, err := repo.GetResults(nil)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrReadingResultFile, err)
		assert.Len(t, results, 0)
	})

	t.Run("Invalid json", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().ResultFilePath).Return([]byte(`invalid json`), nil)

		repo := repository.NewResultRepositoryJsonImpl(mu)
		results, err := repo.GetResults(nil)

		assert.Error(t, err)
		assert.Len(t, results, 0)
	})

	t.Run("Filter no match", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().ResultFilePath).Return([]byte(`{"path/to/file": [{"ID": "file", "Purl": ["pkg:example/package"]}]}`), nil)

		filter := mocks.MockResultFilter{}
		filter.EXPECT().IsValid(mock.Anything).Return(false)

		repo := repository.NewResultRepositoryJsonImpl(mu)
		results, err := repo.GetResults(&filter)

		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})
}
