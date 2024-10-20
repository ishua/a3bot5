package mcore

import (
	"context"

	"github.com/ishua/a3bot5/mcore/pkg/schema"
)

type ClientMcore struct {
	addr string
}

func (c *ClientMcore) AddMsg(ctx context.Context, msg schema.TelegramMsg) error {
	return nil
}

func (с *ClientMcore) GetMsg(ctx context.Context, queueName string) (schema.TelegramMsg, error) {
	return schema.TelegramMsg{}, nil
}
