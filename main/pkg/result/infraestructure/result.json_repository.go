package infraestructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/result/common"
	"integration-git/main/pkg/result/domain"
	"integration-git/main/pkg/utils"
	"io"
	"os"
)

var (
	ErrReadingResultFile = errors.New("error reading result file")
	ErrParsingResultFile = errors.New("error parsing result file")
)

type Match struct {
	ID string `json:"id"`
}

type JsonResultRepository struct {
}

func NewJsonResultRepository() *JsonResultRepository {
	return &JsonResultRepository{}
}

func (r *JsonResultRepository) GetResults(filter common.ResultFilter) ([]domain.Result, error) {
	// Path to your JSON file
	resultFilePath := config.Get().Scanoss.ResultFilePath

	resultByte, err := utils.ReadFile(resultFilePath)
	if err != nil {
		return []domain.Result{}, ErrReadingResultFile
	}

	scanResults, err := r.parseScanResults(resultByte)

	// Filter scan results
	if filter != nil {
		var filteredResults []domain.Result
		for _, result := range scanResults {
			if filter.IsValid(result) {
				filteredResults = append(filteredResults, result)
			}
		}
		return filteredResults, nil
	}
	if err != nil {
		return []domain.Result{}, ErrParsingResultFile
	}

	return scanResults, err

}

func (r *JsonResultRepository) readResultFile(resultFilePath string) ([]byte, error) {
	// Open and read the JSON file
	jsonFile, err := os.Open(resultFilePath)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return []byte{}, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return []byte{}, err
	}

	return byteValue, nil
}

func (r *JsonResultRepository) parseScanResults(resultByte []byte) ([]domain.Result, error) {

	// Unmarshal JSON data into a map with dynamic keys
	var intermediateMap map[string][]Match
	if err := json.Unmarshal(resultByte, &intermediateMap); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return []domain.Result{}, err
	}

	// Create a slice to store ScanResult objects
	var scanResults []domain.Result

	// Iterate over the dynamic keys and convert to ScanResult
	for key, matches := range intermediateMap {
		for _, match := range matches {
			scanResult := domain.Result{ // Use the dynamic key for the File field
			}
			scanResult.SetFile(key)
			scanResult.SetMatchType(match.ID)
			scanResults = append(scanResults, scanResult)
		}
	}

	return scanResults, nil
}
