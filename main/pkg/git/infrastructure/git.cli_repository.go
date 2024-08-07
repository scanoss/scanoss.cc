package infrastructure

import (
	"errors"
	"integration-git/main/pkg/git/domain"
	"os/exec"
	"path"
	"strings"
)

var (
	ErrExecutingGitStatus    = errors.New("error while executing git status command")
	ErrInvalidGitStatusInput = errors.New("invalid git status input")
)

type GitCliRepository struct{}

func NewGitRepository() *GitCliRepository {
	return &GitCliRepository{}
}

func (r *GitCliRepository) GetFiles() ([]domain.File, error) {
	cmd := exec.Command("git", "status", "-s")
	output, err := cmd.Output()
	if err != nil {
		return nil, ErrExecutingGitStatus
	}

	return r.GetFileStatusFromOutputGitCli(string(output))
}

func (r *GitCliRepository) GetFileStatusFromOutputGitCli(s string) ([]domain.File, error) {

	lines := strings.Split(s, "\n")
	var files []domain.File
	for _, line := range lines {

		if len(line) == 0 {
			continue
		}

		if len(line) < 4 {
			return nil, ErrInvalidGitStatusInput //<X><Y> <PATH>
		}

		status := line[0:2]
		filepath := line[3:] //TODO: Handle case when there is a rename src/main.c -> src/main_new_name.c
		filename := path.Base(filepath)

		// Get the state and path name
		stagingStatus, err := domain.NewFileStatus(string(status[0]))
		if err != nil {
			return nil, err
		}

		workingStatus, err := domain.NewFileStatus(string(status[1]))
		if err != nil {
			return nil, err
		}

		f, err := domain.NewFile(filename, filepath, stagingStatus, workingStatus)
		if err != nil {
			return nil, err
		}
		files = append(files, *f)
	}

	return files, nil
}
