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
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type ScanossSettingsJsonRepository struct {
	fr    utils.FileReader
	mutex sync.RWMutex
}

func NewScanossSettingsJsonRepository(fr utils.FileReader) ScanossSettingsRepository {
	return &ScanossSettingsJsonRepository{
		fr: fr,
	}
}

func (r *ScanossSettingsJsonRepository) Init() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	sf, err := r.Read()
	if err != nil {
		log.Fatal().Err(err).Msgf("Error reading settings file: %v", config.GetInstance().ScanSettingsFilePath)
		return err
	}

	entities.ScanossSettingsJson = &entities.ScanossSettings{
		SettingsFile: &sf,
	}

	return nil
}

func (r *ScanossSettingsJsonRepository) Save() error {
	sf := r.GetSettings()
	if err := utils.WriteJsonFile(config.GetInstance().ScanSettingsFilePath, sf); err != nil {
		return err
	}
	return nil
}

func (r *ScanossSettingsJsonRepository) Read() (entities.SettingsFile, error) {
	if config.GetInstance() == nil {
		return entities.SettingsFile{}, fmt.Errorf("config is not initialized")
	}
	scanSettingsFileBytes, err := r.fr.ReadFile(config.GetInstance().ScanSettingsFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return entities.SettingsFile{}, nil
		}
		return entities.SettingsFile{}, err
	}

	scanossSettings, err := utils.JSONParse[entities.SettingsFile](scanSettingsFileBytes)
	if err != nil {
		return entities.SettingsFile{}, err
	}

	return scanossSettings, nil
}

func (r *ScanossSettingsJsonRepository) GetSettings() *entities.SettingsFile {
	return entities.ScanossSettingsJson.SettingsFile
}

func (r *ScanossSettingsJsonRepository) HasUnsavedChanges() (bool, error) {
	originalBom, err := r.Read()
	if err != nil {
		return false, err
	}

	currentContent := r.GetSettings()

	return originalBom.Equal(currentContent)
}

func (r *ScanossSettingsJsonRepository) AddBomEntry(newEntry entities.ComponentFilter, filterAction string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	sf := r.GetSettings()
	var targetList *[]entities.ComponentFilter

	switch filterAction {
	case "remove":
		targetList = &sf.Bom.Remove
	case "include":
		targetList = &sf.Bom.Include
	case "replace":
		targetList = &sf.Bom.Replace
	default:
		return fmt.Errorf("invalid filter action: %s", filterAction)
	}

	r.removeDuplicatesFromAllLists(newEntry)

	*targetList = append(*targetList, newEntry)

	return nil
}

func (r *ScanossSettingsJsonRepository) removeDuplicatesFromAllLists(newEntry entities.ComponentFilter) {
	sf := r.GetSettings()

	sf.Bom.Remove = removeDuplicatesFromList(sf.Bom.Remove, newEntry)
	sf.Bom.Include = removeDuplicatesFromList(sf.Bom.Include, newEntry)
	sf.Bom.Replace = removeDuplicatesFromList(sf.Bom.Replace, newEntry)
}

func removeDuplicatesFromList(list []entities.ComponentFilter, newEntry entities.ComponentFilter) []entities.ComponentFilter {
	result := make([]entities.ComponentFilter, 0, len(list))
	for _, entry := range list {
		if !isDuplicate(entry, newEntry) {
			result = append(result, entry)
		}
	}
	return result
}

func isDuplicate(entry, newEntry entities.ComponentFilter) bool {
	if newEntry.Path == "" {
		return entry.Purl == newEntry.Purl
	}
	return entry.Purl == newEntry.Purl && entry.Path == newEntry.Path && entry.ReplaceWith == newEntry.ReplaceWith && entry.License == newEntry.License
}

func (r *ScanossSettingsJsonRepository) ClearAllFilters() error {
	sf := r.GetSettings()
	sf.Bom.Include = []entities.ComponentFilter{}
	sf.Bom.Remove = []entities.ComponentFilter{}
	sf.Bom.Replace = []entities.ComponentFilter{}
	return nil
}

func (r *ScanossSettingsJsonRepository) GetDeclaredPurls() []string {
	sf := r.GetSettings()

	includedComponents := extractPurlsFromBom(sf.Bom.Include)
	removedComponents := extractPurlsFromBom(sf.Bom.Remove)
	replacedComponents := extractPurlsFromBom(sf.Bom.Replace)

	totalLength := len(includedComponents) + len(removedComponents) + len(replacedComponents)
	declaredPurls := make([]string, 0, totalLength)

	declaredPurls = append(declaredPurls, includedComponents...)
	declaredPurls = append(declaredPurls, removedComponents...)
	declaredPurls = append(declaredPurls, replacedComponents...)

	return declaredPurls
}

func extractPurlsFromBom(componentFilters []entities.ComponentFilter) []string {
	if len(componentFilters) == 0 {
		return []string{}
	}

	extractedPurls := make([]string, 0)

	for _, cf := range componentFilters {
		extractedPurls = append(extractedPurls, cf.Purl)
	}

	return extractedPurls
}
