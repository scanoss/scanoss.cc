package repository

import "github.com/scanoss/scanoss.cc/backend/entities"

type FileRepository interface {
	ReadLocalFile(filePath string) (entities.File, error)
	ReadRemoteFileByMD5(path string, md5 string) (entities.File, error)
	GetComponentByFilePath(filePath string) (entities.Component, error)
}
