package infraestructure

import (
	"integration-git/main/pkg/result/common"
	"integration-git/main/pkg/result/domain"
)

type InMemoryResultRepository struct {
}

func NewInMemoryResultRepository() *InMemoryResultRepository {
	return &InMemoryResultRepository{}
}

func (r *InMemoryResultRepository) GetResults(filter common.ResultFilter) ([]domain.Result, error) {

	var results []domain.Result

	result1 := domain.Result{}
	result1.SetFile("/inc/format_utils.h")
	result1.SetMatchType("file")

	result2 := domain.Result{}
	result2.SetFile("/src/format_utils.c")
	result2.SetMatchType("snippet")

	result3 := domain.Result{}
	result3.SetFile("/external/inc/crc32c.h")
	result3.SetMatchType("none")

	results = append(results, result1)
	results = append(results, result2)
	results = append(results, result3)

	// Filter scan results
	if filter != nil {
		var filteredResults []domain.Result
		for _, result := range results {
			if filter.IsValid(result) {
				filteredResults = append(filteredResults, result)
			}
		}
		return filteredResults, nil
	}
	return results, nil

}
