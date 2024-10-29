package service

import "github.com/scanoss/scanoss.lui/backend/main/entities"

type LicenseService interface {
	GetAll() ([]entities.License, error)
}
