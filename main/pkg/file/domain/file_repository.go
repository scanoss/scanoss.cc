package domain

type FileRepository interface {
	ReadRemoteFile(path string) (File, error)

	ReadLocalFile(path string) (File, error)
}
