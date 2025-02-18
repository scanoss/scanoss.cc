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

// Global cache variable (ensure proper concurrency control)
var purlCache sync.Map

func NewResultMapper(scanossSettings *entities.ScanossSettings) ResultMapper {
	return &ResultMapperImpl{
		scanossSettings: scanossSettings,
	}
}

func (m ResultMapperImpl) MapToResultDTO(result entities.Result) entities.ResultDTO {
	bomEntry := m.scanossSettings.SettingsFile.GetBomEntryFromResult(result)
	return entities.ResultDTO{
		MatchType:        entities.MatchType(result.MatchType),
		Path:             result.Path,
		DetectedPurl:     (*result.Purl)[0],
		DetectedPurlUrl:  m.mapDetectedPurlUrl(result),
		ConcludedPurl:    bomEntry.ReplaceWith,
		ConcludedPurlUrl: m.mapConcludedPurlUrl(result),
		ConcludedName:    m.mapConcludedName(result),
		WorkflowState:    m.mapWorkflowState(result),
		FilterConfig:     m.mapFilterConfig(result),
		Comment:          bomEntry.Comment,
	}
}

func (m ResultMapperImpl) mapConcludedPurl(result entities.Result) string {
	return m.scanossSettings.SettingsFile.GetBomEntryFromResult(result).ReplaceWith
}

func (m ResultMapperImpl) MapToResultDTOList(results []entities.Result) []entities.ResultDTO {
	output := make([]entities.ResultDTO, len(results))
	numWorkers := runtime.NumCPU() // or any number that fits your workload and machine
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

func (m *ResultMapperImpl) mapConcludedPurlUrl(result entities.Result) string {
	concludedPurl := m.mapConcludedPurl(result)
	if concludedPurl == "" {
		return ""
	}

	purlObject, err := purlutils.PurlFromString(concludedPurl)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing concluded purl")
		return ""
	}

	// Workaround for github purls until purlutils is updated
	purlName := purlObject.Name
	if purlObject.Type == "github" && purlObject.Namespace != "" {
		purlName = fmt.Sprintf("%s/%s", purlObject.Namespace, purlObject.Name)
	}

	return m.getProjectURL(concludedPurl, func() (string, error) {
		return purlutils.ProjectUrl(purlName, purlObject.Type)
	})
}

func (m ResultMapperImpl) getProjectURL(purl string, compute func() (string, error)) string {
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

func (m ResultMapperImpl) mapDetectedPurlUrl(result entities.Result) string {
	detectedPurl := (*result.Purl)[0]

	purlObject, err := purlutils.PurlFromString(detectedPurl)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing detected purl")
		return ""
	}
	purlName := purlObject.Name
	if purlObject.Type == "github" && purlObject.Namespace != "" {
		purlName = fmt.Sprintf("%s/%s", purlObject.Namespace, purlObject.Name)
	}

	return m.getProjectURL(detectedPurl, func() (string, error) {
		return purlutils.ProjectUrl(purlName, purlObject.Type)
	})
}

func (m ResultMapperImpl) mapConcludedName(result entities.Result) string {
	replacedPurl := m.scanossSettings.SettingsFile.GetBomEntryFromResult(result).ReplaceWith
	if replacedPurl == "" {
		return ""
	}

	purlName, err := purlutils.PurlNameFromString(replacedPurl)
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
