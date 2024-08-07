package domain

type FileRepository interface {
	ReadLocalFile(path string) (File, error)
	ReadRemoteFile(path string) (File, error)
}
