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

package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/repository/mocks"
	internal_test "github.com/scanoss/scanoss.cc/internal"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestFileRepositoryImpl(t *testing.T) {
	cleanup := internal_test.InitializeTestEnvironment(t)
	defer cleanup()

	t.Run("ReadLocalFile", func(t *testing.T) {
		testFilePath := "test.js"
		testContent := []byte("function main() {\n\tconsole.log('Hello, World!');\n}")
		currentPath := t.TempDir()
		absolutePath := filepath.Join(currentPath, testFilePath)
		err := os.WriteFile(absolutePath, testContent, 0644)
		assert.NoError(t, err)

		config.GetInstance().ScanRoot = currentPath

		repo := NewFileRepositoryImpl()

		file, err := repo.ReadLocalFile(testFilePath)
		assert.NoError(t, err)
		assert.Equal(t, testContent, file.GetContent())
	})

	t.Run("ReadRemoteFileByMD5", func(t *testing.T) {
		testFilePath := "test.js"
		testContent := []byte("function main() {\n\tconsole.log('Hello, World!');\n}")
		md5 := "test-md5"

		repo := mocks.NewMockFileRepository(t)

		repo.EXPECT().ReadRemoteFileByMD5(testFilePath, md5).Return(*entities.NewFile(
			"",
			testFilePath,
			testContent,
		), nil)

		file, err := repo.ReadRemoteFileByMD5(testFilePath, md5)

		assert.NoError(t, err)
		assert.Equal(t, testContent, file.GetContent())

		repo.AssertExpectations(t)
	})
}
