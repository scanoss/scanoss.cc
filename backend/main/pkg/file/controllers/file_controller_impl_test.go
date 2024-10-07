package controllers

import (
	"errors"
	"testing"

	componentEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	componentMocks "github.com/scanoss/scanoss.lui/backend/main/pkg/component/service/mocks"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/service/mocks"
	"github.com/stretchr/testify/assert"
)

func resetMocks(mockFileService *mocks.MockFileService, mockComponentService *componentMocks.MockComponentService) {
	mockFileService.ExpectedCalls = nil
	mockFileService.Calls = nil
	mockComponentService.ExpectedCalls = nil
	mockComponentService.Calls = nil
}

func TestFileController(t *testing.T) {
	mockFileService := mocks.NewMockFileService(t)
	mockComponentService := componentMocks.NewMockComponentService(t)
	controller := NewFileController(mockFileService, mockComponentService)

	filePath := "components/button.tsx"
	fileHash := "hash"
	fileContent := []byte("function main() {\n\tconsole.log('Hello, World!');\n}")
	file := entities.NewFile("", filePath, fileContent)

	t.Run("GetRemoteFile", func(t *testing.T) {
		mockComponentService.EXPECT().GetComponentByFilePath(filePath).Return(componentEntities.Component{
			FileHash: fileHash,
		}, nil)
		mockFileService.EXPECT().GetRemoteFileContent(filePath, fileHash).Return(*file, nil)

		response, err := controller.GetRemoteFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, "button.tsx", response.Name)
		assert.Equal(t, "components/button.tsx", response.Path)
		assert.Equal(t, string(fileContent), response.Content)

		mockComponentService.AssertExpectations(t)
		mockFileService.AssertExpectations(t)

		resetMocks(mockFileService, mockComponentService)
	})

	t.Run("ComponentServiceError", func(t *testing.T) {
		path := "components/button.tsx"
		mockComponentService.EXPECT().GetComponentByFilePath(path).Return(componentEntities.Component{}, errors.New("component service error"))

		result, err := controller.GetRemoteFile(path)

		assert.Error(t, err)
		assert.Equal(t, entities.FileDTO{}, result)

		mockComponentService.AssertExpectations(t)
		resetMocks(mockFileService, mockComponentService)
	})

	t.Run("GetLocalFile", func(t *testing.T) {
		path := "components/button.tsx"
		fileContent := []byte("function main() {\n\tconsole.log('Hello, World!');\n}")
		file := entities.NewFile("", path, fileContent)

		mockFileService.EXPECT().GetLocalFileContent(path).Return(*file, nil)

		result, err := controller.GetLocalFile(path)

		assert.NoError(t, err)
		assert.Equal(t, "button.tsx", result.Name)
		assert.Equal(t, "components/button.tsx", result.Path)
		assert.Equal(t, string(fileContent), result.Content)

		mockFileService.AssertExpectations(t)
	})
}
