package repository

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

//go:generate mockery --name ResultRepository --with-expecter
type ResultRepository interface {
	GetResults(filters entities.ResultFilter) ([]entities.Result, error)
}
