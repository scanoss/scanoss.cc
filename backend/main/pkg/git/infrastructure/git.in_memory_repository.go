package infrastructure

import (
	"integration-git/main/pkg/git/domain"
)

type GitInMemoryRepository struct {
	files []domain.File
}

func NewGitInMemoryRepository() *GitInMemoryRepository {
	return &GitInMemoryRepository{}
}

func (r *GitInMemoryRepository) GetFiles() ([]domain.File, error) {
	return r.files, nil
}

func (r *GitInMemoryRepository) AddFile(file domain.File) {
	r.files = append(r.files, file)
}

func (r *GitInMemoryRepository) AddBulk(files []domain.File) {
	r.files = append(r.files, files...)
}

func (r *GitInMemoryRepository) RemoveFile(file domain.File) {
	for i, f := range r.files {
		if f == file {
			r.files = append(r.files[:i], r.files[i+1:]...)
		}
	}
}
