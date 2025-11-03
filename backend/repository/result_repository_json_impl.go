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
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type ResultRepositoryJsonImpl struct {
	fr           utils.FileReader
	cache        []entities.Result
	pathIndex    map[string]entities.Result
	lastModified time.Time
	mutex        sync.RWMutex
}

func NewResultRepositoryJsonImpl(fr utils.FileReader) (*ResultRepositoryJsonImpl, error) {
	repo := &ResultRepositoryJsonImpl{
		fr: fr,
	}

	// Initial cache load
	if err := repo.refreshCache(); err != nil {
		log.Error().Err(err).Msg("Error loading initial cache")
		return repo, err
	}

	config.GetInstance().RegisterListener(repo.onConfigChange)

	return repo, nil
}

func (r *ResultRepositoryJsonImpl) onConfigChange(newCfg *config.Config) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.refreshCache(); err != nil {
		log.Error().Err(err).Msg("Error refreshing results cache after config change")
	}
}

func (r *ResultRepositoryJsonImpl) GetResults(filter entities.ResultFilter) ([]entities.Result, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if filter == nil {
		return r.cache, nil
	}

	var filteredResults []entities.Result
	for _, result := range r.cache {
		if result.IsEmpty() {
			continue
		}
		if filter.IsValid(result) {
			filteredResults = append(filteredResults, result)
		}
	}

	return filteredResults, nil
}

func (r *ResultRepositoryJsonImpl) refreshCache() error {
	resultFilePath := config.GetInstance().GetResultFilePath()
	resultByte, err := r.fr.ReadFile(resultFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			r.cache = []entities.Result{}
			return nil
		}
		return entities.ErrReadingResultFile
	}

	scanResults, err := r.parseScanResults(resultByte)
	if err != nil {
		if os.IsNotExist(err) {
			r.cache = []entities.Result{}
			return nil
		}
		return entities.ErrParsingResultFile
	}

	fileInfo, err := os.Stat(resultFilePath)
	if err == nil {
		r.lastModified = fileInfo.ModTime()
	}

	r.cache = scanResults
	r.pathIndex = make(map[string]entities.Result)
	for _, result := range scanResults {
		r.pathIndex[result.Path] = result
	}
	return nil
}

func (r *ResultRepositoryJsonImpl) parseScanResults(resultByte []byte) ([]entities.Result, error) {
	var intermediateMap map[string][]entities.Component
	err := json.Unmarshal(resultByte, &intermediateMap)
	if err != nil {
		// Gracefully handle JSON syntax errors
		if typeError, ok := err.(*json.SyntaxError); ok {
			log.Error().Err(typeError).Msgf("JSON file %s syntax error: offset %d", config.GetInstance().GetResultFilePath(), typeError.Offset)
			return []entities.Result{}, nil
		}
		log.Error().Err(err).Msg("Error parsing scan results")
		return []entities.Result{}, err
	}

	var scanResults []entities.Result
	for path, components := range intermediateMap {
		// Create a single Result for each path with all its components
		result := entities.Result{
			Path:    path,
			Matches: components,
		}

		// If there are components, set the first component's details as the main result details
		if len(components) > 0 {
			result.MatchType = components[0].ID
			result.ComponentName = components[0].Component
			result.Purl = &components[0].Purl
		}

		scanResults = append(scanResults, result)
	}

	return scanResults, nil
}

func (r *ResultRepositoryJsonImpl) GetResultByPath(path string) *entities.Result {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result, ok := r.pathIndex[path]
	if !ok {
		return nil
	}

	return &result
}
