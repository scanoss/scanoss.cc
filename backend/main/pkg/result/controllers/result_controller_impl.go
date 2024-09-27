package controllers

import (
	"sort"

	"github.com/go-playground/validator"
	"github.com/labstack/gommon/log"
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
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		log.Errorf("Validation error: %v", err)
		return []entities.ResultDTO{}, err
	}

	filter := entities.NewResultFilterFactory().Create(dto)
	results, err := c.service.GetResults(filter)
	if err != nil {
		return []entities.ResultDTO{}, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})

	return c.mapper.MapToResultDTOList(results), nil
}
