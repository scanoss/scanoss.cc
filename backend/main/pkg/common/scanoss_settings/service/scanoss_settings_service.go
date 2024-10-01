package service

//go:generate mockery --name ScanossSettingsService --with-expecter
type ScanossSettingsService interface {
	Save() error
	HasUnsavedChanges() (bool, error)
}
