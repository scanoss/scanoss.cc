package handlers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

type ResultHandler struct {
	resultModule *result.Module
}

func NewResultHandler(scanossBomModule *scanoss_settings.Module) *ResultHandler {
	return &ResultHandler{
		resultModule: result.NewModule(scanossBomModule),
	}
}

// *****  RESULT MODULE ***** //
func (rh *ResultHandler) ResultGetAll(dto *entities.RequestResultDTO) []entities.ResultDTO {
	results, _ := rh.resultModule.Controller.GetAll(dto)
	return results
}
