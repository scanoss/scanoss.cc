package infrastructure

import (
	"errors"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/common/fetch"
	"integration-git/main/pkg/file/domain"
	"net/http"
	"os"
	"path"
)

var (
	ErrReadingFile     = errors.New("error reading file")
	ErrFetchingContent = errors.New("Error fetching remote file")
)

type Component struct {
	FileHash string `json:"file_hash"`
	File     string `json:"file"`
}

// fileRepository is a concrete implementation of the FileRepository interface.
type FileRepository struct{}

// NewFileRepository creates a new instance of fileRepository.
func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

// ReadFile reads the content of a file at the given path.
func (r *FileRepository) ReadLocalFile(filePath string) (domain.File, error) {
	currentPath := config.Get().ScanRoot

	absolutePath := path.Join(currentPath, filePath)

	content, err := os.ReadFile(absolutePath)
	if err != nil {
		return domain.File{}, ErrReadingFile
	}

	return *domain.NewFile(currentPath, filePath, content), nil
}

func (r *FileRepository) ReadRemoteFileByMD5(path string, md5 string) (domain.File, error) {
	baseURL := config.Get().ApiUrl
	token := config.Get().ApiToken

	url := fmt.Sprintf("%s/file_contents/%s", baseURL, md5)

	headers := make(map[string]string)
	if token != "" {
		headers["X-Session"] = token
	}

	options := fetch.Options{
		Method:  http.MethodGet,
		Headers: headers,
	}

	body, err := fetch.Text(url, options)
	if err != nil {
		return domain.File{}, fmt.Errorf("failed to fetch file content: %w", err)
	}

	basePath, err := os.Getwd()

	return *domain.NewFile(basePath, path, []byte(body)), nil

}
