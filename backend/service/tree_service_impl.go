// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2025 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package service

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/repository"
	"github.com/scanoss/scanoss.cc/internal/config"
)

type TreeServiceImpl struct {
	resultService       ResultService
	scanossSettingsRepo repository.ScanossSettingsRepository
	workerLimit         int
}

func NewTreeServiceImpl(resultService ResultService, scanossSettingsRepo repository.ScanossSettingsRepository) *TreeServiceImpl {
	maxWorkers := runtime.NumCPU() * 2
	return &TreeServiceImpl{
		resultService:       resultService,
		scanossSettingsRepo: scanossSettingsRepo,
		workerLimit:         maxWorkers,
	}
}

func (s *TreeServiceImpl) GetTree(rootPath string) ([]entities.TreeNode, error) {
	rootInfo, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}

	root := entities.NewTreeNode(rootInfo.Name(), entities.ResultDTO{}, rootInfo.IsDir())

	err = s.buildTree(rootPath, &root)
	if err != nil {
		return nil, err
	}

	return root.Children, nil
}

func (s *TreeServiceImpl) buildTree(path string, node *entities.TreeNode) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(entries))
	childrenChan := make(chan entities.TreeNode, len(entries))

	for _, entry := range entries {
		if isHiddenFileOrFolder(entry.Name()) {
			continue
		}

		wg.Add(1)
		go func(entry os.DirEntry) {
			defer wg.Done()

			absPath := filepath.Join(path, entry.Name())

			info, err := os.Lstat(absPath)
			if err != nil {
				errChan <- err
				return
			}

			if info.Mode()&os.ModeSymlink != 0 {
				return
			}

			scanRoot := config.GetInstance().GetScanRoot()
			resultRelativePath, err := filepath.Rel(scanRoot, absPath)
			if err != nil {
				errChan <- err
				return
			}

			result := s.resultService.GetByPath(resultRelativePath)
			childNode := entities.NewTreeNode(resultRelativePath, result, entry.IsDir())

			if entry.IsDir() {
				if err := s.buildTree(absPath, &childNode); err != nil {
					errChan <- err
					return
				}
				childNode.WorkflowState = calculateFolderWorkflowState(childNode.Children)
				childNode.ScanningSkipState = s.calculateFolderScanningSkipState(childNode)
			} else {
				childNode.ScanningSkipState = s.calculateScanningSkipState(resultRelativePath)
			}

			childrenChan <- childNode
		}(entry)
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(childrenChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	var children []entities.TreeNode
	for child := range childrenChan {
		children = append(children, child)
	}

	node.Children = children
	sortTreeNodes(node.Children)
	return nil
}

func (s *TreeServiceImpl) calculateScanningSkipState(path string) entities.SkipState {
	if s.scanossSettingsRepo.MatchesEffectiveScanningSkipPattern(path) {
		return entities.SkipStateExcluded
	}
	return entities.SkipStateIncluded
}

func (s *TreeServiceImpl) calculateFolderScanningSkipState(node entities.TreeNode) entities.SkipState {
	if s.scanossSettingsRepo.MatchesEffectiveScanningSkipPattern(node.Path) {
		return entities.SkipStateExcluded
	}

	excludedCount := 0
	includedCount := 0

	for _, child := range node.Children {
		if child.ScanningSkipState == entities.SkipStateMixed {
			return entities.SkipStateMixed
		}

		switch child.ScanningSkipState {
		case entities.SkipStateExcluded:
			excludedCount++
		case entities.SkipStateIncluded:
			includedCount++
		}
	}

	if excludedCount > 0 && includedCount > 0 {
		return entities.SkipStateMixed
	} else if excludedCount > 0 && includedCount == 0 {
		return entities.SkipStateExcluded
	} else {
		return entities.SkipStateIncluded
	}
}

func calculateFolderWorkflowState(nodes []entities.TreeNode) entities.WorkflowState {
	if len(nodes) == 0 {
		return entities.Pending
	}

	hasCompleted := false
	hasPending := false

	for _, node := range nodes {
		switch node.WorkflowState {
		case entities.Completed:
			hasCompleted = true
		case entities.Pending:
			hasPending = true
		case entities.Mixed:
			return entities.Mixed
		}

		if hasCompleted && hasPending {
			return entities.Mixed
		}
	}

	if hasCompleted && !hasPending {
		return entities.Completed
	}
	if hasPending && !hasCompleted {
		return entities.Pending
	}
	return entities.Mixed
}

func sortTreeNodes(nodes []entities.TreeNode) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].IsFolder == nodes[j].IsFolder {
			return nodes[i].Path < nodes[j].Path
		}
		return nodes[i].IsFolder
	})
}

func isHiddenFileOrFolder(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}
