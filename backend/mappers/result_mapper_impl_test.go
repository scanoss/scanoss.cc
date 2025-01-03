//go:build unit

// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
