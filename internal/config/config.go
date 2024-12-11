package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
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

func (c *Config) SetApiToken(token string) error {
	c.ApiToken = token
	viper.Set("apitoken", token)
	return viper.WriteConfig()
}

func (c *Config) SetApiUrl(url string) error {
	c.ApiUrl = url
	viper.Set("apiurl", url)
	return viper.WriteConfig()
}

func (c *Config) SetResultFilePath(path string) error {
	c.ResultFilePath = path
	viper.Set("resultfilepath", path)
	return viper.WriteConfig()
}

func (c *Config) SetScanRoot(path string) error {
	c.ScanRoot = path
	viper.Set("scanroot", path)
	return viper.WriteConfig()
}

func (c *Config) SetScanSettingsFilePath(path string) error {
	c.ScanSettingsFilePath = path
	viper.Set("scansettingsfilepath", path)
	return viper.WriteConfig()
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

func InitializeConfig(cfgFile, scanRoot, apiKey, apiUrl, inputFile string) error {
	if cfgFile != "" {
		absCfgFile, _ := filepath.Abs(cfgFile)
		fmt.Println("Using config file:", absCfgFile)

		viper.SetConfigFile(absCfgFile)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		viper.SetConfigName(DEFAULT_CONFIG_FILE_NAME)
		viper.SetConfigType(DEFAULT_CONFIG_FILE_TYPE)
		viper.AddConfigPath(GetDefaultConfigFolder())

		// Default values
		viper.SetDefault("apiUrl", DEFAULT_API_URL)
		viper.SetDefault("apiToken", "")
		viper.SetDefault("resultFilePath", GetDefaultResultFilePath())
		viper.SetDefault("scanRoot", "")
		viper.SetDefault("scanSettingsFilePath", GetDefaultScanSettingsFilePath())

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				defaultConfigDir := GetDefaultConfigFolder()
				if err := os.MkdirAll(defaultConfigDir, 0755); err != nil {
					return fmt.Errorf("error creating config directory: %w", err)
				}
				if err := viper.SafeWriteConfig(); err != nil {
					return fmt.Errorf("error creating config file: %w", err)
				}
				fmt.Println("Created default config file:", viper.ConfigFileUsed())
			} else {
				return fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	cfg := GetInstance()

	// Override with command line flags
	if scanRoot != "" {
		if err := cfg.SetScanRoot(scanRoot); err != nil {
			return fmt.Errorf("error saving scan root: %w", err)
		}
	}
	if apiKey != "" {
		if err := cfg.SetApiToken(apiKey); err != nil {
			return fmt.Errorf("error saving API token: %w", err)
		}
	}
	if apiUrl != "" {
		if err := cfg.SetApiUrl(apiUrl); err != nil {
			return fmt.Errorf("error saving API URL: %w", err)
		}
	}
	if inputFile != "" {
		if err := cfg.SetResultFilePath(inputFile); err != nil {
			return fmt.Errorf("error saving result file path: %w", err)
		}
	}

	// Load config values from viper if not set by flags
	if cfg.ApiToken == "" {
		cfg.ApiToken = viper.GetString("apiToken")
	}
	if cfg.ApiUrl == "" {
		cfg.ApiUrl = viper.GetString("apiUrl")
	}
	if cfg.ResultFilePath == "" {
		cfg.ResultFilePath = viper.GetString("resultFilePath")
	}
	if cfg.ScanRoot == "" {
		cfg.ScanRoot = viper.GetString("scanRoot")
	}
	if cfg.ScanSettingsFilePath == "" {
		cfg.ScanSettingsFilePath = viper.GetString("scanSettingsFilePath")
	}

	return nil
}
