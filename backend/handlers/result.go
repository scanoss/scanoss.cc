package handlers

import (
	"context"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/result"
	module_result_adapter "github.com/scanoss/scanoss.lui/backend/main/pkg/result/adapter"
)

type ResultHandler struct {
	ctx          context.Context
	resultModule *result.Module
}

func NewResultHandler() *ResultHandler {
	return &ResultHandler{
		resultModule: result.NewModule(),
	}
}

// *****  RESULT MODULE ***** //
func (rh *ResultHandler) ResultGetAll(dto *module_result_adapter.RequestResultDTO) []module_result_adapter.ResultDTO {
	results, _ := rh.resultModule.Controller.GetAll(dto)
	return results
}
