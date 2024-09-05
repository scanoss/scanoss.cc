package repository

import "github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"

type ResultRepositoryInMemoryImpl struct {
}

func NewResultRepositoryInMemoryImpl() *ResultRepositoryInMemoryImpl {
	return &ResultRepositoryInMemoryImpl{}
}

func (r *ResultRepositoryInMemoryImpl) GetResults(filter entities.ResultFilter) ([]entities.Result, error) {
	var results []entities.Result

	result1 := entities.Result{}
	result1.SetFile("/inc/format_utils.h")
	result1.SetMatchType("file")

	result2 := entities.Result{}
	result2.SetFile("/src/format_utils.c")
	result2.SetMatchType("snippet")

	result3 := entities.Result{}
	result3.SetFile("/external/inc/crc32c.h")
	result3.SetMatchType("none")

	results = append(results, result1)
	results = append(results, result2)
	results = append(results, result3)

	// Filter scan results
	if filter != nil {
		var filteredResults []entities.Result
		for _, result := range results {
			if filter.IsValid(result) {
				filteredResults = append(filteredResults, result)
			}
		}
		return filteredResults, nil
	}
	return results, nil
}
