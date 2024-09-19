package service

type ScanossBomService interface {
	Save() error
	HasUnsavedChanges() (bool, error)
}
