package config

import (
	"fmt"
	"integration-git/main/pkg/common/config/adapter"
	"integration-git/main/pkg/common/config/domain"
	"os"
	"sync"
)

var (
	config     *domain.Config
	once       sync.Once
	configLock sync.Mutex
)

const GLOBAL_CONFIG_FILE_NAME = "scanoss-lui-settings.json"
const GLOBAL_CONFIG_FOLDER = "scanoss"
const SCANOSS_JSON = "scanoss.json"
const SCAN_SETTINGS_DEFAULT_LOCATION = ".scanoss/scanoss.json"
const ROOT_FOLDER = "."
const DEFAULT_API_URL = "https://api.osskb.org"

// Read reads the configuration from the specified file and sets it in the application.
//
// If the configuration file does not exist, the program will terminate with an error message.
//
// Parameters:
// - filename: The name of the configuration file to read.
//
// Returns a pointer to the Config struct containing the configuration, or an error if reading fails.
func Read(filename string) (*domain.Config, error) {
	cfgReader, _ := adapter.NewConfigServiceReaderFactory().Create(filename)
	cfg, err := cfgReader.ReadConfig(filename)
	if err != nil {
		fmt.Println("Config file does not exist, please add the 'scanoss-lui-settings.json' in $HOME/.scanoss/")
		os.Exit(1)
	}
	SetConfig(&cfg)

	return config, nil
}

func SetConfig(cfg *domain.Config) *domain.Config {
	configLock.Lock()
	defer configLock.Unlock()
	once.Do(func() {
		config = cfg
	})
	return config
}

// Get returns the singleton instance of the configuration
func Get() *domain.Config {
	return config
}

func GetDefaultGlobalFolder() string {
	return GLOBAL_CONFIG_FOLDER
}

func GetDefaultConfigFileName() string {
	return GLOBAL_CONFIG_FILE_NAME
}

func GetDefaultConfigLocation() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir + string(os.PathSeparator) + ROOT_FOLDER + GLOBAL_CONFIG_FOLDER + string(os.PathSeparator) + GLOBAL_CONFIG_FILE_NAME
}

func GetScanSettingDefaultLocation() string {
	workingDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return workingDir + string(os.PathSeparator) + SCAN_SETTINGS_DEFAULT_LOCATION
}

func GetScanossJsonDefaultLocation() string {
	workingDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return workingDir + string(os.PathSeparator) + SCANOSS_JSON
}

func GetDefaultApiURL() string {
	return DEFAULT_API_URL
}
