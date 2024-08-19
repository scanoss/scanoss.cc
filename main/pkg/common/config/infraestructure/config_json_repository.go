package infraestructure

import (
	"encoding/json"
	"errors"
	"integration-git/main/pkg/common/config/domain"
	"os"
)

var (
	ErrReadingFile       = errors.New("error reading file")
	ErrUnmarshallingFile = errors.New("error unmarshalling file file")
)

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
