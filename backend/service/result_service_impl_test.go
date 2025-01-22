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
