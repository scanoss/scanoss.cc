package entities

type MatchType string

const (
	File    MatchType = "file"
	Snippet MatchType = "snippet"
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

type FilterAction string

const (
	Include FilterAction = "include"
	Remove  FilterAction = "remove"
	Replace FilterAction = "replace"
)

type FilterType string

const (
	ByFile FilterType = "by_file"
	ByPurl FilterType = "by_purl"
)

type ResultPurl struct {
	// Original purl detected in results.json
	Detected string `json:"detected,omitempty"`
	// This is the purl that was replaced_with
	Concluded        string `json:"concluded,omitempty"`
	ConcludedPurlUrl string `json:"concluded_purl_url,omitempty"`
}

type ResultDTO struct {
	Path          string        `json:"path"`
	Purl          ResultPurl    `json:"purl,omitempty"`
	MatchType     MatchType     `json:"match_type"`
	WorkflowState WorkflowState `json:"workflow_state,omitempty"`
	FilterConfig  FilterConfig  `json:"filter_config,omitempty"`
	Comment       string        `json:"comment,omitempty"`
}

type RequestResultDTO struct {
	MatchType MatchType `json:"match_type,omitempty" validate:"omitempty,eq=file|eq=snippet"`
	Query     string    `json:"query,omitempty"`
}
