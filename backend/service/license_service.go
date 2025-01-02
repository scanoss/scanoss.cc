package service

import "github.com/scanoss/scanoss.cc/backend/entities"

type LicenseService interface {
	GetAll() ([]entities.License, error)
}
