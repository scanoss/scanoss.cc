package application

import "integration-git/main/pkg/git/domain"

type GetFilesToBeCommitedUseCase struct {
	gs *GitService
}

func NewGetFilesToBeCommitedUseCase(gs *GitService) *GetFilesToBeCommitedUseCase {
	return &GetFilesToBeCommitedUseCase{gs}
}

func (uc *GetFilesToBeCommitedUseCase) Execute() ([]domain.File, error) {
	files, err := uc.gs.GetFilesToBeCommited()

	if err != nil {
		return nil, err
	}

	return files, nil
}
