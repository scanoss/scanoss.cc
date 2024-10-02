package handlers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

type ResultHandler struct {
	resultModule *result.Module
}

func NewResultHandler() *ResultHandler {
	return &ResultHandler{
		resultModule: result.NewModule(),
	}
}

func (rh *ResultHandler) ResultGetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error) {
	results, err := rh.resultModule.Controller.GetAll(dto)
	return results, err
}
