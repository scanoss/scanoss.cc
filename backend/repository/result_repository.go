package repository

import "github.com/scanoss/scanoss.lui/backend/entities"

type ResultRepository interface {
	GetResults(filters entities.ResultFilter) ([]entities.Result, error)
}
