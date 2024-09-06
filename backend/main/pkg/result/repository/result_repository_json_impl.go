package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
)

type ResultRepositoryJsonImpl struct {
}

func NewResultRepositoryJsonImpl() *ResultRepositoryJsonImpl {
	return &ResultRepositoryJsonImpl{}
}

func (r *ResultRepositoryJsonImpl) GetResults(filter entities.ResultFilter) ([]entities.Result, error) {
	// Path to your JSON file
	resultFilePath := config.Get().ResultFilePath
	resultByte, err := utils.ReadFile(resultFilePath)
	if err != nil {
		fmt.Println(err)
		return []entities.Result{}, entities.ErrReadingResultFile
	}

	scanResults, err := r.parseScanResults(resultByte)

	// Filter scan results
	if filter != nil {
		var filteredResults []entities.Result
		for _, result := range scanResults {
			if result.IsEmpty() {
				continue
			}

			if filter.IsValid(result) {
				filteredResults = append(filteredResults, result)
			}
		}
		return filteredResults, nil
	}
	if err != nil {
		return []entities.Result{}, entities.ErrParsingResultFile
	}

	return scanResults, err
}

func (r *ResultRepositoryJsonImpl) readResultFile(resultFilePath string) ([]byte, error) {
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

func (r *ResultRepositoryJsonImpl) parseScanResults(resultByte []byte) ([]entities.Result, error) {

	// Unmarshal JSON data into a map with dynamic keys
	var intermediateMap map[string][]entities.Match
	if err := json.Unmarshal(resultByte, &intermediateMap); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return []entities.Result{}, err
	}

	// Create a slice to store ScanResult objects
	var scanResults []entities.Result

	// Iterate over the dynamic keys and convert to ScanResult
	for key, matches := range intermediateMap {
		for _, match := range matches {
			scanResult := entities.Result{ // Use the dynamic key for the File field
			}
			scanResult.SetFile(key)
			scanResult.SetMatchType(match.ID)
			scanResults = append(scanResults, scanResult)
		}
	}

	return scanResults, nil
}
