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
	"errors"
	"testing"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/repository/mocks"
	"github.com/scanoss/scanoss.cc/backend/service"
	"github.com/stretchr/testify/assert"
)

func TestGetLocalFileContent(t *testing.T) {
	mockFileRepo := mocks.NewMockFileRepository(t)
	mockComponentRepo := mocks.NewMockComponentRepository(t)
	service := service.NewFileService(mockFileRepo, mockComponentRepo)

	expectedFile := entities.NewFile(
		"",
		"test.js",
		[]byte("function main() {\n\tconsole.log('Hello, World!');\n}"),
	)
	mockFileRepo.EXPECT().ReadLocalFile("test.js").Return(*expectedFile, nil)
	file, err := service.GetLocalFile("test.js")

	assert.NoError(t, err)
	assert.Equal(t, entities.FileDTO{
		Path:     "test.js",
		Name:     "test.js",
		Content:  "function main() {\n\tconsole.log('Hello, World!');\n}",
		Language: "javascript",
	}, file)
	mockFileRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}

func TestGetLocalFileContent_Error(t *testing.T) {
	mockFileRepo := mocks.NewMockFileRepository(t)
	mockComponentRepo := mocks.NewMockComponentRepository(t)
	service := service.NewFileService(mockFileRepo, mockComponentRepo)

	mockFileRepo.EXPECT().ReadLocalFile("test.js").Return(entities.File{}, errors.New("file not found"))

	file, err := service.GetLocalFile("test.js")

	assert.Error(t, err)
	assert.Equal(t, entities.FileDTO{}, file)
	mockFileRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}

func TestGetRemoteFileContent(t *testing.T) {
	mockFileRepo := mocks.NewMockFileRepository(t)
	mockComponentRepo := mocks.NewMockComponentRepository(t)
	service := service.NewFileService(mockFileRepo, mockComponentRepo)

	expectedFile := entities.NewFile(
		"",
		"remote.js",
		[]byte("function main() {\n\tconsole.log('Hello, World!');\n}"),
	)

	mockComponentRepo.EXPECT().FindByFilePath("remote.js").Return(entities.Component{FileHash: "test-md5"}, nil)
	mockFileRepo.EXPECT().ReadRemoteFileByMD5("remote.js", "test-md5").Return(*expectedFile, nil)

	file, err := service.GetRemoteFile("remote.js")

	assert.NoError(t, err)
	assert.Equal(t, entities.FileDTO{
		Path:     "remote.js",
		Name:     "remote.js",
		Content:  "function main() {\n\tconsole.log('Hello, World!');\n}",
		Language: "javascript",
	}, file)
	mockFileRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}

func TestGetRemoteFileContent_Error(t *testing.T) {
	mockFileRepo := mocks.NewMockFileRepository(t)
	mockComponentRepo := mocks.NewMockComponentRepository(t)
	service := service.NewFileService(mockFileRepo, mockComponentRepo)

	mockComponentRepo.EXPECT().FindByFilePath("remote.js").Return(entities.Component{}, errors.New("component not found"))

	file, err := service.GetRemoteFile("remote.js")

	assert.Error(t, err)
	assert.Equal(t, entities.FileDTO{}, file)
	mockFileRepo.AssertExpectations(t)
	mockComponentRepo.AssertExpectations(t)
}
