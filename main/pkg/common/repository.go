package common

type FileRepository interface {
	ReadFile(path string) ([]byte, error)
}
