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
