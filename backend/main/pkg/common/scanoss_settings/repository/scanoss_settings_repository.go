package repository

import "github.com/scanoss/scanoss.lui/backend/main/pkg/common/scanoss_settings/entities"

type ScanossSettingsRepository interface {
	Save() error
	Read() (entities.SettingsFile, error)
	HasUnsavedChanges() (bool, error)
}
