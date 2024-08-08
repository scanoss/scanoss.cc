package repository

import (
	"integration-git/main/pkg/result/common"
	"integration-git/main/pkg/result/domain"
)

type ResultRepository interface {
	GetResults(filter common.ResultFilter) ([]domain.Result, error)
}
