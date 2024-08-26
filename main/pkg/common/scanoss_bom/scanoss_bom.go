package scanoss_bom

import (
	"fmt"
	"sync"
)

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

var (
	bom        *BomFile
	once       sync.Once
	configLock sync.Mutex
)

// LoadConfig reads the configuration and sets it as the singleton instance
func Init() (*BomFile, error) {
	fmt.Println("Initializing in memory bom")
	configLock.Lock()
	defer configLock.Unlock()

	once.Do(func() {
		bom = &BomFile{
			Bom: Bom{
				Include: []ComponentFilter{},
				Remove:  []ComponentFilter{},
			},
		}
	})
	return bom, nil
}

func Get() *BomFile {
	return bom
}
