package mappers

import "github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"

//go:generate mockery --name ResultMapper --with-expecter
type ResultMapper interface {
	MapToResultDTO(result entities.Result) entities.ResultDTO
	MapToResultDTOList(results []entities.Result) []entities.ResultDTO
}
