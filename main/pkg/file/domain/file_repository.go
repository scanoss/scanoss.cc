package domain

type FileRepository interface {
	ReadLocalFile(path string) ([]byte, error)
	ReadRemoteFile(path string) ([]byte, error)
}
