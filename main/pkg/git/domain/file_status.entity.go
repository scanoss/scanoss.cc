package domain

import (
	"errors"
	"strings"
)

var (
	ErrStatusMustContainOnlyOneCharacter = errors.New("status must contain only one character")
	ErrInvalidStatus                     = errors.New("invalid status")
)

var validFileStatus = []string{" ", "M", "T", "A", "D", "R", "C", "U", "?"}

// FileStatus represents the status of a file in Git.
type FileStatus struct {
	Status string //Read about file states https://git-scm.com/docs/git-status#_short_format
}

func NewFileStatus(status string) (FileStatus, error) {

	// Ensure the status is exactly one character
	if len(status) != 1 {
		return FileStatus{}, ErrStatusMustContainOnlyOneCharacter
	}

	// Convert status to uppercase to ensure case-insensitive comparison
	status = strings.ToUpper(status)

	if joinedValidStates := strings.Join(validFileStatus, ""); strings.Contains(joinedValidStates, status) == false {
		return FileStatus{}, ErrInvalidStatus
	}

	return FileStatus{Status: status}, nil
}

// IsUnmodified checks if the file is unmodified.
func (fs FileStatus) IsUnmodified() bool {
	return fs.Status == " "
}

// IsModified checks if the file is modified.
func (fs FileStatus) IsModified() bool {
	return fs.Status == "M"
}

// IsTypeChanged checks if the file type has changed.
func (fs FileStatus) IsTypeChanged() bool {
	return fs.Status == "T"
}

// IsAdded checks if the file is added.
func (fs FileStatus) IsAdded() bool {
	return fs.Status == "A"
}

// IsDeleted checks if the file is deleted.
func (fs FileStatus) IsDeleted() bool {
	return fs.Status == "D"
}

// IsRenamed checks if the file is renamed.
func (fs FileStatus) IsRenamed() bool {
	return fs.Status == "R"
}

// IsCopied checks if the file is copied.
func (fs FileStatus) IsCopied() bool {
	return fs.Status == "C"
}

// IsUpdatedButUnmerged checks if the file is updated but unmerged.
func (fs FileStatus) IsUpdatedButUnmerged() bool {
	return fs.Status == "U"
}
