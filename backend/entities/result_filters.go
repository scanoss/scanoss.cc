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

package entities

import "strings"

type ResultFilterAND struct {
	filters []ResultFilter
}

func NewResultFilterAND() *ResultFilterAND {
	return &ResultFilterAND{}
}

func (f *ResultFilterAND) AddFilter(filter ResultFilter) {
	f.filters = append(f.filters, filter)
}

func (f *ResultFilterAND) IsValid(result Result) bool {
	for _, filter := range f.filters {
		if !filter.IsValid(result) {
			return false
		}
	}
	return true
}

func (f *ResultFilterAND) HasFilters() bool {
	return len(f.filters) > 0
}

type ResultFilterMatchType struct {
	matchType string
}

type ResultQueryFilter struct {
	query string
}

func NewResultFilterMatchType(matchType string) *ResultFilterMatchType {
	return &ResultFilterMatchType{
		matchType: matchType,
	}
}

func NewResultQueryFilter(query string) *ResultQueryFilter {
	return &ResultQueryFilter{
		query: query,
	}
}

func (f *ResultFilterMatchType) IsValid(result Result) bool {
	return f.matchType == result.MatchType
}

func (f *ResultQueryFilter) IsValid(result Result) bool {
	queryLower := strings.ToLower(f.query)
	pathLower := strings.ToLower(result.Path)
	purlLower := strings.ToLower((*result.Purl)[0])

	return strings.Contains(pathLower, queryLower) || strings.Contains(purlLower, queryLower)
}

type ResultFilterFactory struct {
}

func NewResultFilterFactory() *ResultFilterFactory {
	return &ResultFilterFactory{}
}

func (f *ResultFilterFactory) Create(dto *RequestResultDTO) ResultFilter {
	if dto == nil {
		return nil
	}
	filterAND := NewResultFilterAND()

	if dto.MatchType != "" {
		filterAND.AddFilter(NewResultFilterMatchType(string(dto.MatchType)))
	}

	if dto.Query != "" {
		filterAND.AddFilter(NewResultQueryFilter(dto.Query))
	}

	return filterAND
}
