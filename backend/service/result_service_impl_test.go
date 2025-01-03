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

package service_test

import (
	"testing"

	"github.com/scanoss/scanoss.cc/backend/entities"
	mapperMocks "github.com/scanoss/scanoss.cc/backend/mappers/mocks"
	repoMocks "github.com/scanoss/scanoss.cc/backend/repository/mocks"
	"github.com/scanoss/scanoss.cc/backend/service"
	internal_test "github.com/scanoss/scanoss.cc/internal"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetResults(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	mu := internal_test.NewMockUtils()
	mu.On("ReadFile", config.GetInstance().ResultFilePath).Return([]byte(`{"path/to/file": [{"ID": "file", "Purl": ["pkg:example/package"]}]}`), nil)

	mockRepo := repoMocks.NewMockResultRepository(t)
	resultMapper := mapperMocks.NewMockResultMapper(t)
	mockRepo.EXPECT().GetResults(mock.AnythingOfType("*entities.ResultFilterAND")).Return([]entities.Result{
		{
			Path:      "path/to/file",
			MatchType: "file",
			Purl:      &[]string{"pkg:example/package"},
		},
	}, nil)

	resultMapper.EXPECT().MapToResultDTOList([]entities.Result{
		{
			Path:      "path/to/file",
			MatchType: "file",
			Purl:      &[]string{"pkg:example/package"},
		},
	}).Return([]entities.ResultDTO{{
		Path:         "path/to/file",
		MatchType:    "file",
		DetectedPurl: "pkg:example/package",
	}})

	service := service.NewResultServiceImpl(mockRepo, resultMapper)

	dto := &entities.RequestResultDTO{}

	results, err := service.GetAll(dto)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, results[0], entities.ResultDTO{
		Path:             "path/to/file",
		MatchType:        "file",
		ConcludedPurl:    "",
		WorkflowState:    "",
		FilterConfig:     entities.FilterConfig{},
		Comment:          "",
		DetectedPurl:     "pkg:example/package",
		ConcludedPurlUrl: "",
		ConcludedName:    "",
	})

	mockRepo.AssertExpectations(t)
	resultMapper.AssertExpectations(t)
}
