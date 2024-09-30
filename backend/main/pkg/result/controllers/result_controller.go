package controllers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

//go:generate mockery --name ResultController --with-expecter
type ResultController interface {
	GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error)
}
