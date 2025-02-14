package service

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/scanoss/scanoss.cc/backend/entities"
)

type TreeServiceImpl struct {
	resultService ResultService
}

func NewTreeService(resultService ResultService) TreeService {
	return &TreeServiceImpl{
		resultService: resultService,
	}
}

func (s *TreeServiceImpl) GetChildren(parentPath string, offset, limit int) (*entities.TreeNodeResponse, error) {
	// Get results with path prefix filter
	results, err := s.resultService.GetAll(&entities.RequestResultDTO{
		Query: parentPath + "/", // This should be handled by the result service to filter by path prefix
	})
	if err != nil {
		return nil, err
	}

	// Use a map to efficiently track unique directories
	dirMap := make(map[string]bool)
	var children []entities.TreeNode
	var totalCount int

	for _, result := range results {
		// Skip if not a direct child
		if !strings.HasPrefix(result.Path, parentPath+"/") {
			continue
		}

		// Get the relative path from parent
		relPath := strings.TrimPrefix(result.Path, parentPath+"/")
		parts := strings.Split(relPath, "/")

		if len(parts) == 1 {
			// Direct file child
			totalCount++
			if totalCount > offset && len(children) < limit {
				children = append(children, s.createNode(result.Path, result))
			}
		} else {
			// Directory child
			dirPath := filepath.Join(parentPath, parts[0])
			if !dirMap[dirPath] {
				dirMap[dirPath] = true
				totalCount++
				if totalCount > offset && len(children) < limit {
					children = append(children, s.createFolderNode(dirPath))
				}
			}
		}
	}

	return &entities.TreeNodeResponse{
		Nodes:      children,
		TotalCount: totalCount,
	}, nil
}

func (s *TreeServiceImpl) GetRootNodes(offset, limit int) (*entities.TreeNodeResponse, error) {
	results, err := s.resultService.GetAll(&entities.RequestResultDTO{})
	if err != nil {
		return nil, err
	}

	// Use a map to track unique root paths and their states
	rootMap := make(map[string]*entities.TreeNode)
	var totalCount int

	for _, result := range results {
		parts := strings.Split(result.Path, "/")
		if len(parts) == 1 {
			// Root file
			if _, exists := rootMap[result.Path]; !exists {
				node := s.createNode(result.Path, result)
				rootMap[result.Path] = &node
				totalCount++
			}
		} else {
			// Root directory
			rootPath := parts[0]
			if _, exists := rootMap[rootPath]; !exists {
				node := s.createFolderNode(rootPath)
				rootMap[rootPath] = &node
				totalCount++
			}
		}
	}

	// Convert map to slice with pagination
	var rootNodes []entities.TreeNode
	count := 0
	for _, node := range rootMap {
		if count >= offset && len(rootNodes) < limit {
			rootNodes = append(rootNodes, *node)
		}
		count++
	}

	return &entities.TreeNodeResponse{
		Nodes:      rootNodes,
		TotalCount: totalCount,
	}, nil
}

func (s *TreeServiceImpl) SearchNodes(query string, offset, limit int) (*entities.TreeNodeResponse, error) {
	results, err := s.resultService.GetAll(&entities.RequestResultDTO{Query: query})
	if err != nil {
		return nil, err
	}

	// Use a map to track unique paths and their states
	nodeMap := make(map[string]*entities.TreeNode)
	var totalCount int

	for _, result := range results {
		// Add the file node
		if _, exists := nodeMap[result.Path]; !exists {
			node := s.createNode(result.Path, result)
			nodeMap[result.Path] = &node
			totalCount++
		}

		// Add parent directory nodes
		dir := filepath.Dir(result.Path)
		for dir != "." {
			if _, exists := nodeMap[dir]; !exists {
				node := s.createFolderNode(dir)
				nodeMap[dir] = &node
				totalCount++
			}
			dir = filepath.Dir(dir)
		}
	}

	// Convert map to slice with pagination
	var nodes []entities.TreeNode
	count := 0
	for _, node := range nodeMap {
		if count >= offset && len(nodes) < limit {
			nodes = append(nodes, *node)
		}
		count++
	}

	return &entities.TreeNodeResponse{
		Nodes:      nodes,
		TotalCount: totalCount,
	}, nil
}

func (s *TreeServiceImpl) GetNode(path string) (*entities.TreeNode, error) {
	// Try to get it as a file first
	results, err := s.resultService.GetAll(&entities.RequestResultDTO{Query: path})
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		if result.Path == path {
			node := s.createNode(path, result)
			return &node, nil
		}
	}

	// If not found as a file, check if it's a directory
	dirResults, err := s.resultService.GetAll(&entities.RequestResultDTO{Query: path + "/"})
	if err != nil {
		return nil, err
	}

	if len(dirResults) > 0 {
		node := s.createFolderNode(path)
		return &node, nil
	}

	return nil, fmt.Errorf("path not found: %s", path)
}

func (s *TreeServiceImpl) createNode(path string, result entities.ResultDTO) entities.TreeNode {
	return entities.TreeNode{
		ID:         path,
		Name:       filepath.Base(path),
		Path:       path,
		IsFolder:   false,
		State:      entities.TreeNodeState(result.WorkflowState),
		Result:     &result,
		ChildCount: 0,
	}
}

func (s *TreeServiceImpl) createFolderNode(path string) entities.TreeNode {
	// Get child count
	results, err := s.resultService.GetAll(&entities.RequestResultDTO{Query: path + "/"})
	childCount := 0
	if err == nil {
		for _, result := range results {
			if strings.HasPrefix(result.Path, path+"/") {
				childCount++
			}
		}
	}

	return entities.TreeNode{
		ID:         path,
		Name:       filepath.Base(path),
		Path:       path,
		IsFolder:   true,
		State:      s.computeFolderState(path),
		ChildCount: childCount,
	}
}

func (s *TreeServiceImpl) computeFolderState(path string) entities.TreeNodeState {
	results, err := s.resultService.GetAll(&entities.RequestResultDTO{Query: path + "/"})
	if err != nil {
		return entities.TreeNodeStateMixed
	}

	hasPending := false
	hasCompleted := false

	for _, result := range results {
		if strings.HasPrefix(result.Path, path+"/") {
			if result.WorkflowState == "pending" {
				hasPending = true
			} else {
				hasCompleted = true
			}
			if hasPending && hasCompleted {
				return entities.TreeNodeStateMixed
			}
		}
	}

	if hasPending {
		return entities.TreeNodeStatePending
	}
	return entities.TreeNodeStateCompleted
}
