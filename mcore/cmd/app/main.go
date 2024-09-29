package main

import (
	"context"
	"log"
	"net/http"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/ishua/a3bot5/mcore/internal/clients/myredis"
	"github.com/ishua/a3bot5/mcore/internal/domain"
	"github.com/ishua/a3bot5/mcore/internal/mcsrv"
)

type MyConfig struct {
	Redis string `default:"redis:6379" env:"REDIS" usage:"connect str to redis"`
	Debug bool   `default:"false" usage:"turn on debug mode"`
}

var cfg MyConfig

//очередь для хранения заданий
//ручку для того что бы положить задание
//ручку что бы взять задание

func main() {

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files: []string{"conf/mcore_config.yaml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})
	if err := loader.Load(); err != nil {
		log.Fatal(err)
	}

	log.Println("init redis: " + cfg.Redis)
	mq := myredis.NewMessageQueue(cfg.Redis)

	ctx := context.Background()
	if err := mq.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	md := domain.NewMyDomain(mq, mq)

	mux := http.NewServeMux()
	server := mcsrv.NewSrvHandlers(md)

	mux.HandleFunc("POST /add-msg/", server.AddMsg)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
