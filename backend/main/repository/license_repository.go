package repository

import "github.com/scanoss/scanoss.lui/backend/main/entities"

type LicenseRepository interface {
	GetAll() ([]entities.License, error)
}
