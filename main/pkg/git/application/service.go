package application

type GitService struct {
	gr GitRepository
}

func NewGitService(r GitRepository) *GitService {
	return &GitService{gr: r}
}
