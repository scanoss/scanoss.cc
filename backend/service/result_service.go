package service

import (
	"context"

	"github.com/scanoss/scanoss.lui/backend/entities"
)

type ResultService interface {
	GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error)
	SetContext(ctx context.Context)
}
