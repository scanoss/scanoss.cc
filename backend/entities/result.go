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

package entities

import (
	"errors"
	"strconv"
	"strings"
)

const (
	MatchTypeNone       = "none"
	MatchTypeDependency = "dependency"
)

var (
	ErrReadingResultFile = errors.New("error reading result file")
	ErrParsingResultFile = errors.New("error parsing result file")
)

type Result struct {
	Path          string      `json:"path"`
	MatchType     string      `json:"match_type"`
	// TODO: Consider changing Purl from *[]string to []string.
	// The pointer adds unnecessary complexity (nil checks everywhere).
	// []string with omitempty already omits empty slices from JSON.
	Purl          *[]string   `json:"purl,omitempty"`
	ComponentName string      `json:"component"`
	Matches       []Component `json:"matches,omitempty"`
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) IsEmpty() bool {
	return r.MatchType == MatchTypeNone
}

func (r *Result) IsDependency() bool {
	return r.MatchType == MatchTypeDependency
}

func (r *Result) IsValid() bool {
	return r.Path != "" && r.MatchType != ""
}

func (r *Result) GetFileName() string {
	parts := strings.Split(r.Path, "/")
	i := len(parts) - 1
	fileName := parts[i]

	return fileName
}

func (r *Result) GetMatchPercentage() float64 {
	if r.MatchType == string(MatchTypeFile) {
		return 100.0
	}

	// For snippet matches, parse the percentage from the matched field
	// Default to 0 if not found or invalid
	matches := r.Matches
	if len(matches) > 0 {
		matchStr := matches[0].Matched
		if matchStr != "" {
			matchStr = strings.TrimSuffix(matchStr, "%")
			if percentage, err := strconv.ParseFloat(matchStr, 64); err == nil {
				return percentage
			}
		}
	}
	return 0.0
}

type ResultFilter interface {
	IsValid(result Result) bool
}

type ResultLicense struct {
	Name   string `json:"name"`
	Source string `json:"source,omitempty"`
	URL    string `json:"url,omitempty"`
}

type MatchType string

const (
	MatchTypeFile    MatchType = "file"
	MatchTypeSnippet MatchType = "snippet"
)

type WorkflowState string

const (
	Pending   WorkflowState = "pending"
	Completed WorkflowState = "completed"
	Mixed     WorkflowState = "mixed"
)

type FilterConfig struct {
	Action FilterAction `json:"action,omitempty"`
	Type   FilterType   `json:"type,omitempty"`
}

type FilterType string

const (
	ByFile   FilterType = "by_file"
	ByPurl   FilterType = "by_purl"
	ByFolder FilterType = "by_folder"
)

type ResultDTO struct {
	Path             string        `json:"path"`
	MatchType        MatchType     `json:"match_type"`
	WorkflowState    WorkflowState `json:"workflow_state,omitempty"`
	FilterConfig     FilterConfig  `json:"filter_config,omitempty"`
	Comment          string        `json:"comment,omitempty"`
	DetectedPurl     string        `json:"detected_purl,omitempty"`
	DetectedPurlUrl  string        `json:"detected_purl_url,omitempty"`
	ConcludedPurl    string        `json:"concluded_purl,omitempty"`
	ConcludedPurlUrl string        `json:"concluded_purl_url,omitempty"`
	ConcludedName    string        `json:"concluded_name,omitempty"`
}

type RequestResultDTO struct {
	MatchType MatchType  `json:"match_type,omitempty" validate:"omitempty,eq=file|eq=snippet"`
	Query     string     `json:"query,omitempty"`
	Sort      SortConfig `json:"sort,omitempty" validate:"dive"`
}
