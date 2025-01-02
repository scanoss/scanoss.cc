package repository

import "github.com/scanoss/scanoss.cc/backend/entities"

type ScanossSettingsRepository interface {
	Init() error
	Save() error
	Read() (entities.SettingsFile, error)
	HasUnsavedChanges() (bool, error)
	AddBomEntry(newEntry entities.ComponentFilter, filterAction string) error
	ClearAllFilters() error
	GetSettings() *entities.SettingsFile
	GetDeclaredPurls() []string
}
