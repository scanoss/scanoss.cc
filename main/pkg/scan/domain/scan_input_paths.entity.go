package domain

import (
	"errors"
	"path/filepath"
)

var (
	ErrPathIsAbsolute = errors.New("path is absolute")
)

type ScanInputPaths struct {
	basePath string
	paths    []string
}

func NewScanInputPaths(rootPath string, files []string) (ScanInputPaths, error) {

	for _, file := range files {
		if filepath.IsAbs(file) {
			return ScanInputPaths{}, ErrPathIsAbsolute
		}
	}

	return ScanInputPaths{
		rootPath,
		files,
	}, nil
}

func (sif *ScanInputPaths) GetAbsolutePaths() []string {
	var f []string
	for _, file := range sif.paths {
		f = append(f, filepath.Join(sif.basePath, file))
	}
	return f
}

func (sif *ScanInputPaths) GetRelativePaths() []string {
	return sif.paths
}

func (sif *ScanInputPaths) GetBasePath() string {
	return sif.basePath
}
