package mappers

import "github.com/scanoss/scanoss.cc/backend/entities"

type ResultMapper interface {
	MapToResultDTO(result entities.Result) entities.ResultDTO
	MapToResultDTOList(results []entities.Result) []entities.ResultDTO
}
