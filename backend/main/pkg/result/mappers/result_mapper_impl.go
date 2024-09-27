package mappers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/service"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

type ResultMapperImpl struct {
	scanossSettingsService service.ScanossSettingsService
}

func NewResultMapper(scanossSettingsService service.ScanossSettingsService) ResultMapper {
	return &ResultMapperImpl{
		scanossSettingsService: scanossSettingsService,
	}
}

func (m ResultMapperImpl) MapToResultDTO(result entities.Result) entities.ResultDTO {
	return entities.ResultDTO{
		MatchType:     entities.MatchType(result.MatchType),
		Path:          result.Path,
		WorkflowState: m.mapWorkflowState(result),
		FilterConfig:  m.mapFilterConfig(result),
	}
}

func (m ResultMapperImpl) MapToResultDTOList(results []entities.Result) []entities.ResultDTO {
	output := make([]entities.ResultDTO, len(results))

	for i, v := range results {
		output[i] = m.MapToResultDTO(v)
	}

	return output
}

func (m *ResultMapperImpl) mapWorkflowState(result entities.Result) entities.WorkflowState {
	settingsFile := m.scanossSettingsService.GetSettingsFile()

	return settingsFile.GetResultWorkflowState(result)
}

func (m *ResultMapperImpl) mapFilterConfig(result entities.Result) entities.FilterConfig {
	settingsFile := m.scanossSettingsService.GetSettingsFile()

	return settingsFile.GetResultFilterConfig(result)
}
