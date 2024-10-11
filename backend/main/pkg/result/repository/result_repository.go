package repository

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

type ResultRepository interface {
	GetResults(filters entities.ResultFilter) ([]entities.Result, error)
}
