package infrastructure

import (
	"errors"
	"os"
)

var (
	ErrReadingFile = errors.New("error reading file")
)

// fileRepository is a concrete implementation of the FileRepository interface.
type FileRepository struct{}

// NewFileRepository creates a new instance of fileRepository.
func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

// ReadFile reads the content of a file at the given path.
func (r *FileRepository) ReadLocalFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, ErrReadingFile
	}
	return data, nil
}

func (r *FileRepository) ReadRemoteFile(path string) ([]byte, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, ErrReadingFile
	}
	return data, nil
}
