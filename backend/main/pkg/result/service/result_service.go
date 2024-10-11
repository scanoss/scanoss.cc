package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

type ResultService interface {
	GetResults(filters entities.ResultFilter) ([]entities.Result, error)
}
