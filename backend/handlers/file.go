package handlers

import (
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/entities"
)

type FileHandler struct {
	fileModule *file.Module
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		fileModule: file.NewModule(),
	}
}

func (fh *FileHandler) FileGetLocalContent(path string) entities.FileDTO {
	f, _ := fh.fileModule.Controller.GetLocalFile(path)
	return f
}

func (fh *FileHandler) FileGetRemoteContent(path string) entities.FileDTO {
	f, _ := fh.fileModule.Controller.GetRemoteFile(path)
	return f
}
