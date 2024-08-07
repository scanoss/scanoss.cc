package adapter

import (
	"fmt"
	"integration-git/main/pkg/common/config/application/services"
	"integration-git/main/pkg/common/config/application/use_cases"
	"integration-git/main/pkg/common/config/infraestructure"
	"path/filepath"
	"strings"
)

type ConfigFactory struct {
	s services.ConfigService
}

func NewConfigServiceReaderFactory() *ConfigFactory {
	return &ConfigFactory{}
}

func (f ConfigFactory) Create(path string) (use_cases.ReadConfigFileUseCase, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return use_cases.NewReadJsonConfigFileUseCase(services.NewConfigService(infraestructure.NewConfigJsonRepository())), nil
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
}
