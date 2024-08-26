package repositories

import (
	"errors"
	"fmt"
	"integration-git/main/pkg/common/config"
	"integration-git/main/pkg/common/scanoss_bom"
	"integration-git/main/pkg/component/entities"
	"integration-git/main/pkg/utils"
	"sort"
)

var (
	ErrOpeningResult       = errors.New("error opening result file")
	ErrReadingResult       = errors.New("error reading result file")
	ErrInvalidFilterAction = errors.New("invalid filter action")
)

var licenseSourceOrder = map[string]int{
	"component_declared": 1,
	"file_header":        2,
	"scancode":           3,
	"license_file":       4,
	"file_spdx_tag":      5,
}

type JSONComponentRepository struct {
}

func NewComponentRepository() *JSONComponentRepository {
	return &JSONComponentRepository{}
}

func (r *JSONComponentRepository) FindByFilePath(path string) (entities.Component, error) {
	resultFilePath := config.Get().ResultFilePath

	resultFileBytes, err := utils.ReadFile(resultFilePath)
	if err != nil {
		return entities.Component{}, err
	}
	results, err := utils.JSONParse[map[string][]entities.Component](resultFileBytes)
	if err != nil {
		return entities.Component{}, err
	}

	// Gets components from results
	components := results[path]

	// Order component licenses by source
	if len(components[0].Licenses) > 0 {
		r.orderComponentLicensesBySourceType(&components[0])
	}

	return components[0], nil
}

func (r *JSONComponentRepository) orderComponentLicensesBySourceType(component *entities.Component) {
	sort.Slice(component.Licenses, func(i, j int) bool {
		return licenseSourceOrder[component.Licenses[i].Source] < licenseSourceOrder[component.Licenses[j].Source]
	})
}

func (r *JSONComponentRepository) InsertComponentFilter(dto *entities.ComponentFilterDTO) error {
	newFilter := &scanoss_bom.ComponentFilter{
		Path:    dto.Path,
		Purl:    dto.Purl,
		Usage:   scanoss_bom.ComponentFilterUsage(dto.Usage),
		Version: dto.Version,
	}

	/*	scanSettingsFileBytes, err := utils.ReadFile(config.Get().ScanSettingsFilePath)
		if err != nil {
			return err
		}*/

	/*	parsedFile, err := utils.JSONParse[entities.ScanSettingsFile](scanSettingsFileBytes)
		if err != nil {
			return err
		}*/

	bom := scanoss_bom.Get()

	fmt.Println(bom.Bom.Include)
	fmt.Println(bom.Bom.Remove)

	if err := insertNewComponentFilter(bom, newFilter, dto.Action); err != nil {
		return err
	}

	/*	if err := utils.WriteJsonFile(config.Get().ScanSettingsFilePath, bom); err != nil {
		return err
	}*/

	return nil
}

func findComponent(newComponent *scanoss_bom.ComponentFilter, components []scanoss_bom.ComponentFilter) *scanoss_bom.ComponentFilter {
	for _, c := range components {
		if newComponent.Path == c.Path && newComponent.Purl == c.Purl && newComponent.Version == c.Version && newComponent.Usage == c.Usage {
			return &c
		}
	}
	return nil
}

func deleteComponent(newComponent *scanoss_bom.ComponentFilter, components *[]scanoss_bom.ComponentFilter) {
	for i := range *components {
		if (*components)[i].Path == newComponent.Path && (*components)[i].Version == newComponent.Version && (*components)[i].Usage == newComponent.Usage && (*components)[i].Purl == newComponent.Purl {
			*components = append((*components)[:i], (*components)[i+1:]...)
			break
		}
	}
}

func insertComponent(newComponent *scanoss_bom.ComponentFilter, a *[]scanoss_bom.ComponentFilter, b *[]scanoss_bom.ComponentFilter) {
	if findComponent(newComponent, *a) == nil {
		*a = append((*a), *newComponent)

		// Delete if exists
		c := findComponent(newComponent, *b)
		if c != nil {
			deleteComponent(newComponent, b)
		}
	}
}

func insertNewComponentFilter(bomFile *scanoss_bom.BomFile, newFilter *scanoss_bom.ComponentFilter, action entities.FilterAction) error {
	switch action {
	case entities.Include:
		insertComponent(newFilter, &bomFile.Bom.Include, &bomFile.Bom.Remove)
	case entities.Remove:
		insertComponent(newFilter, &bomFile.Bom.Remove, &bomFile.Bom.Include)
	default:
		return ErrInvalidFilterAction
	}

	return nil
}
