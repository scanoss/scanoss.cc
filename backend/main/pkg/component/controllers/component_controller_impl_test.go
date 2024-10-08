package controllers_test

import (
	"testing"

	"github.com/go-playground/validator"
	scanossSettingsEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	serviceMocks "github.com/scanoss/scanoss.lui/backend/main/pkg/component/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFilterComponent_Integration(t *testing.T) {
	mockService := serviceMocks.NewMockComponentService(t)
	mockService.EXPECT().GetInitialFilters().Return([]scanossSettingsEntities.ComponentFilter{}, []scanossSettingsEntities.ComponentFilter{})

	controller := controllers.NewComponentController(mockService)

	t.Run("No errors", func(t *testing.T) {
		dto := entities.ComponentFilterDTO{
			Path:   "test/path",
			Purl:   "test:purl",
			Usage:  "file",
			Action: entities.Include,
		}
		mockService.EXPECT().FilterComponent(dto).Return(nil)

		err := controller.FilterComponent(dto)

		assert.NoError(t, err)
	})

	t.Run("Validation error", func(t *testing.T) {
		dto := entities.ComponentFilterDTO{
			Path:   "test/path",
			Purl:   "pkg:purl.com/test",
			Usage:  "file",
			Action: entities.FilterAction("unsupported"),
		}
		err := controller.FilterComponent(dto)
		assert.Error(t, err)
		assert.IsType(t, err, validator.ValidationErrors{})
	})
}
