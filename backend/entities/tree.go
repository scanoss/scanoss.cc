package entities

import (
	"path/filepath"
)

type TreeNode struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Path              string        `json:"path"`
	IsFolder          bool          `json:"isFolder"`
	WorkflowState     WorkflowState `json:"workflowState"`
	ScanningSkipState SkipState     `json:"scanningSkipState"`
	Children          []TreeNode    `json:"children"`
}

type SkipState string

const (
	SkipStateIncluded SkipState = "included"
	SkipStateExcluded SkipState = "excluded"
	SkipStateMixed    SkipState = "mixed"
)

func NewTreeNode(path string, result ResultDTO, isFolder bool) TreeNode {
	return TreeNode{
		ID:                path,
		Name:              filepath.Base(path),
		Path:              path,
		IsFolder:          isFolder,
		WorkflowState:     result.WorkflowState,
		ScanningSkipState: SkipStateIncluded,
		Children:          make([]TreeNode, 0),
	}
}
