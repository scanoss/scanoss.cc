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
	DEFAULT_CONFIG_FILE_NAME      = "scanoss-lui-settings"
	DEFAULT_CONFIG_FILE_TYPE      = "json"
	ROOT_FOLDER                   = "."
	SCANOSS_HIDDEN_FOLDER         = ".scanoss"
	SCANOSS_PREMIUM_API_URL       = "https://api.scanoss.com"
)

type Config struct {
	ApiToken             string `json:"apitoken"`
	ApiUrl               string `json:"apiurl"`
	ResultFilePath       string `json:"resultfilepath,omitempty"`
	ScanRoot             string `json:"scanroot,omitempty"`
	ScanSettingsFilePath string `json:"scansettingsfilepath,omitempty"`
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
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

func GetDefaultConfigFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		return ""
	}

	return filepath.Join(homeDir, ROOT_FOLDER, SCANOSS_HIDDEN_FOLDER)
}
