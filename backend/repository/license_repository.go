package repository

import "github.com/scanoss/scanoss.cc/backend/entities"

type LicenseRepository interface {
	GetAll() ([]entities.License, error)
}
