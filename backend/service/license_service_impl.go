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
	"sort"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/backend/repository"
)

type LicenseServiceImpl struct {
	repo repository.LicenseRepository
}

func NewLicenseServiceImpl(repo repository.LicenseRepository) LicenseService {
	return &LicenseServiceImpl{
		repo: repo,
	}
}

func (s *LicenseServiceImpl) GetAll() ([]entities.License, error) {
	licenses, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	sort.Slice(licenses, func(i, j int) bool {
		return licenses[i].LicenseId < licenses[j].LicenseId
	})

	return licenses, nil
}
