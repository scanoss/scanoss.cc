package infraestructure

import (
	"errors"
	"integration-git/main/pkg/component/domain"
	"integration-git/main/pkg/utils"
	"sort"
)

var (
	ErrOpeningResult = errors.New("error opening result file")
	ErrReadingResult = errors.New("error reading result file")
)

type JSONComponentRepository struct {
	resultFilePath     string
	licenseSourceOrder map[string]int
}

func NewComponentRepository(resultFilePath string) *JSONComponentRepository {
	return &JSONComponentRepository{
		resultFilePath: resultFilePath,
		licenseSourceOrder: map[string]int{
			"component_declared": 1,
			"file_header":        2,
			"scancode":           3,
			"license_file":       4,
			"file_spdx_tag":      5,
		},
	}
}

func (r *JSONComponentRepository) FindByFilePath(path string) (domain.Component, error) {

	resultFileBytes, err := utils.ReadFile(r.resultFilePath)
	if err != nil {
		return domain.Component{}, err
	}
	results, err := utils.JSONParse[domain.Component](resultFileBytes)
	if err != nil {
		return domain.Component{}, err
	}

	// Gets components from results
	components := results[path]

	// Order component licenses by source
	if len(components[0].Licenses) > 0 {
		r.orderComponentLicensesBySourceType(&components[0])
	}

	return components[0], nil
}

func (r *JSONComponentRepository) orderComponentLicensesBySourceType(component *domain.Component) {
	sort.Slice(component.Licenses, func(i, j int) bool {
		return r.licenseSourceOrder[component.Licenses[i].Source] < r.licenseSourceOrder[component.Licenses[j].Source]
	})
}
