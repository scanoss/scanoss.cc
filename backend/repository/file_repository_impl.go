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
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/config"
	"github.com/scanoss/scanoss.cc/internal/fetch"
)

type FileRepositoryImpl struct{}

func NewFileRepositoryImpl() FileRepository {
	return &FileRepositoryImpl{}
}

func (r *FileRepositoryImpl) ReadLocalFile(path string) (entities.File, error) {
	scanRootPath := config.GetInstance().ScanRoot

	absolutePath := filepath.Join(scanRootPath, path)

	content, err := os.ReadFile(absolutePath)
	if err != nil {
		return entities.File{}, fmt.Errorf("%s does not exist. Try changing the scan root from the status bar", path)
	}

	return *entities.NewFile(scanRootPath, path, content), nil
}

func (r *FileRepositoryImpl) ReadRemoteFileByMD5(path string, md5 string) (entities.File, error) {
	baseURL := config.GetInstance().ApiUrl
	token := config.GetInstance().ApiToken

	url := fmt.Sprintf("%s/file_contents/%s", baseURL, md5)

	headers := make(map[string]string)
	if token != "" {
		headers["X-Session"] = token
	}

	options := fetch.Options{
		Method:  http.MethodGet,
		Headers: headers,
	}

	body, err := fetch.Text(url, options)
	if err != nil {
		return entities.File{}, fmt.Errorf("failed to fetch file content: %w", err)
	}

	basePath, err := os.Getwd()
	if err != nil {
		return entities.File{}, fmt.Errorf("failed to get current working directory: %w", err)
	}

	return *entities.NewFile(basePath, path, []byte(body)), nil
}

func (r *FileRepositoryImpl) GetComponentByFilePath(filePath string) (entities.Component, error) {
	return entities.Component{}, nil
}
