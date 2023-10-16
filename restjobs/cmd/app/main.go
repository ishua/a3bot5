package main

import (
	"context"
	"log"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/ishua/a3bot5/restjobs/internal/clients/cbrapi"
	"github.com/ishua/a3bot5/restjobs/internal/domain"
	"github.com/redis/go-redis/v9"
)

type MyConfig struct {
	Redis       string `default:"redis:6379" env:"REDIS" usage:"connect str to redis"`
	SubChannel  string `default:"restjobs" usage:"channel for subscribe jobs"`
	TbotChannel string `default:"tbot" usage:"channel for jobs for tbot"`
}

type Pubsub struct {
	*redis.Client
}

var (
	cfg MyConfig
	rdb *Pubsub
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
	rdb = &Pubsub{redis.NewClient(&redis.Options{
		Addr:     cfg.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})}
}
func (r Pubsub) Pub(ctx context.Context, value string) {
	err := r.Publish(ctx, cfg.TbotChannel, value).Err()
	if err != nil {
		log.Println("publish error: " + err.Error())
	}
}

func main() {
	m := domain.NewModel(&cbrapi.CbrClient{}, rdb)

	ctx := context.Background()
	// There is no error because go-redis automatically reconnects on error.
	pubsub := rdb.Subscribe(ctx, cfg.SubChannel)
	// Close the subscription when we are done.
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		go m.DoJob(ctx, msg.Payload)
	}
}
