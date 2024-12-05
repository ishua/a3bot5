package myredis

import (
	"context"
	"log"

	"github.com/ishua/a3bot5/mcore/pkg/schema"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
	channel string
}

type Telegramer interface {
	send(ctx context.Context, msg schema.TelegramMsg) error
}

func NewRedisClient(adr, channel string) *RedisClient {
	return &RedisClient{
		redis.NewClient(&redis.Options{
			Addr:     adr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		channel,
	}
}

func (c *RedisClient) AddMsg(ctx context.Context, msg schema.TelegramMsg) error {
	return nil
}

func (c *RedisClient) ListeningQueue(ctx context.Context, t Telegramer) {
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
				err = t.send(ctx, m)
				if err != nil {
					log.Printf("redis queue try to send:%s %s", m.String(), err.Error())
				}
			case <-ctx.Done():
				log.Println("stopping listen redis")
				pubsub.Close()
				return
			}
		}
	}()
}
