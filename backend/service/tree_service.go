package service

import "github.com/scanoss/scanoss.cc/backend/entities"

type TreeService interface {
	GetChildren(parentPath string, offset, limit int) (*entities.TreeNodeResponse, error)
	GetRootNodes(offset, limit int) (*entities.TreeNodeResponse, error)
	SearchNodes(query string, offset, limit int) (*entities.TreeNodeResponse, error)
	GetNode(path string) (*entities.TreeNode, error)
}
