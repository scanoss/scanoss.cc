package application

import (
	"integration-git/main/pkg/git/domain"
)

type GitService struct {
	gr domain.GitRepository
}

func NewGitService(gr domain.GitRepository) *GitService {
	return &GitService{gr}
}

func (gs *GitService) GetFilesToBeCommited() ([]domain.File, error) {
	var filesToBeCommited []domain.File

	files, err := gs.gr.GetFiles()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.ToBeCommitted() {
			filesToBeCommited = append(filesToBeCommited, file)
		}
	}

	return filesToBeCommited, nil
}
