package repo

import "context"

type GitFile struct {
	path string
	git  Gitter
}

type Gitter interface {
	Pull(ctx context.Context) error
	CommitAndPush(ctx context.Context, path string)
}

func NewGitFile(path string, git Gitter) *GitFile {
	return &GitFile{
		path: path,
		git:  git,
	}
}
