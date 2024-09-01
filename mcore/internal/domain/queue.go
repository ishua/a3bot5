package domain

import "context"

type AddQueue interface {
	Add(ctx context.Context, queueName, payload string) error
}

type GetLastQueue interface {
	GetLast(ctx context.Context, queueName string) ([]byte, error)
}
