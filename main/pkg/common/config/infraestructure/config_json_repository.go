package infraestructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"integration-git/main/pkg/common/config/domain"
	"os"
)

var (
	ErrReadingFile       = errors.New("error reading file")
	ErrUnmarshallingFile = errors.New("error unmarshalling file file")
)

// fileRepository is a concrete implementation of the FileRepository interface.
type ConfigJsonRepository struct{}

// NewFileRepository creates a new instance of fileRepository.
func NewConfigJsonRepository() *ConfigJsonRepository {
	return &ConfigJsonRepository{}
}

func (r *ConfigJsonRepository) Read(path string) (domain.Config, error) {
	fileData, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("Unable to read configuration file %v\n", err)
		os.Exit(1)
	}

	// Marshal the default config into JSON
	defaultConfig, err := json.Marshal(getDefaultConfigFile())
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

func getDefaultConfigFile() domain.Config {
	workingDir, _ := os.Getwd()
	var defaultConfigFile domain.Config = domain.Config{
		ScanRoot:             "",
		ResultFilePath:       workingDir + string(os.PathSeparator) + ".scanoss" + string(os.PathSeparator) + "results.json",
		ApiToken:             "",
		ApiUrl:               "https://api.osskb.org",
		ScanSettingsFilePath: workingDir + string(os.PathSeparator) + "scanoss.json",
	}

	return defaultConfigFile
}
