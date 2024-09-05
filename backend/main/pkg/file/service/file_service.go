package service

import "github.com/scanoss/scanoss.lui/backend/main/pkg/file/entities"

type FileService interface {
	GetLocalFileContent(path string) (entities.File, error)
	GetRemoteFileContent(path string, md5 string) (entities.File, error)
}
