package controllers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/service"
)

type ResultControllerImpl struct {
	service service.ResultService
}

func NewResultController(service service.ResultService) ResultController {
	return &ResultControllerImpl{
		service: service,
	}
}

func (c *ResultControllerImpl) GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error) {
	// convert DTO to filters
	filter := entities.NewResultFilterFactory().Create(dto)
	results, err := c.service.GetResults(filter)

	if err != nil {
		return []entities.ResultDTO{}, err
	}

	var output []entities.ResultDTO

	for _, result := range results {
		output = append(output, entities.ResultDTO{
			MatchType: entities.MatchType(result.GetMatchType()),
			File:      result.GetFile(),
		})
	}

	return output, nil
}
