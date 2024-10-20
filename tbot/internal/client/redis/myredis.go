package myredis

import (
	"context"

	"github.com/ishua/a3bot5/mcore/pkg/schema"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func (c *RedisClient) AddMsg(ctx context.Context, msg schema.TelegramMsg) error {
	return nil
}

func (с *RedisClient) GetMsg(ctx context.Context, queueName string) (schema.TelegramMsg, error) {
	return schema.TelegramMsg{}, nil
}
