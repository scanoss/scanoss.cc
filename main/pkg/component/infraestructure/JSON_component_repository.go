package infraestructure

import (
	"errors"
	"integration-git/main/pkg/component/domain"
	"integration-git/main/pkg/utils"
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

	resultFileBytes, err := utils.ReadFile(r.resultFilePath)
	if err != nil {
		return domain.Component{}, err
	}
	results, err := utils.JSONParse[domain.Component](resultFileBytes)
	if err != nil {
		return domain.Component{}, err
	}

	components := results[path]
	return components[0], nil
}
