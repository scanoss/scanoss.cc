package infraestructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"integration-git/main/pkg/component/domain"
	"io"
	"os"
)

var (
	ErrOpeningResult = errors.New("error opening result file")
	ErrReadingResult = errors.New("error reading result file")
)

type JSONComponentRepository struct {
	resultFilePath string
}

func NewComponentRepository(resultFilePath string) *JSONComponentRepository {
	return &JSONComponentRepository{
		resultFilePath: resultFilePath,
	}
}

func (r *JSONComponentRepository) FindByFilePath(path string) (domain.Component, error) {

	resultFileBytes, err := r.readResultFile()
	if err != nil {
		return domain.Component{}, err
	}
	results, err := r.parseScanResults(resultFileBytes)
	if err != nil {
		return domain.Component{}, err
	}

	components := results[path]
	return components[0], nil
}

func (r *JSONComponentRepository) readResultFile() ([]byte, error) {
	// Open and read the JSON file
	jsonFile, err := os.Open(r.resultFilePath)
	if err != nil {
		return []byte{}, ErrOpeningResult
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return []byte{}, ErrReadingResult
	}

	return byteValue, nil
}

func (r *JSONComponentRepository) parseScanResults(resultByte []byte) (map[string][]domain.Component, error) {
	var intermediateMap map[string][]domain.Component

	if err := json.Unmarshal(resultByte, &intermediateMap); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return intermediateMap, err
	}

	return intermediateMap, nil
}
