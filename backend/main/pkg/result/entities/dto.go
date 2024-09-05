package entities

type MatchType string

const (
	File    MatchType = "file"
	Snippet MatchType = "snippet"
)

type ResultDTO struct {
	File      string    `json:"file"`
	MatchType MatchType `json:"matchType"`
}

type RequestResultDTO struct {
	MatchType MatchType `json:"matchType,omitempty"`
}
