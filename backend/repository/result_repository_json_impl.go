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
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type ResultRepositoryJsonImpl struct {
	fr           utils.FileReader
	cache        []entities.Result
	lastModified time.Time
	mutex        sync.RWMutex
	watcher      *fsnotify.Watcher
}

func NewResultRepositoryJsonImpl(fr utils.FileReader) *ResultRepositoryJsonImpl {
	repo := &ResultRepositoryJsonImpl{
		fr: fr,
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("Error creating watcher")
		return repo
	}

	resultFilePath := config.GetInstance().GetResultFilePath()
	if err := watcher.Add(resultFilePath); err != nil {
		log.Error().Err(err).Msgf("Error watching file: %s", resultFilePath)
		watcher.Close()
		return repo
	}

	repo.watcher = watcher

	// Initial cache load
	if err := repo.refreshCache(); err != nil {
		log.Error().Err(err).Msg("Error loading initial cache")
	}

	config.GetInstance().RegisterListener(repo.onConfigChange)

	go repo.watchForChanges()

	return repo
}

func (r *ResultRepositoryJsonImpl) onConfigChange(newCfg *config.Config) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := r.refreshCache(); err != nil {
		log.Error().Err(err).Msg("Error refreshing results cache after config change")
	}
}

func (r *ResultRepositoryJsonImpl) watchForChanges() {
	for {
		select {
		case event, ok := <-r.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				r.mutex.Lock()
				if err := r.refreshCache(); err != nil {
					log.Error().Err(err).Msg("Error refreshing cache after file change")
				}
				r.mutex.Unlock()
			}
		case err, ok := <-r.watcher.Errors:
			if !ok {
				return
			}
			log.Error().Err(err).Msg("Watcher error")
		}
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
		return entities.ErrParsingResultFile
	}

	fileInfo, err := os.Stat(resultFilePath)
	if err == nil {
		r.lastModified = fileInfo.ModTime()
	}

	r.cache = scanResults
	return nil
}

func (r *ResultRepositoryJsonImpl) parseScanResults(resultByte []byte) ([]entities.Result, error) {
	var intermediateMap map[string][]entities.Component
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
			scanResult.ComponentName = match.Component
			scanResult.Matches = []entities.Component{match}
			scanResults = append(scanResults, scanResult)
		}
	}

	return scanResults, nil
}

func (r *ResultRepositoryJsonImpl) GetResultByPath(path string) (entities.Result, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, result := range r.cache {
		if result.Path == path {
			return result, nil
		}
	}

	return entities.Result{}, fmt.Errorf("result not found: %s", path)
}

func (r *ResultRepositoryJsonImpl) Close() {
	if r.watcher != nil {
		r.watcher.Close()
	}
}
