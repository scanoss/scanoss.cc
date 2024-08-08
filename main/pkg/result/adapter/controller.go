package adapter

import (
	"integration-git/main/pkg/result/application/services"
	"integration-git/main/pkg/result/application/use_cases"
	"integration-git/main/pkg/result/common"
	"integration-git/main/pkg/result/filter_adapter"
)

type ResultController struct {
	getResultsUseCase *use_cases.GetResultsUseCase
}

func NewResultController(service *services.ResultService) *ResultController {
	return &ResultController{
		getResultsUseCase: use_cases.NewGetResultsUseCase(service),
	}
}

func (c *ResultController) GetAll(dto *common.RequestResultDTO) ([]common.ResultDTO, error) {
	// convert DTO to filters
	filter := filter_adapter.NewResultFilterFactory().Create(dto)
	results, err := c.getResultsUseCase.GetResults(filter)

	if err != nil {
		return []common.ResultDTO{}, err
	}

	var output []common.ResultDTO

	for _, result := range results {
		output = append(output, common.ResultDTO{
			MatchType: common.MatchType(result.GetMatchType()),
			File:      result.GetFile(),
		})
	}

	return output, nil
}
