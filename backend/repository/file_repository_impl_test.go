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
