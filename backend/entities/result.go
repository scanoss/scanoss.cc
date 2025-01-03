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
	ConcludedPurl    string        `json:"concluded_purl,omitempty"`
	ConcludedPurlUrl string        `json:"concluded_purl_url,omitempty"`
	ConcludedName    string        `json:"concluded_name,omitempty"`
}

type RequestResultDTO struct {
	MatchType MatchType `json:"match_type,omitempty" validate:"omitempty,eq=file|eq=snippet"`
	Query     string    `json:"query,omitempty"`
}
