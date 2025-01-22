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
	"strings"
)

const (
	NoMatch = "none"
)

var (
	ErrReadingResultFile = errors.New("error reading result file")
	ErrParsingResultFile = errors.New("error parsing result file")
)

type Result struct {
	Path          string    `json:"path"`
	MatchType     string    `json:"match_type"`
	Purl          *[]string `json:"purl,omitempty"`
	ComponentName string    `json:"component"`
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) IsEmpty() bool {
	return r.MatchType == NoMatch
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

type ResultFilter interface {
	IsValid(result Result) bool
}

type Match struct {
	ID            string   `json:"id"`
	Purl          []string `json:"purl,omitempty"`
	ComponentName string   `json:"component"`
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
)

type FilterConfig struct {
	Action FilterAction `json:"action,omitempty"`
	Type   FilterType   `json:"type,omitempty"`
}

type FilterType string

const (
	ByFile FilterType = "by_file"
	ByPurl FilterType = "by_purl"
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
	MatchType MatchType `json:"match_type,omitempty" validate:"omitempty,eq=file|eq=snippet"`
	Query     string    `json:"query,omitempty"`
}
