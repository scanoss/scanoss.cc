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
