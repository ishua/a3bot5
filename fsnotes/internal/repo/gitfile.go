package repo

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type GitFile struct {
	path string
	git  Gitter
}

type Gitter interface {
	Pull(ctx context.Context) error
	CommitAndPush(ctx context.Context, path []string) error
}

func NewGitFile(path string, git Gitter) *GitFile {
	return &GitFile{
		path: path,
		git:  git,
	}
}

func (g *GitFile) ReadRows(ctx context.Context) ([]string, error) {
	err := g.git.Pull(ctx)
	if err != nil {
		return nil, fmt.Errorf("gitfile ReadRows pull: %w", err)
	}

	f, err := os.Open(g.path)
	if err != nil {
		return nil, fmt.Errorf("gitfile readfile: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var ret []string
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("gitfile scaner: %w", err)
	}

	return ret, nil
}

func (g *GitFile) AddRows(ctx context.Context, str []string) error {
	err := g.git.Pull(ctx)
	if err != nil {
		return fmt.Errorf("gitfile AddRows pull: %w", err)
	}

	f, err := os.OpenFile(g.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("gitfile AddRows OpenFile: %w", err)
	}

	for _, s := range str {
		_, err = f.WriteString(s + "\n")
		if err != nil {
			return fmt.Errorf("gitfile AddRows WriteString: %w", err)
		}
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("gitfile AddRows file.Close: %w", err)
	}

	err = g.git.CommitAndPush(ctx, []string{g.path})
	if err != nil {
		return fmt.Errorf("gitfile AddRows CommitAndPush: %w", err)
	}
	return nil
}
