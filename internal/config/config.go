package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrReadingFile       = errors.New("error reading file")
	ErrUnmarshallingFile = errors.New("error unmarshalling file")
)

const (
	DEFAULT_API_URL               = "https://api.osskb.org"
	DEFAULT_RESULTS_FILE          = "results.json"
	DEFAULT_SCANOSS_SETTINGS_FILE = "scanoss.json"
	GLOBAL_CONFIG_FILE_NAME       = "scanoss-lui-settings.json"
	ROOT_FOLDER                   = "."
	SCANOSS_HIDDEN_FOLDER         = ".scanoss"
	SCANOSS_PREMIUM_API_URL       = "https://api.scanoss.com"
)

type Config struct {
	ApiToken             string `json:"apiToken"`
	ApiUrl               string `json:"apiUrl"`
	ResultFilePath       string `json:"resultFilePath,omitempty"`
	ScanRoot             string `json:"scanRoot,omitempty"`
	ScanSettingsFilePath string `json:"scanSettingsFilePath,omitempty"`
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{
			ApiUrl: DEFAULT_API_URL,
		}
	})
	return instance
}

func Get() *Config {
	return GetInstance()
}

func GetDefaultResultFilePath() string {
	workingDir, _ := os.Getwd()
	return filepath.Join(workingDir, SCANOSS_HIDDEN_FOLDER, DEFAULT_RESULTS_FILE)
}

func GetDefaultScanSettingsFilePath() string {
	workingDir, _ := os.Getwd()
	return filepath.Join(workingDir, DEFAULT_SCANOSS_SETTINGS_FILE)
}

func GetDefaultConfigLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		return GLOBAL_CONFIG_FILE_NAME
	}

	return filepath.Join(homeDir, ROOT_FOLDER, SCANOSS_HIDDEN_FOLDER, GLOBAL_CONFIG_FILE_NAME)
}
