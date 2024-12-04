package repository

import (
	"encoding/json"
	"fmt"

	"github.com/scanoss/scanoss.lui/backend/entities"
	"github.com/scanoss/scanoss.lui/internal/config"
	"github.com/scanoss/scanoss.lui/internal/utils"
)

type ResultRepositoryJsonImpl struct {
	fr utils.FileReader
}

func NewResultRepositoryJsonImpl(fr utils.FileReader) *ResultRepositoryJsonImpl {
	return &ResultRepositoryJsonImpl{
		fr: fr,
	}
}

func (r *ResultRepositoryJsonImpl) GetResults(filter entities.ResultFilter) ([]entities.Result, error) {
	// Path to your JSON file
	resultFilePath := config.Get().ResultFilePath
	resultByte, err := r.fr.ReadFile(resultFilePath)
	if err != nil {
		fmt.Println(err)
		return []entities.Result{}, entities.ErrReadingResultFile
	}

	scanResults, err := r.parseScanResults(resultByte)
	if err != nil {
		return []entities.Result{}, entities.ErrParsingResultFile
	}

	if filter == nil {
		return scanResults, nil
	}

	// Filter scan results
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

func (r *ResultRepositoryJsonImpl) parseScanResults(resultByte []byte) ([]entities.Result, error) {
	var intermediateMap map[string][]entities.Match
	if err := json.Unmarshal(resultByte, &intermediateMap); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return []entities.Result{}, err
	}

	var scanResults []entities.Result

	for key, matches := range intermediateMap {
		for _, match := range matches {
			scanResult := entities.Result{}
			scanResult.Path = key
			scanResult.MatchType = match.ID
			scanResult.Purl = &match.Purl
			scanResult.ComponentName = match.ComponentName
			scanResults = append(scanResults, scanResult)
		}
	}

	return scanResults, nil
}
