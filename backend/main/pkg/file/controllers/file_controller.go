package controllers

import "github.com/scanoss/scanoss.lui/backend/main/pkg/file/entities"

type FileController interface {
	GetRemoteFile(path string) (entities.FileDTO, error)
	GetLocalFile(path string) (entities.FileDTO, error)
}
