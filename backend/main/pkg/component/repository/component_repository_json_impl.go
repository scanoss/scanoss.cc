package repository

import (
	"errors"
	"sort"

	scanossSettingsEntities "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/config"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/component/entities"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/utils"
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

func NewJSONComponentRepository() *JSONComponentRepository {
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
	newFilter := &scanossSettingsEntities.ComponentFilter{
		Path:  dto.Path,
		Purl:  dto.Purl,
		Usage: scanossSettingsEntities.ComponentFilterUsage(dto.Usage),
	}
	settingsFile := scanossSettingsEntities.ScanossSettingsJson.SettingsFile

	if err := insertNewComponentFilter(settingsFile, newFilter, dto.Action); err != nil {
		return err
	}

	return nil
}

func findComponent(newComponent *scanossSettingsEntities.ComponentFilter, components []scanossSettingsEntities.ComponentFilter) *scanossSettingsEntities.ComponentFilter {
	for _, c := range components {
		if newComponent.Path == c.Path && newComponent.Purl == c.Purl && newComponent.Usage == c.Usage {
			return &c
		}
	}
	return nil
}

func deleteComponent(newComponent *scanossSettingsEntities.ComponentFilter, components *[]scanossSettingsEntities.ComponentFilter) {
	for i := range *components {
		if (*components)[i].Path == newComponent.Path && (*components)[i].Usage == newComponent.Usage && (*components)[i].Purl == newComponent.Purl {
			*components = append((*components)[:i], (*components)[i+1:]...)
			break
		}
	}
}

func insertComponent(newComponent *scanossSettingsEntities.ComponentFilter, a *[]scanossSettingsEntities.ComponentFilter, b *[]scanossSettingsEntities.ComponentFilter) {
	if findComponent(newComponent, *a) == nil {
		*a = append((*a), *newComponent)

		// Delete if exists
		c := findComponent(newComponent, *b)
		if c != nil {
			deleteComponent(newComponent, b)
		}
	}
}

func insertNewComponentFilter(settingsFile *scanossSettingsEntities.SettingsFile, newFilter *scanossSettingsEntities.ComponentFilter, action entities.FilterAction) error {
	switch action {
	case entities.Include:
		insertComponent(newFilter, &settingsFile.Bom.Include, &settingsFile.Bom.Remove)
	case entities.Remove:
		insertComponent(newFilter, &settingsFile.Bom.Remove, &settingsFile.Bom.Include)
	default:
		return ErrInvalidFilterAction
	}

	return nil
}
