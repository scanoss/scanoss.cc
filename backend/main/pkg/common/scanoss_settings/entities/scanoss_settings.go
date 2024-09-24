package entities

import "reflect"

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
