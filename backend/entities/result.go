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
