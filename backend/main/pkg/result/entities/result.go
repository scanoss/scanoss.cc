package entities

import "errors"

type Result struct {
	file      string
	matchType string
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) SetFile(fileName string) {
	r.file = fileName
}

func (r *Result) SetMatchType(matchType string) {
	r.matchType = matchType
}

func (r *Result) GetMatchType() string {
	return r.matchType
}

func (r *Result) GetFile() string {
	return r.file
}

func (r *Result) IsEmpty() bool {
	if r.matchType == "none" {
		return true
	}
	return false
}

type ResultFilter interface {
	IsValid(result Result) bool
}

type Match struct {
	ID string `json:"id"`
}

var (
	ErrReadingResultFile = errors.New("error reading result file")
	ErrParsingResultFile = errors.New("error parsing result file")
)
