package domain

import (
	"context"
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

	//считать запись
	//проверить есть ли запись с текущей датой
	// если есть запись добавить строку
	// если нет с текущей, добавить строку с датой, добавить запись

	return nil
}
