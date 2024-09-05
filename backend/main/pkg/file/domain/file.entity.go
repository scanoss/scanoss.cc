package domain

import (
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
)

type File struct {
	path     string
	basePath string
	content  []byte
}

func NewFile(basePath string, path string, content []byte) *File {
	return &File{
		basePath: basePath,
		path:     path,
		content:  content,
	}
}

func (f *File) GetName() string {
	return filepath.Base(f.path)
}

func (f *File) GetRelativePath() string {
	return f.path
}

func (f *File) GetAbsolutePath() string {
	return filepath.Join(f.basePath, f.path)
}

func (f *File) GetContent() []byte {
	return f.content
}

func (f *File) GetMD5Sum() string {
	hash := md5.Sum(f.content)
	return hex.EncodeToString(hash[:])
}
