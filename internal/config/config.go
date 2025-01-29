// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	DEFAULT_CONFIG_FILE_NAME      = "scanoss-cc-settings"
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
	Debug                bool   `json:"debug,omitempty"`
	mu                   sync.RWMutex
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func (c *Config) getDefaultResultFilePath(scanRoot string) string {
	return filepath.Join(scanRoot, SCANOSS_HIDDEN_FOLDER, DEFAULT_RESULTS_FILE)
}

func (c *Config) getDefaultScanSettingsFilePath(scanRoot string) string {
	return filepath.Join(scanRoot, DEFAULT_SCANOSS_SETTINGS_FILE)
}

func (c *Config) GetApiToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ApiToken
}

func (c *Config) GetApiUrl() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ApiUrl
}

func (c *Config) GetResultFilePath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ResultFilePath
}

func (c *Config) GetScanRoot() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ScanRoot
}

func (c *Config) GetScanSettingsFilePath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ScanSettingsFilePath
}

func (c *Config) SetApiToken(token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ApiToken = token
	viper.Set("apitoken", token)
	return viper.WriteConfig()
}

func (c *Config) SetApiUrl(url string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ApiUrl = url
	viper.Set("apiurl", url)
	return viper.WriteConfig()
}

func (c *Config) SetResultFilePath(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ResultFilePath = path
}

func (c *Config) SetScanRoot(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ScanRoot = path
	c.ResultFilePath = c.getDefaultResultFilePath(path)
	c.ScanSettingsFilePath = c.getDefaultScanSettingsFilePath(path)
}

func (c *Config) SetScanSettingsFilePath(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ScanSettingsFilePath = path
}

func (c *Config) GetDefaultConfigFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("Error getting user home directory")
		return ""
	}

	return filepath.Join(homeDir, ROOT_FOLDER, SCANOSS_HIDDEN_FOLDER)
}

func (c *Config) setupLogger(debug bool) error {
	logsDir := filepath.Join(c.GetDefaultConfigFolder(), "logs")
	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating logs directory: %w", err)
	}

	logFileName := fmt.Sprintf("scanoss-cc-%s.log", time.Now().Format(time.DateOnly))
	logFile, err := os.OpenFile(filepath.Join(logsDir, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)

	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return nil
}

func (c *Config) InitializeConfig(cfgFile, scanRoot, apiKey, apiUrl, inputFile, scanossSettingsFilePath string, originalWorkDir string, debug bool) error {
	if err := c.setupLogger(debug); err != nil {
		return fmt.Errorf("error setting up logger: %w", err)
	}

	if cfgFile != "" {
		absCfgFile, _ := filepath.Abs(cfgFile)
		log.Debug().Msgf("Using config file: %s", absCfgFile)

		viper.SetConfigFile(absCfgFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal().Err(err).Msgf("Error reading config file %v", err.Error())
		}
	} else {
		viper.SetConfigName(DEFAULT_CONFIG_FILE_NAME)
		viper.SetConfigType(DEFAULT_CONFIG_FILE_TYPE)
		viper.AddConfigPath(c.GetDefaultConfigFolder())

		// Default values
		viper.SetDefault("apiurl", DEFAULT_API_URL)
		viper.SetDefault("apitoken", "")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				defaultConfigDir := c.GetDefaultConfigFolder()
				if err := os.MkdirAll(defaultConfigDir, os.ModePerm); err != nil {
					return fmt.Errorf("error creating config directory: %w", err)
				}
				if err := viper.SafeWriteConfig(); err != nil {
					return fmt.Errorf("error creating config file: %w", err)
				}
				log.Debug().Msgf("Created default config file: %s", viper.ConfigFileUsed())
			} else {
				return fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	c.SetApiToken(viper.GetString("apitoken"))
	c.SetApiUrl(viper.GetString("apiurl"))

	c.Debug = debug

	// Override with command line flags
	if apiKey != "" {
		if err := c.SetApiToken(apiKey); err != nil {
			return fmt.Errorf("error saving API token: %w", err)
		}
	}
	if apiUrl != "" {
		if err := c.SetApiUrl(apiUrl); err != nil {
			return fmt.Errorf("error saving API URL: %w", err)
		}
	}
	if scanRoot != "" {
		c.SetScanRoot(scanRoot)
	}
	if inputFile != "" {
		c.SetResultFilePath(inputFile)
	}
	if scanossSettingsFilePath != "" {
		c.SetScanSettingsFilePath(scanossSettingsFilePath)
	}

	// Set default values if not set via config file or command line args
	if c.ScanRoot == "" {
		if originalWorkDir != "" {
			c.SetScanRoot(originalWorkDir)
		} else {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("error getting working directory: %w", err)
			}
			c.SetScanRoot(wd)
		}
	}
	if c.ResultFilePath == "" {
		c.SetResultFilePath(c.getDefaultResultFilePath(originalWorkDir))
	}
	if c.ScanSettingsFilePath == "" {
		c.SetScanSettingsFilePath(c.getDefaultScanSettingsFilePath(originalWorkDir))
	}

	return nil
}
