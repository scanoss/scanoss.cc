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

package repository

import "github.com/scanoss/scanoss.cc/backend/entities"

type ResultRepositoryInMemoryImpl struct {
}

func NewResultRepositoryInMemoryImpl() *ResultRepositoryInMemoryImpl {
	return &ResultRepositoryInMemoryImpl{}
}

func (r *ResultRepositoryInMemoryImpl) GetResults(filter entities.ResultFilter) ([]entities.Result, error) {
	var results []entities.Result

	result1 := entities.Result{}
	result1.Path = "/inc/format_utils.h"
	result1.MatchType = "file"

	result2 := entities.Result{}
	result2.Path = "/src/format_utils.c"
	result2.MatchType = "snippet"

	result3 := entities.Result{}
	result3.Path = "/external/inc/crc32c.h"
	result3.MatchType = "none"

	results = append(results, result1)
	results = append(results, result2)
	results = append(results, result3)

	// Filter scan results
	if filter != nil {
		var filteredResults []entities.Result
		for _, result := range results {
			if filter.IsValid(result) {
				filteredResults = append(filteredResults, result)
			}
		}
		return filteredResults, nil
	}
	return results, nil
}
