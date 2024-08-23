package controllers_test

import (
	internal_test "integration-git/main/internal"
	"integration-git/main/pkg/common/config"
	entities2 "integration-git/main/pkg/common/scanoss_bom/application/entities"
	modules "integration-git/main/pkg/common/scanoss_bom/module"
	"integration-git/main/pkg/component/controllers"
	"integration-git/main/pkg/component/entities"
	"integration-git/main/pkg/component/repositories"
	"integration-git/main/pkg/component/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterComponent_Integration(t *testing.T) {
	mockedConfig, cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	config.SetConfig(mockedConfig)

	repo := repositories.NewComponentRepository()
	useCase := usecases.NewComponentUseCase(repo)
	controller := controllers.NewComponentController(useCase)

	modules.NewScanossBomModule().Init(&entities2.BomFile{})
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
			expectedError: repositories.ErrInvalidFilterAction,
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
