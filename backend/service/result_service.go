package service

import (
	"context"

	"github.com/scanoss/scanoss.cc/backend/entities"
)

type ResultService interface {
	GetAll(dto *entities.RequestResultDTO) ([]entities.ResultDTO, error)
	SetContext(ctx context.Context)
}
