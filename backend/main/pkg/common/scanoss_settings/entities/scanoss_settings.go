package entities

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/result/entities"
)

type ComponentFilterUsage string

var (
	ScanossSettingsJson *ScanossSettings
)

type ScanossSettings struct {
	SettingsFile *SettingsFile
}

type SettingsFile struct {
	Bom Bom `json:"bom"`
}

type Bom struct {
	Include []ComponentFilter `json:"include,omitempty"`
	Remove  []ComponentFilter `json:"remove,omitempty"`
	Replace []ComponentFilter `json:"replace,omitempty"`
}

const (
	File    ComponentFilterUsage = "file"
	Snippet ComponentFilterUsage = "snippet"
)

type ComponentFilter struct {
	Path          string               `json:"path,omitempty"`
	Purl          string               `json:"purl"`
	Usage         ComponentFilterUsage `json:"usage,omitempty"`
	Comment       string               `json:"comment,omitempty"`
	ReplaceWith   string               `json:"replace_with,omitempty"`
	ComponentName string               `json:"component_name,omitempty"`
}

type InitialFilters struct {
	Include []ComponentFilter
	Remove  []ComponentFilter
	Replace []ComponentFilter
}

func (sf *SettingsFile) Equal(other *SettingsFile) (bool, error) {
	originalSettingsFileJson, err := json.Marshal(sf)
	if err != nil {
		return false, fmt.Errorf("error parsing json: %v", err)
	}

	currentContentJSON, err := json.Marshal(other)
	if err != nil {
		return false, fmt.Errorf("error parsing json: %v", err)
	}

	hasChanges := string(originalSettingsFileJson) != string(currentContentJSON)

	return hasChanges, nil
}

func (sf *SettingsFile) GetResultWorkflowState(result entities.Result) entities.WorkflowState {
	included, _ := sf.IsResultIncluded(result)
	removed, _ := sf.IsResultRemoved(result)
	replaced, _ := sf.IsResultReplaced(result)

	if included || removed || replaced {
		return entities.Completed
	}

	return entities.Pending
}

func (sf *SettingsFile) IsResultIncluded(result entities.Result) (bool, int) {
	return sf.IsResultInList(result, sf.Bom.Include)
}

func (sf *SettingsFile) IsResultRemoved(result entities.Result) (bool, int) {
	return sf.IsResultInList(result, sf.Bom.Remove)
}

func (sf *SettingsFile) IsResultReplaced(result entities.Result) (bool, int) {
	return sf.IsResultInList(result, sf.Bom.Replace)
}

func (sf *SettingsFile) IsResultInList(result entities.Result, list []ComponentFilter) (bool, int) {
	i := slices.IndexFunc(list, func(cf ComponentFilter) bool {
		if cf.Path != "" && cf.Purl != "" {
			return cf.Path == result.Path && slices.Contains(*result.Purl, cf.Purl)
		}
		return slices.Contains(*result.Purl, cf.Purl)
	})

	return i != -1, i
}

func (sf *SettingsFile) GetResultFilterConfig(result entities.Result) entities.FilterConfig {
	var filterAction entities.FilterAction
	var filterType entities.FilterType

	if included, i := sf.IsResultIncluded(result); included {
		filterAction = entities.Include
		filterType = getResultFilterType(sf.Bom.Include[i])
	} else if removed, i := sf.IsResultRemoved(result); removed {
		filterAction = entities.Remove
		filterType = getResultFilterType(sf.Bom.Remove[i])
	} else if replaced, i := sf.IsResultReplaced(result); replaced {
		filterAction = entities.Replace
		filterType = getResultFilterType(sf.Bom.Replace[i])
	}

	return entities.FilterConfig{
		Action: filterAction,
		Type:   filterType,
	}
}

func getResultFilterType(cf ComponentFilter) entities.FilterType {
	if cf.Path != "" && cf.Purl != "" {
		return entities.ByFile
	}
	return entities.ByPurl
}

func (sf *SettingsFile) GetBomEntryFromResult(result entities.Result) ComponentFilter {
	if included, i := sf.IsResultIncluded(result); included {
		return sf.Bom.Include[i]
	}

	if removed, i := sf.IsResultRemoved(result); removed {
		return sf.Bom.Remove[i]
	}

	if replaced, i := sf.IsResultReplaced(result); replaced {
		return sf.Bom.Replace[i]
	}

	return ComponentFilter{}
}
