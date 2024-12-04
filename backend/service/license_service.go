package service

import "github.com/scanoss/scanoss.lui/backend/entities"

type LicenseService interface {
	GetAll() ([]entities.License, error)
}
