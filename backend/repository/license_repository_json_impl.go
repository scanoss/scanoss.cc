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
