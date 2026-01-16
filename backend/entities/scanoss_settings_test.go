//go:build unit

// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2026 SCANOSS.COM
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

package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentFilter_Priority(t *testing.T) {
	tests := []struct {
		name     string
		filter   ComponentFilter
		expected int
	}{
		{
			name:     "path + purl = score 4",
			filter:   ComponentFilter{Path: "src/file.js", Purl: "pkg:npm/lodash"},
			expected: 4,
		},
		{
			name:     "purl only = score 2",
			filter:   ComponentFilter{Purl: "pkg:npm/lodash"},
			expected: 2,
		},
		{
			name:     "path only (folder) = score 1",
			filter:   ComponentFilter{Path: "src/vendor/"},
			expected: 1,
		},
		{
			name:     "path only (file without purl) = score 1",
			filter:   ComponentFilter{Path: "src/file.js"},
			expected: 1,
		},
		{
			name:     "empty filter = score 0",
			filter:   ComponentFilter{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.Priority()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComponentFilter_AppliesTo(t *testing.T) {
	purl := []string{"pkg:npm/lodash@1.0.0"}

	tests := []struct {
		name     string
		filter   ComponentFilter
		result   Result
		expected bool
	}{
		{
			name:     "filter with no constraints applies to any result",
			filter:   ComponentFilter{},
			result:   Result{Path: "src/file.js", Purl: &purl},
			expected: true,
		},
		{
			name:     "path only filter applies to result with matching path",
			filter:   ComponentFilter{Path: "src/"},
			result:   Result{Path: "src/file.js", Purl: &purl},
			expected: true,
		},
		{
			name:     "purl only filter applies to result with matching purl",
			filter:   ComponentFilter{Purl: "pkg:npm/lodash@1.0.0"},
			result:   Result{Path: "lib/file.js", Purl: &purl},
			expected: true,
		},
		{
			name:     "path+purl filter requires both to apply",
			filter:   ComponentFilter{Path: "src/file.js", Purl: "pkg:npm/lodash@1.0.0"},
			result:   Result{Path: "src/file.js", Purl: &purl},
			expected: true,
		},
		{
			name:     "path+purl filter does not apply if path wrong",
			filter:   ComponentFilter{Path: "src/file.js", Purl: "pkg:npm/lodash@1.0.0"},
			result:   Result{Path: "lib/file.js", Purl: &purl},
			expected: false,
		},
		{
			name:     "path+purl filter does not apply if purl wrong",
			filter:   ComponentFilter{Path: "src/file.js", Purl: "pkg:npm/react@18.0.0"},
			result:   Result{Path: "src/file.js", Purl: &purl},
			expected: false,
		},
		{
			name:     "path only filter applies to result with nil purl",
			filter:   ComponentFilter{Path: "src/"},
			result:   Result{Path: "src/file.js", Purl: nil},
			expected: true,
		},
		{
			name:     "purl filter does not apply when result has nil purl",
			filter:   ComponentFilter{Purl: "pkg:npm/lodash@1.0.0"},
			result:   Result{Path: "src/file.js", Purl: nil},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.filter.AppliesTo(tt.result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestComponentFilter_Compare(t *testing.T) {
	tests := []struct {
		name     string
		r1       ComponentFilter
		r2       ComponentFilter
		expected int
	}{
		{
			name:     "path+purl beats purl only",
			r1:       ComponentFilter{Path: "src/file.js", Purl: "pkg:npm/lodash"},
			r2:       ComponentFilter{Purl: "pkg:npm/lodash"},
			expected: -1, // r1 has higher priority, we need first in the list
		},
		{
			name:     "purl only beats path only",
			r1:       ComponentFilter{Purl: "pkg:npm/lodash"},
			r2:       ComponentFilter{Path: "src/vendor/"},
			expected: -1,
		},
		{
			name:     "longer path wins at same score",
			r1:       ComponentFilter{Path: "src/vendor/deep/"},
			r2:       ComponentFilter{Path: "src/"},
			expected: -1, // r1 has longer path, so higher priority
		},
		{
			name:     "file path beats folder path at same score (longer path wins)",
			r1:       ComponentFilter{Path: "src/main.c"},
			r2:       ComponentFilter{Path: "src/"},
			expected: -1, // r1 is longer (10 chars vs 4 chars), so r1 comes first
		},
		{
			name:     "same path and score",
			r1:       ComponentFilter{Path: "src/vendor/"},
			r2:       ComponentFilter{Path: "src/vendor/"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.r1.Compare(tt.r2))
		})
	}
}

func TestIsResultInList_FolderMatching(t *testing.T) {
	sf := &SettingsFile{}
	purl := []string{"pkg:npm/lodash"}

	tests := []struct {
		name     string
		result   Result
		list     []ComponentFilter
		expected bool
		index    int
	}{
		{
			name: "folder rule matches file in folder",
			result: Result{
				Path: "src/vendor/lib.js",
				Purl: &purl,
			},
			list: []ComponentFilter{
				{Path: "src/vendor/"},
			},
			expected: true,
			index:    0,
		},
		{
			name: "folder rule matches nested file",
			result: Result{
				Path: "src/vendor/deep/nested/file.js",
				Purl: &purl,
			},
			list: []ComponentFilter{
				{Path: "src/vendor/"},
			},
			expected: true,
			index:    0,
		},
		{
			name: "folder rule does not match different path",
			result: Result{
				Path: "lib/vendor/file.js",
				Purl: &purl,
			},
			list: []ComponentFilter{
				{Path: "src/vendor/"},
			},
			expected: false,
			index:    -1,
		},
		{
			name: "folder rule does not match similar prefix without trailing slash",
			result: Result{
				Path: "src/vendorlib/file.js",
				Purl: &purl,
			},
			list: []ComponentFilter{
				{Path: "src/vendor/"},
			},
			expected: false,
			index:    -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, idx := sf.IsResultInList(tt.result, tt.list)
			assert.Equal(t, tt.expected, found)
			assert.Equal(t, tt.index, idx)
		})
	}
}

func TestIsResultInList_PrioritySystem(t *testing.T) {
	sf := &SettingsFile{}
	purl := []string{"pkg:npm/lodash@4.17.21"}

	t.Run("path+purl beats path-only folder rule", func(t *testing.T) {
		result := Result{
			Path: "src/main.c",
			Purl: &purl,
		}
		list := []ComponentFilter{
			{Path: "src/"},                      // Score 1 - folder
			{Path: "src/main.c", Purl: purl[0]}, // Score 4 - file + purl
		}

		found, idx := sf.IsResultInList(result, list)
		assert.True(t, found)
		assert.Equal(t, 1, idx, "Should match the path+purl rule (index 1)")
	})

	t.Run("longer path wins at same priority score", func(t *testing.T) {
		result := Result{
			Path: "src/vendor/lib.js",
			Purl: &purl,
		}
		list := []ComponentFilter{
			{Path: "src/"},        // Score 1, path length 4
			{Path: "src/vendor/"}, // Score 1, path length 11
		}

		found, idx := sf.IsResultInList(result, list)
		assert.True(t, found)
		assert.Equal(t, 1, idx, "Should match the longer path rule (src/vendor/)")
	})

	t.Run("purl beats path-only", func(t *testing.T) {
		result := Result{
			Path: "src/file.js",
			Purl: &purl,
		}
		list := []ComponentFilter{
			{Path: "src/"},  // Score 1
			{Purl: purl[0]}, // Score 2
		}

		found, idx := sf.IsResultInList(result, list)
		assert.True(t, found)
		assert.Equal(t, 1, idx, "Should match the purl rule (score 2 > score 1)")
	})

	t.Run("empty list returns false", func(t *testing.T) {
		result := Result{
			Path: "src/file.js",
			Purl: &purl,
		}

		found, idx := sf.IsResultInList(result, []ComponentFilter{})
		assert.False(t, found)
		assert.Equal(t, -1, idx)
	})
}

func TestGetResultFilterType(t *testing.T) {
	tests := []struct {
		name     string
		filter   ComponentFilter
		expected FilterType
	}{
		{
			name:     "folder rule (ends with /)",
			filter:   ComponentFilter{Path: "src/vendor/"},
			expected: ByFolder,
		},
		{
			name:     "file rule (path + purl)",
			filter:   ComponentFilter{Path: "src/file.js", Purl: "pkg:npm/lodash"},
			expected: ByFile,
		},
		{
			name:     "component rule (purl only)",
			filter:   ComponentFilter{Purl: "pkg:npm/lodash"},
			expected: ByPurl,
		},
		{
			name:     "folder with purl still detected as folder",
			filter:   ComponentFilter{Path: "src/vendor/", Purl: "pkg:npm/lodash"},
			expected: ByFolder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getResultFilterType(tt.filter)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetResultFilterConfig_WithFolders(t *testing.T) {
	purl := []string{"pkg:npm/lodash@1.0.0"}

	sf := &SettingsFile{
		Bom: Bom{
			Include: []ComponentFilter{
				{Path: "src/vendor/"},
			},
			Remove: []ComponentFilter{
				{Path: "tests/"},
			},
			Replace: []ComponentFilter{
				{Path: "lib/", ReplaceWith: "pkg:npm/new-lib"},
			},
		},
	}

	t.Run("include folder rule", func(t *testing.T) {
		result := Result{Path: "src/vendor/lib.js", Purl: &purl}
		config := sf.GetResultFilterConfig(result)
		assert.Equal(t, Include, config.Action)
		assert.Equal(t, ByFolder, config.Type)
	})

	t.Run("remove folder rule", func(t *testing.T) {
		result := Result{Path: "tests/unit/test.js", Purl: &purl}
		config := sf.GetResultFilterConfig(result)
		assert.Equal(t, Remove, config.Action)
		assert.Equal(t, ByFolder, config.Type)
	})

	t.Run("replace folder rule", func(t *testing.T) {
		result := Result{Path: "lib/utils.js", Purl: &purl}
		config := sf.GetResultFilterConfig(result)
		assert.Equal(t, Replace, config.Action)
		assert.Equal(t, ByFolder, config.Type)
	})
}
