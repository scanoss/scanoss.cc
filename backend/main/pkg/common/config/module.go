package config

import (
	"os"
	"sync"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config/repository"
)

var (
	configModule *ConfigModule
	once         sync.Once
	configLock   sync.Mutex
)

type ConfigModule struct {
	Config *entities.Config
	repo   repository.ConfigRepository
	path   string
}

func NewConfigModule(path string) *ConfigModule {
	return &ConfigModule{
		path: path,
	}
}

const GLOBAL_CONFIG_FILE_NAME = "scanoss-lui-settings.json"
const GLOBAL_CONFIG_FOLDER = "scanoss"
const ROOT_FOLDER = "."
const DEFAULT_API_URL = "https://api.osskb.org"
const SCNOSS_PREMIUM_API_URL = "https://api.scanoss.com"

var configPath = ""

func (m *ConfigModule) Init() error {
	configLock.Lock()
	defer configLock.Unlock()
	once.Do(func() {
		repo, _ := repository.Create(m.path)
		m.repo = repo
		err := repo.Init()
		if err != nil {
			fmt.Println("Config file does not exist or cannot be created, please add the 'scanoss-lui-settings.json' in $HOME/.scanoss/")
			os.Exit(1)
		}

		configModule = m
	})

	return nil

}

// LoadConfig reads the configuration and sets it as the singleton instance
func (m *ConfigModule) LoadConfig() error {
	cfg, err := m.repo.Read()
	if err != nil {
		return err
	}
	m.Config = &cfg
	return nil
}

// Get returns the singleton instance of the configuration
func Get() *entities.Config {
	return configModule.Config
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

// GetInstance returns the singleton instance of ConfigModule
func GetInstance() *ConfigModule {
	if configModule == nil {
		panic("ConfigModule not initialized. Call LoadConfig first.")
	}
	return configModule
}

func (m *ConfigModule) Save() {
	m.repo.Save(m.Config)
}

func (m *ConfigModule) GetConfigPath() string {
	return m.path
}

func GetDefaultApiURL() string {
	return DEFAULT_API_URL
}