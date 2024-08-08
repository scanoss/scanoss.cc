package infrastructure

import (
	"crypto/md5"
	"errors"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/file/domain"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var (
	ErrReadingFile     = errors.New("error reading file")
	ErrFetchingContent = errors.New("Error fetching remote file")
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
	// Get the file name from the path
	fileName := filepath.Base(filePath)

	data, err := os.ReadFile(absolutePath)
	if err != nil {
		return domain.File{}, ErrReadingFile
	}
	// Convert bytes to string
	content := string(data)

	file := domain.NewFile()
	file.SetLocalContent(content)
	file.SetPath(filePath)
	file.SetName(fileName)
	return *file, nil
}

func (r *FileRepository) ReadRemoteFile(path string) (domain.File, error) {
	file, err := r.ReadLocalFile(path)
	if err != nil {
		return domain.File{}, err
	}
	fileMD5 := r.md5Hash(file.GetLocalContent())
	remoteFileContent, err := r.fetchContent(fileMD5)
	if err != nil {
		return domain.File{}, err
	}
	file.SetRemoteContent(remoteFileContent)
	return file, nil
}

func (r *FileRepository) md5Hash(content string) string {
	// Create a new hash.Hash computing the MD5 checksum
	hash := md5.New()

	// Write the string content into the hash
	if _, err := io.WriteString(hash, content); err != nil {
		fmt.Println("Error writing content:", err)
		return ""
	}
	// Get the MD5 checksum as a byte slice
	checksum := hash.Sum(nil)
	// Convert the byte slice to a hexadecimal string
	checksumStr := fmt.Sprintf("%x", checksum)
	return checksumStr
}

func (r *FileRepository) fetchContent(fileMD5 string) (string, error) {
	baseUrl := config.Get().Scanoss.ApiUrl
	token := config.Get().Scanoss.ApiToken
	// Create a new HTTP request
	url := fmt.Sprintf("%s/%s/%s", baseUrl, "file_contents", fileMD5)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Add the header if the API key is provided
	if token != "" {
		req.Header.Set("X-Session", token)
	}

	// Create a new HTTP client and perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check if the status code is OK
	if resp.StatusCode != http.StatusOK {
		return "", ErrFetchingContent
	}

	return string(body), nil
}
