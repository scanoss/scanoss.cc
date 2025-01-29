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
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type ResultRepositoryJsonImpl struct {
	fr utils.FileReader
}

func NewResultRepositoryJsonImpl(fr utils.FileReader) *ResultRepositoryJsonImpl {
	return &ResultRepositoryJsonImpl{
		fr: fr,
	}
}

func (r *ResultRepositoryJsonImpl) GetResults(filter entities.ResultFilter) ([]entities.Result, error) {
	resultFilePath := config.GetInstance().GetResultFilePath()
	resultByte, err := r.fr.ReadFile(resultFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []entities.Result{}, nil
		}
		log.Error().Err(err).Msg("Error reading result file")
		return []entities.Result{}, entities.ErrReadingResultFile
	}

	scanResults, err := r.parseScanResults(resultByte)
	if err != nil {
		return []entities.Result{}, entities.ErrParsingResultFile
	}

	if filter == nil {
		return scanResults, nil
	}

	// Filter scan results
	var filteredResults []entities.Result
	for _, result := range scanResults {
		if result.IsEmpty() {
			continue
		}

		if filter.IsValid(result) {
			filteredResults = append(filteredResults, result)
		}
	}
	return filteredResults, nil
}

func (r *ResultRepositoryJsonImpl) parseScanResults(resultByte []byte) ([]entities.Result, error) {
	var intermediateMap map[string][]entities.Match
	if err := json.Unmarshal(resultByte, &intermediateMap); err != nil {
		log.Error().Err(err).Msg("Error parsing scan results")
		return []entities.Result{}, err
	}

	var scanResults []entities.Result

	for key, matches := range intermediateMap {
		for _, match := range matches {
			scanResult := entities.Result{}
			scanResult.Path = key
			scanResult.MatchType = match.ID
			scanResult.Purl = &match.Purl
			scanResult.ComponentName = match.ComponentName
			scanResult.Matches = []entities.Match{match}
			scanResults = append(scanResults, scanResult)
		}
	}

	return scanResults, nil
}
