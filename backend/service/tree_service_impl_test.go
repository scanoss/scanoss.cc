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

package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/service"
	"github.com/scanoss/scanoss.cc/backend/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetTree(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "tree-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // Clean up after test

	// Create test directory structure
	dirs := []string{
		"folder1",
		"folder2",
		"folder2/subfolder",
	}

	files := []string{
		"root.js",
		"folder1/file1.js",
		"folder1/file2.js",
		"folder2/subfolder/file3.js",
	}

	// Create directories
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create empty files
	for _, file := range files {
		err := os.WriteFile(filepath.Join(tmpDir, file), []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	mockResultService := mocks.NewMockResultService(t)
	mockResultService.EXPECT().GetByPath(filepath.Join(tmpDir, "folder1")).Return(entities.ResultDTO{}, nil)
	mockResultService.EXPECT().GetByPath(filepath.Join(tmpDir, "folder2")).Return(entities.ResultDTO{}, nil)
	mockResultService.EXPECT().GetByPath(filepath.Join(tmpDir, "folder1/file1.js")).Return(entities.ResultDTO{}, nil)
	mockResultService.EXPECT().GetByPath(filepath.Join(tmpDir, "folder1/file2.js")).Return(entities.ResultDTO{}, nil)
	mockResultService.EXPECT().GetByPath(filepath.Join(tmpDir, "folder2/subfolder")).Return(entities.ResultDTO{}, nil)
	mockResultService.EXPECT().GetByPath(filepath.Join(tmpDir, "folder2/subfolder/file3.js")).Return(entities.ResultDTO{}, nil)
	mockResultService.EXPECT().GetByPath(filepath.Join(tmpDir, "root.js")).Return(entities.ResultDTO{}, nil)

	service := service.NewTreeServiceImpl(mockResultService)

	expectedTree := entities.Tree{
		Nodes: []entities.TreeNode{
			{
				ID:       filepath.Base(tmpDir),
				Name:     filepath.Base(tmpDir),
				Path:     filepath.Base(tmpDir),
				IsFolder: true,
				Children: []entities.TreeNode{
					{
						ID:       "folder1",
						Name:     "folder1",
						Path:     "folder1",
						IsFolder: true,
						Children: []entities.TreeNode{
							{
								ID:       "file1.js",
								Name:     "file1.js",
								Path:     "file1.js",
								IsFolder: false,
								Children: []entities.TreeNode{},
							},
							{
								ID:       "file2.js",
								Name:     "file2.js",
								Path:     "file2.js",
								IsFolder: false,
								Children: []entities.TreeNode{},
							},
						},
					},
					{
						ID:       "folder2",
						Name:     "folder2",
						Path:     "folder2",
						IsFolder: true,
						Children: []entities.TreeNode{
							{
								ID:       "subfolder",
								Name:     "subfolder",
								Path:     "subfolder",
								IsFolder: true,
								Children: []entities.TreeNode{
									{
										ID:       "file3.js",
										Name:     "file3.js",
										Path:     "file3.js",
										IsFolder: false,
										Children: []entities.TreeNode{},
									},
								},
							},
						},
					},
					{
						ID:       "root.js",
						Name:     "root.js",
						Path:     "root.js",
						IsFolder: false,
						Children: []entities.TreeNode{},
					},
				},
			},
		},
	}

	tree, err := service.GetTree(tmpDir)

	assert.NoError(t, err)
	assert.Equal(t, expectedTree, tree)
	mockResultService.AssertExpectations(t)
}
