//go:build unit

package service_test

import (
	"errors"
	"testing"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/repository/mocks"
	"github.com/scanoss/scanoss.cc/backend/service"
	"github.com/stretchr/testify/assert"
)

func TestGetLocalFileContent(t *testing.T) {
	mockFileRepo := mocks.NewMockFileRepository(t)
	mockComponentRepo := mocks.NewMockComponentRepository(t)
	service := service.NewFileService(mockFileRepo, mockComponentRepo)

	expectedFile := entities.NewFile(
		"",
		"test.js",
		[]byte("function main() {\n\tconsole.log('Hello, World!');\n}"),
	)
	mockFileRepo.EXPECT().ReadLocalFile("test.js").Return(*expectedFile, nil)
	file, err := service.GetLocalFile("test.js")

	assert.NoError(t, err)
	assert.Equal(t, entities.FileDTO{
		Path:     "test.js",
		Name:     "test.js",
		Content:  "function main() {\n\tconsole.log('Hello, World!');\n}",
		Language: "javascript",
	}, file)
	mockFileRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}

func TestGetLocalFileContent_Error(t *testing.T) {
	mockFileRepo := mocks.NewMockFileRepository(t)
	mockComponentRepo := mocks.NewMockComponentRepository(t)
	service := service.NewFileService(mockFileRepo, mockComponentRepo)

	mockFileRepo.EXPECT().ReadLocalFile("test.js").Return(entities.File{}, errors.New("file not found"))

	file, err := service.GetLocalFile("test.js")

	assert.Error(t, err)
	assert.Equal(t, entities.FileDTO{}, file)
	mockFileRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}

func TestGetRemoteFileContent(t *testing.T) {
	mockFileRepo := mocks.NewMockFileRepository(t)
	mockComponentRepo := mocks.NewMockComponentRepository(t)
	service := service.NewFileService(mockFileRepo, mockComponentRepo)

	expectedFile := entities.NewFile(
		"",
		"remote.js",
		[]byte("function main() {\n\tconsole.log('Hello, World!');\n}"),
	)

	mockComponentRepo.EXPECT().FindByFilePath("remote.js").Return(entities.Component{FileHash: "test-md5"}, nil)
	mockFileRepo.EXPECT().ReadRemoteFileByMD5("remote.js", "test-md5").Return(*expectedFile, nil)

	file, err := service.GetRemoteFile("remote.js")

	assert.NoError(t, err)
	assert.Equal(t, entities.FileDTO{
		Path:     "remote.js",
		Name:     "remote.js",
		Content:  "function main() {\n\tconsole.log('Hello, World!');\n}",
		Language: "javascript",
	}, file)
	mockFileRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}

func TestGetRemoteFileContent_Error(t *testing.T) {
	mockFileRepo := mocks.NewMockFileRepository(t)
	mockComponentRepo := mocks.NewMockComponentRepository(t)
	service := service.NewFileService(mockFileRepo, mockComponentRepo)

	mockComponentRepo.EXPECT().FindByFilePath("remote.js").Return(entities.Component{}, errors.New("component not found"))

	file, err := service.GetRemoteFile("remote.js")

	assert.Error(t, err)
	assert.Equal(t, entities.FileDTO{}, file)
	mockFileRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}
