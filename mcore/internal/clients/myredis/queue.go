package myredis

import (
	"context"

	"github.com/ishua/a3bot5/mcore/pkg/schema"
	"github.com/redis/go-redis/v9"
)

type MessageQueue struct {
	*redis.Client
}

func NewMessageQueue(redisAddr string) *MessageQueue {
	return &MessageQueue{redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})}
}

func (m *MessageQueue) AddTelegramMsg(ctx context.Context, msg schema.TelegramMsg) error {

	value, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = m.LPush(ctx, msg.QueueName, value).Result()
	return err
}

func (m *MessageQueue) GetTelegramMsg(ctx context.Context, queueName string) (schema.TelegramMsg, error) {
	value, err := m.RPop(ctx, queueName).Result()
	if err != nil {
		return schema.TelegramMsg{}, err
	}

	ret := schema.TelegramMsg{}
	err = ret.UnmarshalBinary([]byte(value))
	return ret, err
}
