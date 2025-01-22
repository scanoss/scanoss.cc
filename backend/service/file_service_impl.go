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

package service

import (
	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/repository"
)

type FileServiceImpl struct {
	repo          repository.FileRepository
	componentRepo repository.ComponentRepository
}

func NewFileService(repo repository.FileRepository, componentRepo repository.ComponentRepository) FileService {
	return &FileServiceImpl{
		repo:          repo,
		componentRepo: componentRepo,
	}
}

func (c *FileServiceImpl) GetRemoteFile(path string) (entities.FileDTO, error) {
	component, err := c.componentRepo.FindByFilePath(path)
	if err != nil {
		return entities.FileDTO{}, err
	}

	file, err := c.repo.ReadRemoteFileByMD5(path, component.FileHash)
	return entities.FileDTO{
		Name:     file.GetName(),
		Path:     file.GetRelativePath(),
		Content:  string(file.GetContent()),
		Language: file.GetLanguage(),
	}, err
}

func (c *FileServiceImpl) GetLocalFile(path string) (entities.FileDTO, error) {
	file, err := c.repo.ReadLocalFile(path)
	if err != nil {
		return entities.FileDTO{}, err
	}

	return entities.FileDTO{
		Name:     file.GetName(),
		Path:     file.GetRelativePath(),
		Content:  string(file.GetContent()),
		Language: file.GetLanguage(),
	}, err
}
