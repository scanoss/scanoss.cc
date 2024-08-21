package repositories

import (
	"errors"
	"integration-git/main/pkg/common/config"
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
	newFilter := &entities.ComponentFilter{
		Path:    dto.Path,
		Purl:    dto.Purl,
		Usage:   entities.ComponentFilterUsage(dto.Usage),
		Version: dto.Version,
	}

	scanSettingsFileBytes, err := utils.ReadFile(config.Get().ScanSettingsFilePath)
	if err != nil {
		return err
	}

	parsedFile, err := utils.JSONParse[entities.ScanSettingsFile](scanSettingsFileBytes)
	if err != nil {
		return err
	}

	if err := insertNewComponentFilter(&parsedFile, newFilter, dto.Action); err != nil {
		return err
	}

	if err := utils.WriteJsonFile(config.Get().ScanSettingsFilePath, parsedFile); err != nil {
		return err
	}

	return nil
}

func insertNewComponentFilter(file *entities.ScanSettingsFile, newFilter *entities.ComponentFilter, action entities.FilterAction) error {
	switch action {
	case entities.Include:
		file.Bom.Include = append(file.Bom.Include, *newFilter)
	case entities.Remove:
		file.Bom.Remove = append(file.Bom.Remove, *newFilter)
	default:
		return ErrInvalidFilterAction
	}

	return nil
}
