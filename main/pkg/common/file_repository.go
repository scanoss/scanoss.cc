package common

import (
	"errors"
	"os"
)

var (
	ErrReadingFile = errors.New("error reading file")
)

// fileRepository is a concrete implementation of the FileRepository interface.
type fileRepository struct{}

// NewFileRepository creates a new instance of fileRepository.
func NewFileRepository() FileRepository {
	return &fileRepository{}
}

// ReadFile reads the content of a file at the given path.
func (r *fileRepository) ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, ErrReadingFile
	}
	return data, nil
}
