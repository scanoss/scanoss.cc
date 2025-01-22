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
