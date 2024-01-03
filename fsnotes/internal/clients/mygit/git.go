package mygit

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Repo struct {
	*git.Repository
	Token    string
	Username string
}

func NewClient(path string, url string, token string) (*Repo, error) {
	r, err := git.PlainOpen(path)
	username := "fsnotesReader"
	if err != nil {
		r, err = git.PlainClone(path, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: username, // yes, this can be anything except an empty string
				Password: token,
			},
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			return nil, err
		}
	}

	return &Repo{r, token, username}, nil
}

func (r *Repo) Pull(ctx context.Context) error {
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("pull create worktree: %w", err)
	}
	return w.Pull(&git.PullOptions{
		Auth: &http.BasicAuth{
			Username: r.Username, // yes, this can be anything except an empty string
			Password: r.Token,
		},
		RemoteName: "origin"})
}

func (r *Repo) CommitAndPush(ctx context.Context, path string) error {
	return nil
}
