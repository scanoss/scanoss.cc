package entities

import "sync"

type ComponentFilterUsage string

var (
	BomJson    *ScanossBom
	once       sync.Once
	configLock sync.Mutex
)

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

type ScanossBom struct {
	BomFile *BomFile
}

func (sc *ScanossBom) Init(bomFile *BomFile) {
	configLock.Lock()
	defer configLock.Unlock()
	// once.Do(func() {
	// 	sc.BomFile = bomFile
	// 	BomJson = sc
	// })
}
