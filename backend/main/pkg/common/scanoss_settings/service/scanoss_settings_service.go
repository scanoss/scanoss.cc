package service

import "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"

//go:generate mockery --name ScanossSettingsService --with-expecter
type ScanossSettingsService interface {
	Save() error
	HasUnsavedChanges() (bool, error)
	GetSettingsFile() *entities.SettingsFile
}
