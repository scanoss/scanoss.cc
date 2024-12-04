package service

type ScanossSettingsService interface {
	Save() error
	HasUnsavedChanges() (bool, error)
}
