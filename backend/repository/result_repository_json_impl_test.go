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

package repository_test

import (
	"testing"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/entities/mocks"
	"github.com/scanoss/scanoss.cc/backend/repository"
	internal_test "github.com/scanoss/scanoss.cc/internal"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetResults(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	t.Run("No filter", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().GetResultFilePath()).Return([]byte(`{"path/to/file": [{"ID": "file", "Purl": ["pkg:example/package"]}]}`), nil)

		repo, err := repository.NewResultRepositoryJsonImpl(mu, nil)
		results, err := repo.GetResults(nil)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "path/to/file", results[0].Path)
		assert.Equal(t, "file", results[0].MatchType)
		assert.Equal(t, &[]string{"pkg:example/package"}, results[0].Purl)
	})

	t.Run("With filter", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().GetResultFilePath()).Return([]byte(`{"path/to/file": [{"ID": "file", "Purl": ["pkg:example/package"]}]}`), nil)

		filter := mocks.MockResultFilter{}
		filter.EXPECT().IsValid(mock.Anything).Return(true)

		repo, err := repository.NewResultRepositoryJsonImpl(mu, nil)
		results, err := repo.GetResults(&filter)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "path/to/file", results[0].Path)
		assert.Equal(t, "file", results[0].MatchType)
		assert.Equal(t, []string{"pkg:example/package"}, (*results[0].Purl))
	})

	t.Run("Read file error", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().GetResultFilePath()).Return([]byte{}, entities.ErrReadingResultFile)

		_, err := repository.NewResultRepositoryJsonImpl(mu, nil)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrReadingResultFile, err)
		mu.AssertExpectations(t)
	})

	t.Run("Invalid json", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().GetResultFilePath()).Return([]byte(`invalid json`), nil)

		_, err := repository.NewResultRepositoryJsonImpl(mu, nil)

		assert.NoError(t, err)
		mu.AssertExpectations(t)
	})

	t.Run("Filter no match", func(t *testing.T) {
		mu := internal_test.NewMockUtils()
		mu.On("ReadFile", config.GetInstance().GetResultFilePath()).Return([]byte(`{"path/to/file": [{"ID": "file", "Purl": ["pkg:example/package"]}]}`), nil)

		filter := mocks.MockResultFilter{}
		filter.EXPECT().IsValid(mock.Anything).Return(false)

		repo, err := repository.NewResultRepositoryJsonImpl(mu, nil)
		results, err := repo.GetResults(&filter)

		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})
}
