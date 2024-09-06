package repository

import (
	"fmt"
	"path/filepath"
	"strings"
)

func Create(path string) (ConfigRepository, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return NewConfigJsonRepository(path), nil
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}
}
