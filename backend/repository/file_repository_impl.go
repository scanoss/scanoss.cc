package repository

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/scanoss/scanoss.lui/backend/entities"
	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/scanoss/scanoss.lui/internal/fetch"
)

type FileRepositoryImpl struct{}

func NewFileRepositoryImpl() FileRepository {
	return &FileRepositoryImpl{}
}

func (r *FileRepositoryImpl) ReadLocalFile(path string) (entities.File, error) {
	scanRootPath := config.Get().ScanRoot

	absolutePath := filepath.Join(scanRootPath, path)

	content, err := os.ReadFile(absolutePath)
	if err != nil {
		return entities.File{}, entities.ErrReadingFile
	}

	return *entities.NewFile(scanRootPath, path, content), nil
}

func (r *FileRepositoryImpl) ReadRemoteFileByMD5(path string, md5 string) (entities.File, error) {
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
		return entities.File{}, fmt.Errorf("failed to fetch file content: %w", err)
	}

	basePath, err := os.Getwd()
	if err != nil {
		return entities.File{}, fmt.Errorf("failed to get current working directory: %w", err)
	}

	return *entities.NewFile(basePath, path, []byte(body)), nil
}

func (r *FileRepositoryImpl) GetComponentByFilePath(filePath string) (entities.Component, error) {
	return entities.Component{}, nil
}
