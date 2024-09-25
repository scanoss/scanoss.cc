package entities

import (
	"reflect"
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
}

const (
	File    ComponentFilterUsage = "file"
	Snippet ComponentFilterUsage = "snippet"
)

type ComponentFilter struct {
	Path  string               `json:"path,omitempty"`
	Purl  string               `json:"purl"`
	Usage ComponentFilterUsage `json:"usage,omitempty"`
}

func (sf *SettingsFile) Equal(other *SettingsFile) bool {
	return reflect.DeepEqual(sf.Bom, other.Bom)
}

func (sf *SettingsFile) GetResultWorkflowState(result entities.Result) entities.WorkflowState {
	included, _ := sf.isResultIncluded(result)
	removed, _ := sf.isResultRemoved(result)

	if included || removed {
		return entities.Completed
	}

	return entities.Pending
}

func (sf *SettingsFile) GetResultFilterConfig(result entities.Result) entities.FilterConfig {
	var filterAction entities.FilterAction
	var filterType entities.FilterType

	if included, t := sf.isResultIncluded(result); included {
		filterAction = entities.Include
		filterType = t
	} else if removed, t := sf.isResultRemoved(result); removed {
		filterAction = entities.Remove
		filterType = t

	}

	return entities.FilterConfig{
		Action: filterAction,
		Type:   filterType,
	}
}

func (sf *SettingsFile) isResultIncluded(result entities.Result) (bool, entities.FilterType) {
	return sf.isResultInList(result, sf.Bom.Include)
}

func (sf *SettingsFile) isResultRemoved(result entities.Result) (bool, entities.FilterType) {
	return sf.isResultInList(result, sf.Bom.Remove)
}

func (sf *SettingsFile) isResultInList(result entities.Result, list []ComponentFilter) (bool, entities.FilterType) {
	i := slices.IndexFunc(list, func(cf ComponentFilter) bool {
		if cf.Path != "" && cf.Purl != "" {
			return cf.Path == result.Path && slices.Contains(result.Purl, cf.Purl)
		}
		return slices.Contains(result.Purl, cf.Purl)
	})

	if i != -1 {
		cf := list[i]
		if cf.Path != "" && cf.Purl != "" {
			return true, entities.ByFile
		}
		return true, entities.ByPurl
	}

	return false, entities.ByPurl
}
