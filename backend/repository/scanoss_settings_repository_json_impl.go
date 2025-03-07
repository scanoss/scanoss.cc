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
	"slices"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

type ScanossSettingsJsonRepository struct {
	fr                   utils.FileReader
	mutex                sync.RWMutex
	defaultSkipPatterns  []string
	stagedAddPatterns    []string
	stagedRemovePatterns []string
}

func NewScanossSettingsJsonRepository(fr utils.FileReader) ScanossSettingsRepository {
	return &ScanossSettingsJsonRepository{
		fr: fr,
	}
}

func (r *ScanossSettingsJsonRepository) Init() error {
	cfg := config.GetInstance()

	if cfg == nil {
		err := fmt.Errorf("config is not initialized")
		log.Error().Err(err).Msg("Error initializing ScanossSettingsJsonRepository")
		return err
	}

	r.defaultSkipPatterns = r.generateDefaultSkipPatterns()

	r.setSettingsFile(cfg.GetScanSettingsFilePath())

	cfg.RegisterListener(r.onConfigChange)

	return nil
}

// Triggered when the config is changed (e.g. scan root, scan settings file path, etc.)
func (r *ScanossSettingsJsonRepository) onConfigChange(newCfg *config.Config) {
	r.setSettingsFile(newCfg.GetScanSettingsFilePath())
}

func (r *ScanossSettingsJsonRepository) setSettingsFile(path string) {
	sf, err := r.Read()

	if err != nil {
		log.Error().Err(err).Msgf("Error reading settings file: %v", path)
		return
	}

	if entities.ScanossSettingsJson == nil {
		entities.ScanossSettingsJson = &entities.ScanossSettings{
			SettingsFile: &sf,
		}
	}

	entities.ScanossSettingsJson.SettingsFile = &sf
}

func (r *ScanossSettingsJsonRepository) Save() error {
	sf := r.GetSettings()

	if err := utils.WriteJsonFile(config.GetInstance().GetScanSettingsFilePath(), sf); err != nil {
		return err
	}

	return nil
}

func (r *ScanossSettingsJsonRepository) Read() (entities.SettingsFile, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	cfg := config.GetInstance()
	if cfg == nil {
		return entities.SettingsFile{}, fmt.Errorf("config is not initialized")
	}
	scanSettingsFileBytes, err := r.fr.ReadFile(cfg.GetScanSettingsFilePath())
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

func (r *ScanossSettingsJsonRepository) MatchesScanningSkipPattern(path string) bool {
	patterns := r.getScanningSkipPatterns()

	var matchers []gitignore.Pattern
	for _, pattern := range patterns {
		matcher := gitignore.ParsePattern(pattern, nil)
		matchers = append(matchers, matcher)
	}

	ps := gitignore.NewMatcher(matchers)

	isDir := false
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Debug().Err(err).Msgf("Error checking if path is directory: %s", path)
	} else {
		isDir = fileInfo.IsDir()
	}

	pathParts := strings.Split(path, "/")
	return ps.Match(pathParts, isDir)
}

func (r *ScanossSettingsJsonRepository) getScanningSkipPatterns() []string {
	sf := r.GetSettings()
	return append(r.defaultSkipPatterns, sf.Settings.Skip.Patterns.Scanning...)
}

func (r *ScanossSettingsJsonRepository) generateDefaultSkipPatterns() []string {
	defaultSkipPatterns := make([]string, 0, len(entities.DefaultSkippedDirExtensions)+len(entities.DefaultSkippedExtensions)+len(entities.DefaultSkippedDirs)+len(entities.DefaultSkippedFiles))

	for _, dirExtension := range entities.DefaultSkippedDirExtensions {
		defaultSkipPatterns = append(defaultSkipPatterns, fmt.Sprintf("*%s", dirExtension))
	}

	for _, extension := range entities.DefaultSkippedExtensions {
		defaultSkipPatterns = append(defaultSkipPatterns, fmt.Sprintf("*%s", extension))
	}

	for _, dir := range entities.DefaultSkippedDirs {
		defaultSkipPatterns = append(defaultSkipPatterns, fmt.Sprintf("%s/", dir))
	}

	defaultSkipPatterns = append(defaultSkipPatterns, entities.DefaultSkippedFiles...)

	return defaultSkipPatterns
}

func (r *ScanossSettingsJsonRepository) GetDefaultSkipPatterns(sf *entities.SettingsFile) []string {
	return r.defaultSkipPatterns
}

func (r *ScanossSettingsJsonRepository) isExcludedByBroaderPattern(path string) bool {
	effectivePatterns := r.GetEffectiveSkipPatterns()
	isExcluded := r.matchesPatterns(path, effectivePatterns)
	if !isExcluded {
		return false
	}

	exactMatch := slices.Contains(effectivePatterns, path)

	for _, pattern := range effectivePatterns {
		if strings.HasPrefix(pattern, "!") && strings.Contains(path, pattern[1:]) {
			log.Debug().Str("path", path).Str("negationPattern", pattern).Msg("Found potential negation pattern match")
		}
	}

	return !exactMatch // Excluded by a broader pattern
}

func (r *ScanossSettingsJsonRepository) matchesPatterns(path string, patterns []string) bool {
	var matchers []gitignore.Pattern
	for _, pattern := range patterns {
		matcher := gitignore.ParsePattern(pattern, nil)
		matchers = append(matchers, matcher)
	}

	ps := gitignore.NewMatcher(matchers)

	isDir := false
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Debug().Err(err).Msgf("Error checking if path is directory: %s", path)
	} else {
		isDir = fileInfo.IsDir()
	}

	pathParts := strings.Split(path, "/")
	return ps.Match(pathParts, isDir)
}

func (r *ScanossSettingsJsonRepository) AddStagedScanningSkipPattern(pattern string) error {
	sf := r.GetSettings()
	negationPattern := "!" + pattern
	hasNegationPattern := slices.Contains(sf.Settings.Skip.Patterns.Scanning, negationPattern)
	patternExistsInSettings := slices.Contains(sf.Settings.Skip.Patterns.Scanning, pattern)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if patternExistsInSettings {
		log.Debug().Msgf("Pattern already exists in settings file, skipping: %s", pattern)
		return nil
	}

	if slices.Contains(r.stagedAddPatterns, pattern) {
		log.Debug().Msgf("Pattern already staged for addition, skipping: %s", pattern)
		return nil
	}

	if hasNegationPattern {
		r.stagedRemovePatterns = append(r.stagedRemovePatterns, negationPattern)
		return nil
	}

	r.stagedAddPatterns = append(r.stagedAddPatterns, pattern)

	return nil
}

func (r *ScanossSettingsJsonRepository) RemoveStagedScanningSkipPattern(pattern string) error {
	isExcludedByBroader := r.isExcludedByBroaderPattern(pattern)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	sf := r.GetSettings()

	patternExistsInSettings := slices.Contains(sf.Settings.Skip.Patterns.Scanning, pattern)

	if patternExistsInSettings {
		r.stagedRemovePatterns = append(r.stagedRemovePatterns, pattern)
		return nil
	}

	if i := slices.Index(r.stagedAddPatterns, pattern); i >= 0 {
		r.stagedAddPatterns = slices.Delete(r.stagedAddPatterns, i, i+1)
		return nil
	}

	if isExcludedByBroader {
		negationPatterns := r.createNegationPattern(pattern)

		for _, negationPattern := range negationPatterns {
			if !slices.Contains(sf.Settings.Skip.Patterns.Scanning, negationPattern) &&
				!slices.Contains(r.stagedAddPatterns, negationPattern) {
				r.stagedAddPatterns = append(r.stagedAddPatterns, negationPattern)
			}
		}
	}

	return nil
}

func (r *ScanossSettingsJsonRepository) createNegationPattern(path string) []string {
	negationPatterns := []string{"!" + path}

	if strings.Contains(path, "/") {
		parts := strings.Split(path, "/")
		if len(parts) > 1 {
			filename := parts[len(parts)-1]
			negationPatterns = append(negationPatterns, "!"+filename)
			log.Debug().Str("path", path).Str("filenameNegation", "!"+filename).Msg("Added filename-only negation pattern")
		}
	}

	return negationPatterns
}

func (r *ScanossSettingsJsonRepository) CommitStagedSkipPatterns() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	sf := r.GetSettings()

	for _, pattern := range r.stagedAddPatterns {
		if !slices.Contains(sf.Settings.Skip.Patterns.Scanning, pattern) {
			sf.Settings.Skip.Patterns.Scanning = append(sf.Settings.Skip.Patterns.Scanning, pattern)
		}
	}

	for _, pattern := range r.stagedRemovePatterns {
		sf.Settings.Skip.Patterns.Scanning = slices.DeleteFunc(sf.Settings.Skip.Patterns.Scanning, func(p string) bool {
			return p == pattern
		})
	}

	r.stagedAddPatterns = nil
	r.stagedRemovePatterns = nil

	return r.Save()
}

func (r *ScanossSettingsJsonRepository) DiscardStagedSkipPatterns() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.stagedAddPatterns = nil
	r.stagedRemovePatterns = nil
	return nil
}

func (r *ScanossSettingsJsonRepository) GetEffectiveSkipPatterns() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	sf := r.GetSettings()

	effectivePatterns := make([]string, len(r.defaultSkipPatterns))
	copy(effectivePatterns, r.defaultSkipPatterns)

	for _, pattern := range sf.Settings.Skip.Patterns.Scanning {
		if !slices.Contains(r.stagedRemovePatterns, pattern) {
			effectivePatterns = append(effectivePatterns, pattern)
		}
	}

	for _, pattern := range r.stagedAddPatterns {
		if !slices.Contains(effectivePatterns, pattern) {
			effectivePatterns = append(effectivePatterns, pattern)
		}
	}

	return effectivePatterns
}

func (r *ScanossSettingsJsonRepository) MatchesEffectiveScanningSkipPattern(path string) bool {
	patterns := r.GetEffectiveSkipPatterns()

	var matchers []gitignore.Pattern
	for _, pattern := range patterns {
		matcher := gitignore.ParsePattern(pattern, nil)
		matchers = append(matchers, matcher)
	}

	ps := gitignore.NewMatcher(matchers)

	isDir := false
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Debug().Err(err).Msgf("Error checking if path is directory: %s", path)
	} else {
		isDir = fileInfo.IsDir()
	}

	pathParts := strings.Split(path, "/")
	return ps.Match(pathParts, isDir)
}
