package service

import (
	"sort"

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
	licenses, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	sort.Slice(licenses, func(i, j int) bool {
		return licenses[i].LicenseId < licenses[j].LicenseId
	})

	return licenses, nil
}
