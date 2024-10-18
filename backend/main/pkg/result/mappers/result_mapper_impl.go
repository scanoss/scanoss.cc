package mappers

import (
	"github.com/labstack/gommon/log"
	purlutils "github.com/scanoss/go-purl-helper/pkg"
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
		MatchType: entities.MatchType(result.MatchType),
		Path:      result.Path,
		Purl: entities.ResultPurl{
			Detected:         (*result.Purl)[0],
			Concluded:        m.mapConcludedPurl(result),
			ConcludedPurlUrl: m.mapConcludedPurlUrl(result),
		},
		WorkflowState: m.mapWorkflowState(result),
		FilterConfig:  m.mapFilterConfig(result),
		Comment:       m.mapComment(result),
	}
}

func (m ResultMapperImpl) mapComment(result entities.Result) string {
	return m.scanossSettings.SettingsFile.GetBomEntryFromResult(result).Comment
}

func (m ResultMapperImpl) mapConcludedPurl(result entities.Result) string {
	return m.scanossSettings.SettingsFile.GetBomEntryFromResult(result).ReplaceWith
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

func (m *ResultMapperImpl) mapConcludedPurlUrl(result entities.Result) string {
	concludedPurl := m.mapConcludedPurl(result)
	if concludedPurl == "" {
		return ""
	}

	purlObject, err := purlutils.PurlFromString(concludedPurl)
	if err != nil {
		log.Errorf("Error parsing concluded purl: %v", err)
		return ""
	}

	purlUrl, err := purlutils.ProjectUrl(purlObject.Name, purlObject.Type)
	if err != nil {
		log.Errorf("Error getting project url: %v", err)
		return ""
	}

	return purlUrl
}
