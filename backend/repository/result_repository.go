package repository

import "github.com/scanoss/scanoss.cc/backend/entities"

type ResultRepository interface {
	GetResults(filters entities.ResultFilter) ([]entities.Result, error)
}
