package domain

type FileRepository interface {
	ReadRemoteFileByMD5(path string, md5 string) (File, error)

	ReadLocalFile(path string) (File, error)
}
