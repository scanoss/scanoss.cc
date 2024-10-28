package service

import "github.com/scanoss/scanoss.lui/backend/main/entities"

type ResultService interface {
	GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error)
}
