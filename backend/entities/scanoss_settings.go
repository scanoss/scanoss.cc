package entities

import (
	"encoding/json"
	"fmt"
	"slices"
)

type ComponentFilterUsage string

var (
	ScanossSettingsJson *ScanossSettings
)

type ScanossSettings struct {
	SettingsFile *SettingsFile
}

type SettingsFile struct {
	Settings ScanossSettingsSchema `json:"settings,omitempty"`
	Bom      Bom                   `json:"bom,omitempty"`
}

type ScanossSettingsSchema struct {
	Skip SkipSettings `json:"skip,omitempty"`
}

type SkipSettings struct {
	Patterns []string `json:"patterns,omitempty"`
	Sizes    Sizes    `json:"sizes,omitempty"`
}

type Sizes struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

type Bom struct {
	Include []ComponentFilter `json:"include,omitempty"`
	Remove  []ComponentFilter `json:"remove,omitempty"`
	Replace []ComponentFilter `json:"replace,omitempty"`
}

type ComponentFilter struct {
	Path        string               `json:"path,omitempty"`
	Purl        string               `json:"purl"`
	Usage       ComponentFilterUsage `json:"usage,omitempty"`
	Comment     string               `json:"comment,omitempty"`
	ReplaceWith string               `json:"replace_with,omitempty"`
	License     string               `json:"license,omitempty"`
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

func (sf *SettingsFile) GetResultWorkflowState(result Result) WorkflowState {
	included, _ := sf.IsResultIncluded(result)
	removed, _ := sf.IsResultRemoved(result)
	replaced, _ := sf.IsResultReplaced(result)

	if included || removed || replaced {
		return Completed
	}

	return Pending
}

func (sf *SettingsFile) IsResultIncluded(result Result) (bool, int) {
	return sf.IsResultInList(result, sf.Bom.Include)
}

func (sf *SettingsFile) IsResultRemoved(result Result) (bool, int) {
	return sf.IsResultInList(result, sf.Bom.Remove)
}

func (sf *SettingsFile) IsResultReplaced(result Result) (bool, int) {
	return sf.IsResultInList(result, sf.Bom.Replace)
}

func (sf *SettingsFile) IsResultInList(result Result, list []ComponentFilter) (bool, int) {
	i := slices.IndexFunc(list, func(cf ComponentFilter) bool {
		if cf.Path != "" && cf.Purl != "" {
			return cf.Path == result.Path && slices.Contains(*result.Purl, cf.Purl)
		}
		return slices.Contains(*result.Purl, cf.Purl)
	})

	return i != -1, i
}

func (sf *SettingsFile) GetResultFilterConfig(result Result) FilterConfig {
	var filterAction FilterAction
	var filterType FilterType

	if included, i := sf.IsResultIncluded(result); included {
		filterAction = Include
		filterType = getResultFilterType(sf.Bom.Include[i])
	} else if removed, i := sf.IsResultRemoved(result); removed {
		filterAction = Remove
		filterType = getResultFilterType(sf.Bom.Remove[i])
	} else if replaced, i := sf.IsResultReplaced(result); replaced {
		filterAction = Replace
		filterType = getResultFilterType(sf.Bom.Replace[i])
	}

	return FilterConfig{
		Action: filterAction,
		Type:   filterType,
	}
}

func getResultFilterType(cf ComponentFilter) FilterType {
	if cf.Path != "" && cf.Purl != "" {
		return ByFile
	}
	return ByPurl
}

func (sf *SettingsFile) GetBomEntryFromResult(result Result) ComponentFilter {
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
