package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type DefaultFileReader struct{}

func NewDefaultFileReader() *DefaultFileReader {
	return &DefaultFileReader{}
}

func (d DefaultFileReader) ReadFile(filePath string) ([]byte, error) {
	// Open and read the JSON file
	expandedPath := ExpandPath(filePath)
	file, err := os.Open(expandedPath)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	return byteValue, nil
}

func JSONParse[T any](file []byte) (T, error) {
	var intermediateMap T

	if err := json.Unmarshal(file, &intermediateMap); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return intermediateMap, err
	}

	return intermediateMap, nil
}

func WriteJsonFile(path string, in any) error {
	j, err := JSONSerialize(in)
	if err != nil {
		return err
	}

	err = WriteFile(path, j)
	if err != nil {
		return err
	}

	return nil
}

func JSONSerialize(in any) ([]byte, error) {
	out, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return nil, err
	}

	return out, nil
}

func WriteFile(filename string, data []byte) error {
	err := os.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func FileExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
}

func IsWritableFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := fileInfo.Mode()

	return mode&0200 != 0
}

func ExpandPath(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	if path == "~" {
		return homeDir
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}

	return path
}
