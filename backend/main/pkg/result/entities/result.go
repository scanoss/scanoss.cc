package entities

import (
	"errors"
	"strings"
)

type Result struct {
	Path      string
	MatchType string
	Purl      []string
}

const (
	NoMatch = "none"
)

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

//go:generate mockery --name ResultFilter --with-expecter
type ResultFilter interface {
	IsValid(result Result) bool
}

type Match struct {
	ID   string   `json:"id"`
	Purl []string `json:"purl,omitempty"`
}

var (
	ErrReadingResultFile = errors.New("error reading result file")
	ErrParsingResultFile = errors.New("error parsing result file")
)
