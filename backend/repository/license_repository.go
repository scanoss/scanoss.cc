package repository

import "github.com/scanoss/scanoss.lui/backend/entities"

type LicenseRepository interface {
	GetAll() ([]entities.License, error)
}
