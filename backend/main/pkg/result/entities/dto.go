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
	Action FilterAction `json:"action"`
	Type   FilterType   `json:"type"`
}

type FilterAction string

const (
	Include FilterAction = "include"
	Remove  FilterAction = "remove"
)

type FilterType string

const (
	ByFile FilterType = "by_file"
	ByPurl FilterType = "by_purl"
)

type ResultDTO struct {
	Path          string        `json:"path"`
	MatchType     MatchType     `json:"match_type"`
	WorkflowState WorkflowState `json:"workflow_state,omitempty"`
	FilterConfig  FilterConfig  `json:"filter_config,omitempty"`
}

type RequestResultDTO struct {
	MatchType MatchType `json:"match_type,omitempty"`
}
