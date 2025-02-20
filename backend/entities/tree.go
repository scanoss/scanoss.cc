package entities

import (
	"path/filepath"
)

type Tree struct {
	Nodes []TreeNode
}

type TreeNode struct {
	ID            string
	Name          string
	Path          string
	IsFolder      bool
	WorkflowState WorkflowState
	Children      []TreeNode
}

func NewTreeNode(path string, result ResultDTO, isFolder bool) TreeNode {
	return TreeNode{
		ID:            path,
		Name:          filepath.Base(path),
		Path:          path,
		IsFolder:      isFolder,
		WorkflowState: result.WorkflowState,
		Children:      make([]TreeNode, 0),
	}
}
