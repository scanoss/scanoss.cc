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
const ROOT_FOLDER = "."
const DEFAULT_API_URL = "https://api.osskb.org"

var configPath = ""

// LoadConfig reads the configuration and sets it as the singleton instance
func LoadConfig(filename string) (*domain.Config, error) {
	configLock.Lock()
	defer configLock.Unlock()
	configPath = filename
	once.Do(func() {
		cfgReader, _ := adapter.NewConfigServiceReaderFactory().Create(filename)
		cfg, err := cfgReader.ReadConfig(filename)
		if err != nil {
			fmt.Println("Config file does not exist, please add the 'scanoss-lui-settings.json' in $HOME/.scanoss/")
			os.Exit(1)
		}

		config = &cfg
	})

	if config == nil {
		return nil, fmt.Errorf("failed to load config")
	}

	return config, nil
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

func GetConfigPath() string {
	return configPath
}

func GetDefaultApiURL() string {
	return DEFAULT_API_URL
}
