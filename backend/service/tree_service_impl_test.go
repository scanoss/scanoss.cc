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
	repositoryMocks "github.com/scanoss/scanoss.cc/backend/repository/mocks"
	"github.com/scanoss/scanoss.cc/backend/service"
	"github.com/scanoss/scanoss.cc/backend/service/mocks"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// patchConfig temporarily replaces the config's GetScanRoot function to return a custom path
func patchConfig(tmpDir string) func() {
	// Store the original function
	oldInstance := config.GetInstance()
	oldGetScanRoot := oldInstance.GetScanRoot

	// Create a patched version of GetScanRoot
	config.GetInstance().SetScanRoot(tmpDir)

	// Return a cleanup function
	return func() {
		config.GetInstance().SetScanRoot(oldGetScanRoot())
	}
}

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

	// Patch config to return our temporary directory
	unpatch := patchConfig(tmpDir)
	defer unpatch() // Restore original config after test

	mockScanossSettingsRepository := repositoryMocks.NewMockScanossSettingsRepository(t)
	mockScanossSettingsRepository.EXPECT().MatchesEffectiveScanningSkipPattern(mock.AnythingOfType("string")).Return(false).Times(7)

	mockResultService := mocks.NewMockResultService(t)
	mockResultService.EXPECT().GetByPath(mock.AnythingOfType("string")).Return(entities.ResultDTO{}).Times(7)

	service := service.NewTreeServiceImpl(mockResultService, mockScanossSettingsRepository)

	treeNodes, err := service.GetTree(tmpDir)

	assert.NoError(t, err)
	assert.NotNil(t, treeNodes)
	assert.Len(t, treeNodes, 3) // Expected 3 child nodes at the root level

	// Check for the specific folder and file names at root level
	foundFolder1 := false
	foundFolder2 := false
	foundRootJs := false

	for _, node := range treeNodes {
		if node.Name == "folder1" {
			foundFolder1 = true
			assert.True(t, node.IsFolder)
			assert.Len(t, node.Children, 2) // folder1 should have 2 children
		}
		if node.Name == "folder2" {
			foundFolder2 = true
			assert.True(t, node.IsFolder)
			assert.Len(t, node.Children, 1) // folder2 should have 1 child (subfolder)
		}
		if node.Name == "root.js" {
			foundRootJs = true
			assert.False(t, node.IsFolder)
			assert.Len(t, node.Children, 0) // root.js should have no children
		}
	}

	assert.True(t, foundFolder1, "folder1 not found in tree nodes")
	assert.True(t, foundFolder2, "folder2 not found in tree nodes")
	assert.True(t, foundRootJs, "root.js not found in tree nodes")

	mockResultService.AssertExpectations(t)
	mockScanossSettingsRepository.AssertExpectations(t)
}
