package entities

type TreeNodeState string

const (
	TreeNodeStatePending   TreeNodeState = "pending"
	TreeNodeStateCompleted TreeNodeState = "completed"
	TreeNodeStateMixed     TreeNodeState = "mixed"
)

type TreeNode struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Path       string        `json:"path"`
	IsFolder   bool          `json:"isFolder"`
	State      TreeNodeState `json:"state"`
	Result     *ResultDTO    `json:"data,omitempty"`
	ChildCount int           `json:"childCount"`
}

type TreeNodeResponse struct {
	Nodes      []TreeNode `json:"nodes"`
	TotalCount int        `json:"totalCount"`
}
