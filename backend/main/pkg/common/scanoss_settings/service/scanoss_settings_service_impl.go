package service

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/repository"
)

type ScanossSettingsServiceImp struct {
	repository repository.ScanossSettingsRepository
}

func NewScanossSettingsService(r repository.ScanossSettingsRepository) *ScanossSettingsServiceImp {
	return &ScanossSettingsServiceImp{
		repository: r,
	}
}

func (us *ScanossSettingsServiceImp) Save() error {
	return us.repository.Save()
}

func (us *ScanossSettingsServiceImp) HasUnsavedChanges() (bool, error) {
	return us.repository.HasUnsavedChanges()
}
