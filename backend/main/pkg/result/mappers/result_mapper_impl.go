package mappers

import (
	scanossSettingsEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

type ResultMapperImpl struct {
	scanossSettings *scanossSettingsEntities.ScanossSettings
}

func NewResultMapper(scanossSettings *scanossSettingsEntities.ScanossSettings) ResultMapper {
	return &ResultMapperImpl{
		scanossSettings: scanossSettings,
	}
}

func (m ResultMapperImpl) MapToResultDTO(result entities.Result) entities.ResultDTO {
	return entities.ResultDTO{
		MatchType:     entities.MatchType(result.MatchType),
		Path:          result.Path,
		Purl:          (*result.Purl)[0],
		WorkflowState: m.mapWorkflowState(result),
		FilterConfig:  m.mapFilterConfig(result),
		Comment:       m.mapComment(result),
	}
}

func (m ResultMapperImpl) mapComment(result entities.Result) string {
	return m.scanossSettings.SettingsFile.GetResultComment(result)
}

func (m ResultMapperImpl) MapToResultDTOList(results []entities.Result) []entities.ResultDTO {
	output := make([]entities.ResultDTO, len(results))

	for i, v := range results {
		output[i] = m.MapToResultDTO(v)
	}

	return output
}

func (m *ResultMapperImpl) mapWorkflowState(result entities.Result) entities.WorkflowState {
	return m.scanossSettings.SettingsFile.GetResultWorkflowState(result)
}

func (m *ResultMapperImpl) mapFilterConfig(result entities.Result) entities.FilterConfig {
	return m.scanossSettings.SettingsFile.GetResultFilterConfig(result)
}
