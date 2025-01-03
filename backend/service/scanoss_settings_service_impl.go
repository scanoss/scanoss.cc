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

type ScanossSettingsServiceImp struct {
	repository repository.ScanossSettingsRepository
}

func NewScanossSettingsServiceImpl(r repository.ScanossSettingsRepository) *ScanossSettingsServiceImp {
	return &ScanossSettingsServiceImp{
		repository: r,
	}
}

func (s *ScanossSettingsServiceImp) Save() error {
	return s.repository.Save()
}

func (s *ScanossSettingsServiceImp) HasUnsavedChanges() (bool, error) {
	return s.repository.HasUnsavedChanges()
}

func (s *ScanossSettingsServiceImp) GetSettings() *entities.SettingsFile {
	return s.repository.GetSettings()
}
