package repository

import "github.com/scanoss/scanoss.lui/backend/main/entities"

type ResultRepository interface {
	GetResults(filters entities.ResultFilter) ([]entities.Result, error)
}
