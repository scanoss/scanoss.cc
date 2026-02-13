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

package mappers

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	purlutils "github.com/scanoss/go-purl-helper/pkg"
	"github.com/scanoss/scanoss.cc/backend/entities"
)

type ResultMapperImpl struct {
	scanossSettings *entities.ScanossSettings
}

var (
	resultDTOCache sync.Map
	purlCache      sync.Map
)

func NewResultMapper(scanossSettings *entities.ScanossSettings) ResultMapper {
	return &ResultMapperImpl{
		scanossSettings: scanossSettings,
	}
}

func (m *ResultMapperImpl) generateCacheKey(result entities.Result, bomEntry entities.ComponentFilter) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s-%s-%s",
		result.Path,
		strings.Join(*result.Purl, ","),
		result.MatchType,
		bomEntry.ReplaceWith,
		bomEntry.Comment,
		m.scanossSettings.SettingsFile.GetResultWorkflowState(result),
		m.scanossSettings.SettingsFile.GetResultFilterConfig(result),
	)
}

func (m *ResultMapperImpl) MapToResultDTO(result entities.Result) entities.ResultDTO {
	bomEntry := m.scanossSettings.SettingsFile.GetBomEntryFromResult(result)
	cacheKey := m.generateCacheKey(result, bomEntry)

	if cached, ok := resultDTOCache.Load(cacheKey); ok {
		return cached.(entities.ResultDTO)
	}

	var detectedPurl string
	if result.Purl != nil && len(*result.Purl) > 0 {
		detectedPurl = (*result.Purl)[0]
	}

	// If the scanner already applied the replacement, the detected purl will
	// equal the replace_with value. In that case, show the original purl from
	// the BOM entry so the display shows "original -> replacement".
	scannerPreAppliedReplacement := bomEntry.ReplaceWith != "" && detectedPurl == bomEntry.ReplaceWith && bomEntry.Purl != ""
	if scannerPreAppliedReplacement {
		detectedPurl = bomEntry.Purl
	}

	var detectedName string
	if scannerPreAppliedReplacement {
		detectedName = m.componentNameFromPurl(bomEntry.Purl)
	} else {
		detectedName = result.ComponentName
	}

	dto := entities.ResultDTO{
		MatchType:        entities.MatchType(result.MatchType),
		Path:             result.Path,
		DetectedPurl:     detectedPurl,
		DetectedPurlUrl:  m.mapPurlUrl(detectedPurl),
		DetectedName:     detectedName,
		ConcludedPurl:    bomEntry.ReplaceWith,
		ConcludedPurlUrl: m.mapPurlUrl(bomEntry.ReplaceWith),
		ConcludedName:    m.mapConcludedName(result),
		WorkflowState:    m.mapWorkflowState(result),
		FilterConfig:     m.mapFilterConfig(result),
		Comment:          bomEntry.Comment,
	}

	resultDTOCache.Store(cacheKey, dto)
	return dto
}


func (m *ResultMapperImpl) MapToResultDTOList(results []entities.Result) []entities.ResultDTO {
	output := make([]entities.ResultDTO, len(results))
	numWorkers := runtime.NumCPU()
	jobChan := make(chan int, len(results))
	var wg sync.WaitGroup

	for i := 0; i < len(results); i++ {
		jobChan <- i
	}
	close(jobChan)

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for idx := range jobChan {
				output[idx] = m.MapToResultDTO(results[idx])
			}
		}()
	}

	wg.Wait()
	return output
}

func (m *ResultMapperImpl) mapWorkflowState(result entities.Result) entities.WorkflowState {
	return m.scanossSettings.SettingsFile.GetResultWorkflowState(result)
}

func (m *ResultMapperImpl) mapFilterConfig(result entities.Result) entities.FilterConfig {
	return m.scanossSettings.SettingsFile.GetResultFilterConfig(result)
}


func (m *ResultMapperImpl) getProjectURL(purl string, compute func() (string, error)) string {
	if cached, ok := purlCache.Load(purl); ok {
		return cached.(string)
	}
	url, err := compute()
	if err != nil {
		log.Error().Err(err).Msg("Error computing project URL")
		return ""
	}
	purlCache.Store(purl, url)
	return url
}

func (m *ResultMapperImpl) mapPurlUrl(purl string) string {
	if purl == "" {
		return ""
	}

	purlObject, err := purlutils.PurlFromString(purl)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing purl")
		return ""
	}
	purlName := purlObject.Name
	if purlObject.Type == "github" && purlObject.Namespace != "" {
		purlName = fmt.Sprintf("%s/%s", purlObject.Namespace, purlObject.Name)
	}

	return m.getProjectURL(purl, func() (string, error) {
		return purlutils.ProjectUrl(purlName, purlObject.Type)
	})
}

func (m *ResultMapperImpl) mapConcludedName(result entities.Result) string {
	replacedPurl := m.scanossSettings.SettingsFile.GetBomEntryFromResult(result).ReplaceWith
	return m.componentNameFromPurl(replacedPurl)
}

func (m *ResultMapperImpl) componentNameFromPurl(purl string) string {
	if purl == "" {
		return ""
	}

	purlName, err := purlutils.PurlNameFromString(purl)
	if err != nil {
		log.Error().Err(err).Msg("Error getting component name from purl")
		return ""
	}

	nameParts := strings.Split(purlName, "/")
	componentName := nameParts[0]
	if len(nameParts) > 1 {
		componentName = nameParts[1]
	}

	return componentName
}
