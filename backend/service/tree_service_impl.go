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
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/scanoss/scanoss.cc/backend/entities"
)

type TreeServiceImpl struct {
	resultService ResultService
}

func NewTreeServiceImpl(resultService ResultService) *TreeServiceImpl {
	return &TreeServiceImpl{resultService: resultService}
}

func (s *TreeServiceImpl) GetTree(rootPath string) (entities.Tree, error) {
	rootInfo, err := os.Stat(rootPath)
	if err != nil {
		return entities.Tree{}, err
	}

	root := entities.NewTreeNode(rootInfo.Name(), entities.ResultDTO{}, rootInfo.IsDir())

	err = s.buildTree(rootPath, &root)
	if err != nil {
		return entities.Tree{}, err
	}

	return entities.Tree{Nodes: []entities.TreeNode{root}}, nil
}

func (s *TreeServiceImpl) buildTree(path string, node *entities.TreeNode) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if isHiddenFileOrFolder(entry.Name()) {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		result := s.resultService.GetByPath(fullPath)

		log.Info().Msgf("result: %+v", result)

		childNode := entities.NewTreeNode(entry.Name(), result, entry.IsDir())

		if entry.IsDir() {
			err := s.buildTree(fullPath, &childNode)
			if err != nil {
				return err
			}
		}

		node.Children = append(node.Children, childNode)
	}

	sortTreeNodes(node.Children)

	return nil
}

// sortTreeNodes sorts tree nodes by path name
func sortTreeNodes(nodes []entities.TreeNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Path < nodes[j].Path
	})
}

func isHiddenFileOrFolder(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}
