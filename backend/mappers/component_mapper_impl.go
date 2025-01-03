// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
