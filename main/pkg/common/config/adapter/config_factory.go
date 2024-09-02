package adapter

import (
	"fmt"
	"integration-git/main/pkg/common/config/infraestructure"
	"path/filepath"
	"strings"
)

func Create(path string) (infraestructure.ConfigRepository, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return infraestructure.NewConfigJsonRepository(path), nil
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
}
