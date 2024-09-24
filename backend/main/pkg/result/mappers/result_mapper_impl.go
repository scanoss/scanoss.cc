package mappers

import (
	"fmt"

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
	// bomFile, err := m.bomService.GetSettingsFile()
	// // Example logic:
	// if bomFile.IsFileCompleted(result.Path) {
	// 	return entities.Completed
	// }
	// return entities.Pending
	return entities.Pending
}

func (m *ResultMapperImpl) mapFilterConfig(result entities.Result) entities.FilterConfig {
	// TODO: Check if should return an error?
	settingsFile, err := m.scanossSettingsService.GetSettingsFile()
	if err != nil {
		fmt.Println("Error reading scan settings file")
	}

	fmt.Println("BOM FILE CONFIG....")
	fmt.Println(settingsFile)

	// action := entities.Include
	// filterType := entities.ByFile
	// if bomFile.ShouldExclude(result.Path) {
	// 	action = entities.Remove
	// }
	// if bomFile.IsPurlFiltered(result.Purl) {
	// 	filterType = entities.ByPurl
	// }
	// return entities.FilterConfig{
	// 	Action: action,
	// 	Type:   filterType,
	// }
	return entities.FilterConfig{
		Action: entities.Include,
		Type:   entities.ByFile,
	}
}
