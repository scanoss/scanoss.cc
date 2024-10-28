package utils

type FileReader interface {
	ReadFile(filePath string) ([]byte, error)
}
