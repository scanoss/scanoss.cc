package service

import "github.com/scanoss/scanoss.lui/backend/main/entities"

type FileService interface {
	GetRemoteFile(path string) (entities.FileDTO, error)
	GetLocalFile(path string) (entities.FileDTO, error)
}
