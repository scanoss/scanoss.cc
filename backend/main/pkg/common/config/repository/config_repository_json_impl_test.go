package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigJsonRepository(t *testing.T) {
	path := "config.json"
	repo := NewConfigJsonRepository(path)
	assert.Equal(t, path, repo.configPath)
}

func TestConfigJsonRepository_Init(t *testing.T) {
	path := "config.json"
	repo := NewConfigJsonRepository(path)

	err := repo.Init()
	assert.NoError(t, err)

	fileExists := utils.FileExist(path)
	assert.NoError(t, fileExists)

	os.Remove(path)
}

func TestConfigJsonRepository_Save(t *testing.T) {
	path := t.TempDir() + "config.json"
	repo := NewConfigJsonRepository(path)

	err := repo.Init()
	assert.NoError(t, err)

	config := &entities.Config{
		ApiUrl:   "https://api.scanoss.org",
		ApiToken: "test-token",
	}

	repo.Save(config)

	fileData, err := os.ReadFile(path)
	assert.NoError(t, err)

	savedData, err := utils.JSONParse[map[string]string](fileData)
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Equal(t, config.ApiUrl, savedData["apiUrl"])
	assert.Equal(t, config.ApiToken, savedData["apiToken"])

	os.Remove(path)
}

func TestConfigJsonRepository_Read(t *testing.T) {
	path := t.TempDir() + "config.json"
	repo := NewConfigJsonRepository(path)
	config := &entities.Config{
		ApiUrl:   "https://api.scanoss.org",
		ApiToken: "test-token",
	}

	err := repo.Init()
	assert.NoError(t, err)

	repo.Save(config)

	readConfig, err := repo.Read()

	assert.NoError(t, err)
	assert.Equal(t, config.ApiUrl, readConfig.ApiUrl)
	assert.Equal(t, config.ApiToken, readConfig.ApiToken)

	os.Remove(path)
}

func TestConfigJsonRepository_createConfigFile(t *testing.T) {
	path := "/tmp/test_config.json"
	repo := NewConfigJsonRepository(path)

	err := repo.createConfigFile()
	assert.NoError(t, err)

	fileExists := utils.FileExist(path)
	assert.NoError(t, fileExists)

	os.Remove(path)
}

func TestGetDefaultConfigFile(t *testing.T) {
	defaultConfig := getDefaultConfigFile()

	workingDir, _ := os.Getwd()
	expectedResultFilePath := filepath.Join(workingDir, ".scanoss", "results.json")
	expectedScanSettingsFilePath := filepath.Join(workingDir, "scanoss.json")

	assert.Equal(t, "", defaultConfig.ScanRoot)
	assert.Equal(t, expectedResultFilePath, defaultConfig.ResultFilePath)
	assert.Equal(t, "", defaultConfig.ApiToken)
	assert.Equal(t, "https://api.osskb.org", defaultConfig.ApiUrl)
	assert.Equal(t, expectedScanSettingsFilePath, defaultConfig.ScanSettingsFilePath)
}
