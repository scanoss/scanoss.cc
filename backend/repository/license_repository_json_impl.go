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
	_ "embed"

	"github.com/scanoss/scanoss.cc/backend/entities"
	"github.com/scanoss/scanoss.cc/internal/utils"
)

const LICENSES_FILE_PATH = "backend/data/spdx-licenses.json"

//go:embed data/spdx-licenses.json
var licensesFile []byte

type LicenseRepositoryJsonImpl struct {
	fr utils.FileReader
}

func NewLicenseJsonRepository(fr utils.FileReader) LicenseRepository {
	return &LicenseRepositoryJsonImpl{
		fr: fr,
	}
}

func (r *LicenseRepositoryJsonImpl) GetAll() ([]entities.License, error) {
	jsonData, err := utils.JSONParse[entities.LicenseFile](licensesFile)
	if err != nil {
		return []entities.License{}, err
	}

	return jsonData.Licenses, nil
}
