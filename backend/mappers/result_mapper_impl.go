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
	"strings"

	"github.com/rs/zerolog/log"
	purlutils "github.com/scanoss/go-purl-helper/pkg"
	"github.com/scanoss/scanoss.cc/backend/entities"
)

type ResultMapperImpl struct {
	scanossSettings *entities.ScanossSettings
}

func NewResultMapper(scanossSettings *entities.ScanossSettings) ResultMapper {
	return &ResultMapperImpl{
		scanossSettings: scanossSettings,
	}
}

func (m ResultMapperImpl) MapToResultDTO(result entities.Result) entities.ResultDTO {
	return entities.ResultDTO{
		MatchType:        entities.MatchType(result.MatchType),
		Path:             result.Path,
		DetectedPurl:     (*result.Purl)[0],
		DetectedPurlUrl:  m.mapDetectedPurlUrl(result),
		ConcludedPurl:    m.mapConcludedPurl(result),
		ConcludedPurlUrl: m.mapConcludedPurlUrl(result),
		ConcludedName:    m.mapConcludedName(result),
		WorkflowState:    m.mapWorkflowState(result),
		FilterConfig:     m.mapFilterConfig(result),
		Comment:          m.mapComment(result),
	}
}

func (m ResultMapperImpl) mapComment(result entities.Result) string {
	return m.scanossSettings.SettingsFile.GetBomEntryFromResult(result).Comment
}

func (m ResultMapperImpl) mapConcludedPurl(result entities.Result) string {
	return m.scanossSettings.SettingsFile.GetBomEntryFromResult(result).ReplaceWith
}

func (m ResultMapperImpl) MapToResultDTOList(results []entities.Result) []entities.ResultDTO {
	output := make([]entities.ResultDTO, len(results))

	for i, v := range results {
		output[i] = m.MapToResultDTO(v)
	}

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

	purlUrl, err := purlutils.ProjectUrl(purlObject.Name, purlObject.Type)
	if err != nil {
		log.Error().Err(err).Msg("Error getting project url")
		return ""
	}

	return purlUrl
}

func (m ResultMapperImpl) mapDetectedPurlUrl(result entities.Result) string {
	detectedPurl := (*result.Purl)[0]

	purlObject, err := purlutils.PurlFromString(detectedPurl)

	if err != nil {
		log.Error().Err(err).Msg("Error parsing detected purl")
		return ""
	}

	purlUrl, err := purlutils.ProjectUrl(purlObject.Name, purlObject.Type)
	if err != nil {
		log.Error().Err(err).Msg("Error getting detected purl url")
		return ""
	}

	return purlUrl
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
