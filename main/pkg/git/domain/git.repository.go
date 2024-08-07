package domain

type GitRepository interface {
	GetFiles() ([]File, error)
}
