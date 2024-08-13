package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadFile(filePath string) ([]byte, error) {
	// Open and read the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return []byte{}, err
	}

	return byteValue, nil
}

func JSONParse[T any](file []byte) (map[string][]T, error) {
	var intermediateMap map[string][]T

	if err := json.Unmarshal(file, &intermediateMap); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return intermediateMap, err
	}

	return intermediateMap, nil
}
