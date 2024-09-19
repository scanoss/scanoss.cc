package controllers

type ScanossBomController interface {
	Save() error
	HasUnsavedChanges() (bool, error)
}
