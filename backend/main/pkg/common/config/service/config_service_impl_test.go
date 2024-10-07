package service

import (
	"errors"
	"testing"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestReadConfig_Success(t *testing.T) {
	mockRepo := mocks.NewMockConfigRepository(t)
	expectedConfig := entities.Config{
		ApiToken:       "test-api-token",
		ApiUrl:         "http://api.scanoss.com/v1",
		ResultFilePath: ".scanoss/results.json",
		ScanRoot:       "frontend/src",
	}

	mockRepo.EXPECT().Read().Return(expectedConfig, nil)

	service := NewConfigService(mockRepo)

	config, err := service.ReadConfig()

	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)
	mockRepo.AssertExpectations(t)
}

func TestReadConfig_Error(t *testing.T) {
	mockRepo := mocks.NewMockConfigRepository(t)
	expectedError := errors.New("read error")

	mockRepo.EXPECT().Read().Return(entities.Config{}, expectedError)

	service := NewConfigService(mockRepo)

	config, err := service.ReadConfig()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, entities.Config{}, config)
	mockRepo.AssertExpectations(t)
}
