package handler

import (
	"context"
	"integration-git/main/pkg/result"
	module_result_adapter "integration-git/main/pkg/result/common"
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
