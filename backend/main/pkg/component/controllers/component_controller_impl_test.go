package controllers_test

import (
	"testing"

	internal_test "github.com/scanoss/scanoss.lui/backend/main/internal"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	"github.com/stretchr/testify/assert"
)

func TestFilterComponent_Integration(t *testing.T) {
	_, cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	repo := repository.NewComponentRepository()
	useCase := service.NewComponentUseCase(repo)
	controller := controllers.NewComponentController(useCase)

	dto := entities.ComponentFilterDTO{
		Path:    "test/path",
		Purl:    "test:purl",
		Usage:   "file",
		Version: "1.0",
		Action:  entities.Include,
	}

	// Create multiple test cases, one for no errors, and one with errors (for example sending wrong action)

	var tests = []struct {
		name          string
		expectedError error
		dto           entities.ComponentFilterDTO
	}{
		{
			name: "No errors",
			dto:  dto,
		},
		{
			name:          "Wrong action",
			expectedError: repository.ErrInvalidFilterAction,
			dto: entities.ComponentFilterDTO{
				Path:    "test/path",
				Purl:    "pkg:purl.com/test",
				Usage:   "file",
				Version: "1.0.3",
				Action:  entities.FilterAction("unsupported"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := controller.FilterComponent(tc.dto)
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				return
			}

			assert.NoError(t, err)
		})
	}

}
