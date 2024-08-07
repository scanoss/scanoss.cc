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

// NewFileRepository creates a new instance of fileRepository.
func NewConfigJsonRepository() *ConfigJsonRepository {
	return &ConfigJsonRepository{}
}

func (r *ConfigJsonRepository) Read(path string) (domain.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.Config{}, ErrReadingFile
	}
	var cfg domain.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return domain.Config{}, ErrUnmarshallingFile
	}
	return cfg, nil
}
