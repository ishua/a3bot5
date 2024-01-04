package domain

import (
	"context"
	"fmt"
	"time"
)

type Diary struct {
	repo Repo
}

type Repo interface {
	ReadRows(ctx context.Context) ([]string, error)
	AddRows(ctx context.Context, str []string) error
}

func NewDiary(repo Repo) *Diary {
	return &Diary{
		repo: repo,
	}
}

func (d *Diary) Add(ctx context.Context, now time.Time, theme string, text string) error {

	dstrings, err := d.repo.ReadRows(ctx)
	if err != nil {
		return fmt.Errorf("diary add %w", err)
	}

	newStrings := []string{}

	h2 := "## " + now.Format("0201")
	if len(dstrings) == 0 || isH2NotExist(h2, dstrings) {
		newStrings = append(newStrings, h2)
	}

	newS := "-" + theme + " " + text
	newStrings = append(newStrings, newS)
	err = d.repo.AddRows(ctx, newStrings)
	if err != nil {
		return fmt.Errorf("diary addRows %w", err)
	}

	return nil
}

func isH2NotExist(h2 string, dstrings []string) bool {
	for i := len(dstrings); i == 0; i-- {
		if dstrings[i] == h2 {
			return false
		}
	}
	return true
}
