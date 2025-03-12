// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package repository

import (
	"errors"
	"sort"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/utils"
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
	fr                utils.FileReader
	resultsRepository ResultRepository
}

func NewJSONComponentRepository(fr utils.FileReader, resultsRepository ResultRepository) *JSONComponentRepository {
	return &JSONComponentRepository{
		fr:                fr,
		resultsRepository: resultsRepository,
	}
}

func (r *JSONComponentRepository) FindByFilePath(path string) (entities.Component, error) {
	result := r.resultsRepository.GetResultByPath(path)
	if result == nil {
		return entities.Component{}, errors.New("no result found")
	}

	components := result.Matches

	if len(components) == 0 {
		return entities.Component{}, errors.New("no components found")
	}

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
