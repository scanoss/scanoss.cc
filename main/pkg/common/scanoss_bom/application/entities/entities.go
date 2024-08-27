package entities

type ComponentFilterUsage string

const (
	File    ComponentFilterUsage = "file"
	Snippet ComponentFilterUsage = "snippet"
)

type ComponentFilter struct {
	Path    string               `json:"path"`
	Purl    string               `json:"purl"`
	Usage   ComponentFilterUsage `json:"usage,omitempty"`
	Version string               `json:"version"`
}

type Bom struct {
	Include []ComponentFilter `json:"include"`
	Remove  []ComponentFilter `json:"remove"`
}

type BomFile struct {
	Bom Bom `json:"bom"`
}
