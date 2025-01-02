package service

import "github.com/scanoss/scanoss.cc/backend/entities"

type FileService interface {
	GetRemoteFile(path string) (entities.FileDTO, error)
	GetLocalFile(path string) (entities.FileDTO, error)
}
