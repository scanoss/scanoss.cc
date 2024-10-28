package repository

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/scanoss/scanoss.lui/backend/main/config"
	"github.com/scanoss/scanoss.lui/backend/main/entities"
	"github.com/scanoss/scanoss.lui/backend/main/utils"
)

type FileRepositoryImpl struct{}

func NewFileRepositoryImpl() FileRepository {
	return &FileRepositoryImpl{}
}

func (r *FileRepositoryImpl) ReadLocalFile(filePath string) (entities.File, error) {
	currentPath := config.Get().ScanRoot

	absolutePath := path.Join(currentPath, filePath)

	content, err := os.ReadFile(absolutePath)
	if err != nil {
		return entities.File{}, entities.ErrReadingFile
	}

	return *entities.NewFile(currentPath, filePath, content), nil
}

func (r *FileRepositoryImpl) ReadRemoteFileByMD5(path string, md5 string) (entities.File, error) {
	baseURL := config.Get().ApiUrl
	token := config.Get().ApiToken

	url := fmt.Sprintf("%s/file_contents/%s", baseURL, md5)

	headers := make(map[string]string)
	if token != "" {
		headers["X-Session"] = token
	}

	options := utils.Options{
		Method:  http.MethodGet,
		Headers: headers,
	}

	body, err := utils.Text(url, options)
	if err != nil {
		return entities.File{}, fmt.Errorf("failed to fetch file content: %w", err)
	}

	basePath, err := os.Getwd()

	return *entities.NewFile(basePath, path, []byte(body)), nil

}

func (r *FileRepositoryImpl) GetComponentByFilePath(filePath string) (entities.Component, error) {
	return entities.Component{}, nil
}
