package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

type ResultService interface {
	GetResults(filter entities.ResultFilter) ([]entities.Result, error)
}
