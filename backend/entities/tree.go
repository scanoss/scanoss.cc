package entities

import (
	"path/filepath"

	"github.com/google/uuid"
)

type TreeNode struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Path          string        `json:"path"`
	IsFolder      bool          `json:"isFolder"`
	WorkflowState WorkflowState `json:"workflowState"`
	Children      []TreeNode    `json:"children"`
}

func NewTreeNode(path string, result ResultDTO, isFolder bool) TreeNode {
	return TreeNode{
		ID:            uuid.New().String(),
		Name:          filepath.Base(path),
		Path:          path,
		IsFolder:      isFolder,
		WorkflowState: result.WorkflowState,
		Children:      make([]TreeNode, 0),
	}
}
