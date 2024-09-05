package handlers

import (
	"context"

	"github.com/scanoss/scanoss.lui/backend/main/pkg/file"
	"github.com/scanoss/scanoss.lui/backend/main/pkg/file/entities"
)

type FileHandler struct {
	ctx        context.Context
	fileModule *file.Module
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		fileModule: file.NewModule(),
	}
}

// Get local file content
func (fh *FileHandler) FileGetLocalContent(path string) entities.FileDTO {
	f, _ := fh.fileModule.Controller.GetLocalFile(path)
	return f
}

// Get remote file content
func (fh *FileHandler) FileGetRemoteContent(path string) entities.FileDTO {
	f, _ := fh.fileModule.Controller.GetRemoteFile(path)
	return f
}
