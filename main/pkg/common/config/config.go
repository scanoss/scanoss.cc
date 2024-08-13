package config

import (
	"fmt"
	"integration-git/main/pkg/common/config/adapter"
	"integration-git/main/pkg/common/config/domain"
	"sync"
)

var (
	config     *domain.Config
	once       sync.Once
	configLock sync.Mutex
)

// LoadConfig reads the configuration and sets it as the singleton instance
func LoadConfig(filename string) (*domain.Config, error) {
	configLock.Lock()
	defer configLock.Unlock()

	once.Do(func() {
		cfgReader, _ := adapter.NewConfigServiceReaderFactory().Create(filename)
		cfg, err := cfgReader.ReadConfig(filename)
		if err != nil {
			config = nil
			return
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
