package domain

import (
	"errors"
	"os"
)

var (
	ErrPathDoesNotExist = errors.New("path does not exist")
	ErrPathIsDirectory  = errors.New("path is a directory")
	ErrOpeningFile      = errors.New("opening file failed")
)

type stagedFile struct {
	path    string
	name    string
	content []byte
	diff    string
	status  string //TODO Change for enum Added, Modified, etc
}

func NewStagedFile(path string) (*stagedFile, error) {

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return &stagedFile{}, ErrPathDoesNotExist
	}

	if fileInfo.IsDir() {
		return &stagedFile{}, ErrPathIsDirectory
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return &stagedFile{}, ErrOpeningFile
	}

	return &stagedFile{
		path:    path,
		name:    fileInfo.Name(),
		content: file,
		status:  "",
	}, nil

}
