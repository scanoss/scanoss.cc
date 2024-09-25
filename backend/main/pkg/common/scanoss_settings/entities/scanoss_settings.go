package entities

import (
	"fmt"
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
	included, _ := sf.IsResultIncluded(result)
	removed, _ := sf.IsResultRemoved(result)

	if included || removed {
		return entities.Completed
	}

	return entities.Pending
}

func (sf *SettingsFile) IsResultIncluded(result entities.Result) (bool, entities.FilterType) {
	return sf.IsResultInList(result, sf.Bom.Include)
}

func (sf *SettingsFile) IsResultRemoved(result entities.Result) (bool, entities.FilterType) {
	return sf.IsResultInList(result, sf.Bom.Remove)
}

func (sf *SettingsFile) IsResultInList(result entities.Result, list []ComponentFilter) (bool, entities.FilterType) {
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

func (sf *SettingsFile) GetResultFilterConfig(result entities.Result) entities.FilterConfig {
	var filterAction entities.FilterAction
	var filterType entities.FilterType

	if included, t := sf.IsResultIncluded(result); included {
		filterAction = entities.Include
		filterType = t
	} else if removed, t := sf.IsResultRemoved(result); removed {
		filterAction = entities.Remove
		filterType = t

	}

	return entities.FilterConfig{
		Action: filterAction,
		Type:   filterType,
	}
}

func (sf *SettingsFile) AddBomEntry(newEntry ComponentFilter, filterAction string) error {
	var targetList *[]ComponentFilter

	switch filterAction {
	case "remove":
		targetList = &sf.Bom.Remove
	case "include":
		targetList = &sf.Bom.Include
	default:
		return fmt.Errorf("invalid filter action: %s", filterAction)
	}

	sf.removeDuplicatesFromAllLists(newEntry)

	*targetList = append(*targetList, newEntry)

	return nil
}

func (sf *SettingsFile) removeDuplicatesFromAllLists(newEntry ComponentFilter) {
	sf.Bom.Remove = removeDuplicatesFromList(sf.Bom.Remove, newEntry)
	sf.Bom.Include = removeDuplicatesFromList(sf.Bom.Include, newEntry)
}

func removeDuplicatesFromList(list []ComponentFilter, newEntry ComponentFilter) []ComponentFilter {
	result := make([]ComponentFilter, 0, len(list))
	for _, entry := range list {
		if !isDuplicate(entry, newEntry) {
			result = append(result, entry)
		}
	}
	return result
}

func isDuplicate(entry, newEntry ComponentFilter) bool {
	if newEntry.Path == "" {
		return entry.Purl == newEntry.Purl
	}
	return entry.Purl == newEntry.Purl && entry.Path == newEntry.Path
}
