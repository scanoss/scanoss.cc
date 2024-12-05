package repository

import (
	_ "embed"

	"github.com/scanoss/scanoss.lui/backend/entities"
	"github.com/scanoss/scanoss.lui/internal/utils"
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
