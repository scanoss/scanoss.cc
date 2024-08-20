package infraestructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"integration-git/main/pkg/common/config/domain"
	"integration-git/main/pkg/component/entities"
	"integration-git/main/pkg/utils"
	"os"
)

var (
	ErrReadingFile       = errors.New("error reading file")
	ErrUnmarshallingFile = errors.New("error unmarshalling file file")
)

const SCAN_SETTINGS_DEFAULT_LOCATION = ".scanoss/scanoss.json"

// fileRepository is a concrete implementation of the FileRepository interface.
type ConfigJsonRepository struct{}

var defaultConfigFile domain.Config = domain.Config{
	ScanRoot:       "",
	ResultFilePath: "./scanoss/results.json",
	ApiToken:       "",
	ApiUrl:         "",
}

// NewFileRepository creates a new instance of fileRepository.
func NewConfigJsonRepository() *ConfigJsonRepository {
	return &ConfigJsonRepository{}
}

func (r *ConfigJsonRepository) Read(path string) (domain.Config, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return domain.Config{}, ErrReadingFile
	}

	readScanSettingsFile(&defaultConfigFile)

	// Marshal the default config into JSON
	defaultConfig, err := json.Marshal(defaultConfigFile)
	if err != nil {
		return domain.Config{}, err
	}

	// Merge the JSON by unmarshalling the file's content into the default JSON
	var mergedData map[string]json.RawMessage
	if err := json.Unmarshal(defaultConfig, &mergedData); err != nil {
		return domain.Config{}, err
	}
	if err := json.Unmarshal(fileData, &mergedData); err != nil {
		return domain.Config{}, err
	}

	// Marshal the merged JSON back into the config struct
	var cfg domain.Config
	mergedBytes, err := json.Marshal(mergedData)
	if err != nil {
		return domain.Config{}, err
	}
	if err := json.Unmarshal(mergedBytes, &cfg); err != nil {
		return domain.Config{}, err
	}

	return cfg, nil
}

func readScanSettingsFile(cfg *domain.Config) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	scanSettingsFilePath := workingDir + "/" + SCAN_SETTINGS_DEFAULT_LOCATION
	_, err = os.Stat(scanSettingsFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Scan settings file not found in %s scanoss.json does not exist, creating now...", scanSettingsFilePath)
			initScanSettingsFile(scanSettingsFilePath)
		}
		return err
	}

	cfg.ScanSettingsFilePath = scanSettingsFilePath

	return nil
}

func initScanSettingsFile(path string) error {
	defaultScanSettings := entities.ScanSettingsFile{
		Bom: entities.Bom{
			Include: []entities.ComponentFilter{},
			Remove:  []entities.ComponentFilter{},
		},
	}

	err := utils.WriteJsonFile(path, defaultScanSettings)
	if err != nil {
		return err
	}

	return nil
}
