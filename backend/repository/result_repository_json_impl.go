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

import (
	"encoding/json"

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
	// Path to your JSON file
	resultFilePath := config.GetInstance().ResultFilePath
	resultByte, err := r.fr.ReadFile(resultFilePath)
	if err != nil {
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
			scanResults = append(scanResults, scanResult)
		}
	}

	return scanResults, nil
}
