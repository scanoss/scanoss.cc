package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/entities"
	"github.com/scanoss/scanoss.lui/backend/main/repository"
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
	return s.repo.GetAll()
}
