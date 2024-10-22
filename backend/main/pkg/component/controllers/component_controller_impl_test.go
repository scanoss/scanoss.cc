package controllers_test

import (
	"testing"

	"github.com/go-playground/validator"
	internal_test "github.com/scanoss/scanoss.lui/backend/main/internal"
	scanossSettingsEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/controllers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	mapperMocks "github.com/scanoss/scanoss.lui/backend/main/pkg/component/mappers/mocks"
	serviceMocks "github.com/scanoss/scanoss.lui/backend/main/pkg/component/service/mocks"
	"github.com/stretchr/testify/assert"
)

func init() {
	internal_test.InitValidatorForTests()
}

func TestFilterComponent_Integration(t *testing.T) {

	mockService := serviceMocks.NewMockComponentService(t)
	mockMapper := mapperMocks.NewMockComponentMapper(t)

	mockService.EXPECT().GetInitialFilters().Return(scanossSettingsEntities.InitialFilters{})

	controller := controllers.NewComponentController(mockService, mockMapper)

	t.Run("No errors", func(t *testing.T) {
		dto := entities.ComponentFilterDTO{
			Path:   "test/path",
			Purl:   "pkg:github/test",
			Usage:  "file",
			Action: entities.Include,
		}
		mockService.EXPECT().FilterComponents([]entities.ComponentFilterDTO{dto}).Return(nil)

		err := controller.FilterComponents([]entities.ComponentFilterDTO{dto})

		assert.NoError(t, err)
	})

	t.Run("Validation error", func(t *testing.T) {
		dto := entities.ComponentFilterDTO{
			Path:   "test/path",
			Purl:   "pkg:github/test2",
			Usage:  "file",
			Action: entities.FilterAction("unsupported"),
		}
		err := controller.FilterComponents([]entities.ComponentFilterDTO{dto})
		assert.Error(t, err)
		assert.IsType(t, err, validator.ValidationErrors{})
	})
}
