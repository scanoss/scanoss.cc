package controllers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

type ResultController interface {
	GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error)
}
