package domain

import (
	"errors"
)

var (
	ErrFileNameCannotBeEmpty = errors.New("file name can not be empty")
	ErrFilePathCannotBeEmpty = errors.New("file path can not be empty")
)

// TODO: Add NewPath and OldPath. To be used when the state is renamed, otherwise both are equal
// File represents a file in the Git repository
type File struct {
	Name          string
	Path          string
	StagingStatus FileStatus
	WorkingStatus FileStatus
}

// NewFile creates a new instance of File.
func NewFile(name, path string, stagingStatus, workingStatus FileStatus) (*File, error) {
	if name == "" {
		return nil, ErrFileNameCannotBeEmpty
	}
	if path == "" {
		return nil, ErrFilePathCannotBeEmpty
	}

	return &File{
		Name:          name,
		Path:          path,
		WorkingStatus: workingStatus,
		StagingStatus: stagingStatus,
	}, nil
}

func (f *File) GetStagingStatus() FileStatus {
	return f.StagingStatus
}

func (f *File) GetWorkingStatus() FileStatus {
	return f.WorkingStatus
}

func (f *File) ToBeCommitted() bool {
	if f.GetStagingStatus().IsAdded() || f.GetStagingStatus().IsModified() {
		return true
	}
	return false
}
