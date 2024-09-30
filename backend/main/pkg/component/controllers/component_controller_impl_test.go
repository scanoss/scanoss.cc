package controllers_test

import (
	"fmt"
	"testing"

	internalTest "github.com/scanoss/scanoss.lui/backend/main/internal"
	scanossSettingsRepository "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/repository"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestFilterComponent_Integration(t *testing.T) {
	cleanup := internalTest.InitializeTestEnvironment(t)
	defer cleanup()

	fr := utils.NewDefaultFileReader()
	repo := repository.NewJSONComponentRepository(fr)
	ssRepo := scanossSettingsRepository.NewScanossSettingsJsonRepository(fr)
	service := service.NewComponentService(repo, ssRepo)
	controller := controllers.NewComponentController(service)

	dto := entities.ComponentFilterDTO{
		Path:   "test/path",
		Purl:   "test:purl",
		Usage:  "file",
		Action: entities.Include,
	}

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
			expectedError: fmt.Errorf("%s: %s", repository.ErrInvalidFilterAction, entities.FilterAction("unsupported")),
			dto: entities.ComponentFilterDTO{
				Path:   "test/path",
				Purl:   "pkg:purl.com/test",
				Usage:  "file",
				Action: entities.FilterAction("unsupported"),
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
