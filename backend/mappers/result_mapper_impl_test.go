//go:build unit

// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
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

package mappers_test

import (
	"testing"

	"github.com/scanoss/scanoss.cc/backend/entities"
	mappers "github.com/scanoss/scanoss.cc/backend/mappers"
	internal_test "github.com/scanoss/scanoss.cc/internal"
	"github.com/stretchr/testify/assert"
)

func TestMapToResultDTO(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	testCases := []struct {
		result   entities.Result
		expected entities.ResultDTO
	}{
		{
			result: entities.Result{
				MatchType: "file",
				Path:      "path/to/file",
				Purl:      &[]string{"pkg:purl-test"},
			},
			expected: entities.ResultDTO{
				Path:          "path/to/file",
				MatchType:     "file",
				WorkflowState: entities.Completed,
				FilterConfig: entities.FilterConfig{
					Action: entities.Include,
					Type:   entities.ByFile,
				},
			},
		},
		{
			result: entities.Result{
				Path:      "another/path/to/file",
				MatchType: "file",
				Purl:      &[]string{"pkg:another-purl-test"},
			},
			expected: entities.ResultDTO{
				Path:          "another/path/to/file",
				MatchType:     "file",
				WorkflowState: entities.Pending,
				FilterConfig:  entities.FilterConfig{},
			},
		},
		{
			result: entities.Result{
				Path:      "path/to/removed/file",
				MatchType: "snippet",
				Purl:      &[]string{"pkg:removed-file"},
			},
			expected: entities.ResultDTO{
				Path:          "path/to/removed/file",
				MatchType:     "snippet",
				WorkflowState: entities.Completed,
				FilterConfig: entities.FilterConfig{
					Action: entities.Remove,
					Type:   entities.ByFile,
				},
			},
		},
	}

	settings := &entities.ScanossSettings{
		SettingsFile: &entities.SettingsFile{
			Bom: entities.Bom{
				Include: []entities.ComponentFilter{
					{
						Path: "path/to/file",
						Purl: "pkg:purl-test",
					},
				},
				Remove: []entities.ComponentFilter{{
					Purl: "pkg:purl-removed",
				}},
			},
		},
	}

	mapper := mappers.NewResultMapper(settings)

	for _, tc := range testCases {
		dto := mapper.MapToResultDTO(tc.result)
		assert.Equal(t, tc.expected.MatchType, dto.MatchType)
		assert.Equal(t, tc.expected.Path, dto.Path)
	}
}
