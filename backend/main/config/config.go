package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/scanoss/scanoss.lui/backend/main/utils"
)

var (
	ErrReadingFile       = errors.New("error reading file")
	ErrUnmarshallingFile = errors.New("error unmarshalling file")
)

const (
	DEFAULT_API_URL         = "https://api.osskb.org"
	GLOBAL_CONFIG_FILE_NAME = "scanoss-lui-settings.json"
	GLOBAL_CONFIG_FOLDER    = ".scanoss"
	ROOT_FOLDER             = "."
	SCNOSS_PREMIUM_API_URL  = "https://api.scanoss.com"
)

type Config struct {
	ApiToken             string `json:"apiToken"`
	ApiUrl               string `json:"apiUrl"`
	ResultFilePath       string `json:"resultFilePath,omitempty"`
	ScanRoot             string `json:"scanRoot,omitempty"`
	ScanSettingsFilePath string `json:"scanSettingsFilePath,omitempty"`
}

type ConfigModule struct {
	Config     *Config
	configPath string
	mutex      sync.RWMutex
}

var instance *ConfigModule
var once sync.Once

func GetInstance() *ConfigModule {
	once.Do(func() {
		instance = &ConfigModule{
			Config: &Config{
				ApiUrl: DEFAULT_API_URL,
			},
		}
	})
	return instance
}

func Get() *Config {
	return GetInstance().Config
}

func NewConfigModule(configPath string) *ConfigModule {
	GetInstance().configPath = configPath
	return GetInstance()
}

func (c *ConfigModule) Init() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := utils.FileExist(c.configPath)
	if err != nil {
		fmt.Println("Configuration file does not exists. Creating default config file...")
		err = c.createConfigFile()
		if err != nil {
			return err
		}
		if !utils.IsWritableFile(c.configPath) {
			return fmt.Errorf("file %s is not writable", c.configPath)
		}
		return nil
	}

	if !utils.IsWritableFile(c.configPath) {
		return fmt.Errorf("file %s is not writable", c.configPath)
	}

	return nil
}

func (c *ConfigModule) LoadConfig() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cfg, err := c.Read()
	if err != nil {
		return err
	}
	GetInstance().Config = &cfg

	return nil
}

func (c *ConfigModule) Read() (Config, error) {
	fileData, err := os.ReadFile(c.configPath)

	if err != nil {
		fmt.Printf("Unable to read configuration file %v\n", err)
		os.Exit(1)
	}

	// Marshal the default config into JSON
	defaultConfig, err := json.Marshal(getDefaultConfig())
	if err != nil {
		return Config{}, err
	}

	// Merge the JSON by unmarshalling the file's content into the default JSON
	var mergedData map[string]json.RawMessage
	if err := json.Unmarshal(defaultConfig, &mergedData); err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal(fileData, &mergedData); err != nil {
		return Config{}, err
	}

	// Marshal the merged JSON back into the config struct
	var cfg Config
	mergedBytes, err := json.Marshal(mergedData)
	if err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal(mergedBytes, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *ConfigModule) Save() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return utils.WriteJsonFile(c.configPath, c.Config)
}

func GetDefaultConfigLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		return GLOBAL_CONFIG_FILE_NAME
	}

	return filepath.Join(homeDir, ROOT_FOLDER, GLOBAL_CONFIG_FOLDER, GLOBAL_CONFIG_FILE_NAME)
}

func (c *ConfigModule) createConfigFile() error {
	// Get the directory path
	dirPath := filepath.Dir(c.configPath)
	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(dirPath, 0700)
		if err != nil {
			return err
		}
	}

	// Create the file with read/write permissions (owner and group only)
	file, err := os.OpenFile(c.configPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("Failed to create configuration file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	defaultConfig := getDefaultConfig()
	content := fmt.Sprintf(`{ "apiToken": "" , "apiUrl": "%s"}`, defaultConfig.ApiUrl)
	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func getDefaultConfig() *Config {
	workingDir, _ := os.Getwd()
	return &Config{
		ResultFilePath:       filepath.Join(workingDir, ".scanoss", "results.json"),
		ApiUrl:               DEFAULT_API_URL,
		ScanSettingsFilePath: filepath.Join(workingDir, "scanoss.json"),
	}
}
