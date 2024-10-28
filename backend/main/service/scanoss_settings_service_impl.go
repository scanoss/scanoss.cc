package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/entities"
	"github.com/scanoss/scanoss.lui/backend/main/repository"
)

type ScanossSettingsServiceImp struct {
	repository repository.ScanossSettingsRepository
}

func NewScanossSettingsServiceImpl(r repository.ScanossSettingsRepository) *ScanossSettingsServiceImp {
	return &ScanossSettingsServiceImp{
		repository: r,
	}
}

func (s *ScanossSettingsServiceImp) Save() error {
	return s.repository.Save()
}

func (s *ScanossSettingsServiceImp) HasUnsavedChanges() (bool, error) {
	return s.repository.HasUnsavedChanges()
}

func (s *ScanossSettingsServiceImp) GetSettings() *entities.SettingsFile {
	return s.repository.GetSettings()
}
