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
