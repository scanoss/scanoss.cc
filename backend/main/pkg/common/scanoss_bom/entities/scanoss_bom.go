package entities

type ComponentFilterUsage string

var (
	BomJson *ScanossBom
)

const (
	File    ComponentFilterUsage = "file"
	Snippet ComponentFilterUsage = "snippet"
)

type ComponentFilter struct {
	Path  string               `json:"path,omitempty"`
	Purl  string               `json:"purl"`
	Usage ComponentFilterUsage `json:"usage,omitempty"`
}

type Bom struct {
	Include []ComponentFilter `json:"include,omitempty"`
	Remove  []ComponentFilter `json:"remove,omitempty"`
}

type BomFile struct {
	Bom Bom `json:"bom"`
}

type ScanossBom struct {
	BomFile *BomFile
}
