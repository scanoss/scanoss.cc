package mappers

import "github.com/scanoss/scanoss.lui/backend/main/entities"

type ResultMapper interface {
	MapToResultDTO(result entities.Result) entities.ResultDTO
	MapToResultDTOList(results []entities.Result) []entities.ResultDTO
}
