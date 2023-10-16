package main

import (
	"context"
	"fmt"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/redis/go-redis/v9"
)

type MyConfig struct {
	Redis       string `default:"redis:6379" env:"REDIS" usage:"connect str to redis"`
	SubChannel  string `default:"restjobs" usage:"channel for subscribe jobs"`
	TbotChannel string `default:"tbot" usage:"channel for jobs for tbot"`
}

var (
	cfg MyConfig
	rdb *redis.Client
)

// init config
func init() {
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files: []string{"conf/restjobs_config.yaml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}
}

// init redis
func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}

func main() {
	ctx := context.Background()
	// There is no error because go-redis automatically reconnects on error.
	pubsub := rdb.Subscribe(ctx, cfg.SubChannel)
	// Close the subscription when we are done.
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(msg.Channel, msg.Payload)
	}

}
