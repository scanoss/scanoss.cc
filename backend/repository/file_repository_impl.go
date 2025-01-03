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
