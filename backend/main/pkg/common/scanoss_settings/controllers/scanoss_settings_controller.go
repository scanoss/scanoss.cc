package controllers

type ScanossSettingsController interface {
	Save() error
	HasUnsavedChanges() (bool, error)
}
