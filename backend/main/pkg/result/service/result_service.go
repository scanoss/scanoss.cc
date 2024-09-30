package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

//go:generate mockery --name ResultService --with-expecter
type ResultService interface {
	GetResults(filters entities.ResultFilter) ([]entities.Result, error)
}
