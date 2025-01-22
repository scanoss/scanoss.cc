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

package mappers

import "github.com/scanoss/scanoss.cc/backend/entities"

type ComponentMapperImpl struct{}

func NewComponentMapper() ComponentMapper {
	return &ComponentMapperImpl{}
}

func (m *ComponentMapperImpl) MapToComponentDTO(componentEntity entities.Component) entities.ComponentDTO {
	dto := entities.ComponentDTO{
		ID:          componentEntity.ID,
		URL:         componentEntity.URL,
		Status:      componentEntity.Status,
		File:        componentEntity.File,
		Component:   componentEntity.Component,
		FileHash:    componentEntity.FileHash,
		Lines:       componentEntity.Lines,
		Purl:        componentEntity.Purl,
		FileURL:     componentEntity.FileURL,
		OssLines:    componentEntity.OssLines,
		Matched:     componentEntity.Matched,
		URLStats:    componentEntity.URLStats,
		Server:      componentEntity.Server,
		SourceHash:  componentEntity.SourceHash,
		URLHash:     componentEntity.URLHash,
		Vendor:      componentEntity.Vendor,
		Version:     componentEntity.Version,
		ReleaseDate: componentEntity.ReleaseDate,
		Provenance:  componentEntity.Provenance,
		Latest:      componentEntity.Latest,
	}

	var licenses []entities.ComponentLicense
	if len(componentEntity.Licenses) > 0 {
		for _, l := range componentEntity.Licenses {
			// patentHints ,copyleft0 checklistURL string, osadlUpdated string, source string, url string, incompatibleWith string
			licenses = append(licenses, entities.NewLicenseDTO(l.Name, l.PatentHints, l.Copyleft, l.ChecklistURL, l.OsadlUpdated, l.Source, l.URL, l.IncompatibleWith))
		}
		dto.Licenses = licenses
	}

	return dto
}
