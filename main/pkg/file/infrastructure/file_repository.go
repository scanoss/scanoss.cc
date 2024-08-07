package infrastructure

import (
	"errors"
	"fmt"
	"integration-git/main/pkg/file/domain"
	"os"
	"path"
	"path/filepath"
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
func (r *FileRepository) ReadLocalFile(filePath string) (domain.File, error) {
	currentPath, err := os.Getwd()
	absolutePath := path.Join(currentPath, filePath)
	fmt.Printf("Absolute PAth: %s", absolutePath)
	// Get the file name from the path
	fileName := filepath.Base(filePath)

	data, err := os.ReadFile(absolutePath)
	if err != nil {
		return domain.File{}, ErrReadingFile
	}
	// Convert bytes to string
	content := string(data)

	file := domain.NewFile()
	file.SetContent(content)
	file.SetPath(filePath)
	file.SetName(fileName)
	return *file, nil
}

func (r *FileRepository) ReadRemoteFile(path string) (domain.File, error) {
	/*	apiUrl := config.GetConfig().Scanoss.ApiUrl
		apiToken := config.GetConfig().Scanoss.ApiToken

		data, err := os.ReadFile(path)
		if err != nil {
			return domain.File{}, ErrReadingFile
		}*/
	return domain.File{}, nil
}
