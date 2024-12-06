package myredis

import (
	"context"
	"fmt"
	"log"

	"github.com/ishua/a3bot5/mcore/pkg/schema"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
	channel    string
	telegramer Telegramer
}

type Telegramer interface {
	Send2Telegram(ctx context.Context, msg schema.TelegramMsg)
}

func NewRedisClient(adr, channel string, t Telegramer) *RedisClient {
	return &RedisClient{
		redis.NewClient(&redis.Options{
			Addr:     adr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		channel,
		t,
	}
}

func (c *RedisClient) AddMsg(ctx context.Context, msg schema.TelegramMsg) error {
	payload, err := msg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal err: %w", err)
	}
	return c.Publish(ctx, msg.QueueName, payload).Err()
}

func (c *RedisClient) ListeningQueue(ctx context.Context) {
	pubsub := c.Subscribe(ctx, c.channel)
	go func() {
		log.Println("start listen reddis channel: " + c.channel)
		ch := pubsub.Channel()

		for {
			select {
			case msg := <-ch:
				var m schema.TelegramMsg
				err := m.UnmarshalBinary([]byte(msg.Payload))
				if err != nil {
					log.Println("wrong unmarshal msg " + err.Error())
					continue
				}
				c.telegramer.Send2Telegram(ctx, m)
			case <-ctx.Done():
				log.Println("stopping listen redis")
				pubsub.Close()
				return
			}
		}
	}()
}
