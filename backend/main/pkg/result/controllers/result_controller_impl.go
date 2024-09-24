package controllers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/mappers"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/service"
)

type ResultControllerImpl struct {
	service service.ResultService
	mapper  mappers.ResultMapper
}

func NewResultController(service service.ResultService, mapper mappers.ResultMapper) ResultController {
	return &ResultControllerImpl{
		service: service,
		mapper:  mapper,
	}
}

func (c *ResultControllerImpl) GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error) {
	filter := entities.NewResultFilterFactory().Create(dto)
	results, err := c.service.GetResults(filter)
	if err != nil {
		return []entities.ResultDTO{}, err
	}

	return c.mapper.MapToResultDTOList(results), nil
}
