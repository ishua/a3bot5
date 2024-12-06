package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/ishua/a3bot5/mcore/internal/clients/myredis"
	"github.com/ishua/a3bot5/mcore/internal/domain"
	"github.com/ishua/a3bot5/mcore/internal/mcsrv"
)

type MyConfig struct {
	Redis      string   `default:"redis:6379" env:"REDIS" usage:"connect str to redis"`
	ListenPort string   `default:":8080" usage:"port where start http rest"`
	Debug      bool     `default:"false" usage:"turn on debug mode"`
	Secrets    []string `default:"mysecret,mysecret2" usage:"secrets for http connect header 'secret'"`
	RootPath   string   `default:"/api" usage:"path begin from this string"`
}

var (
	cfg       MyConfig
	addMsgUrl string
	getMsgUrl string
	pingUrl   string
)

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

	addMsgUrl = fmt.Sprintf("%s/add-msg/", cfg.RootPath)
	getMsgUrl = fmt.Sprintf("%s/get-msg/", cfg.RootPath)
	pingUrl = fmt.Sprintf("%s/ping/", cfg.RootPath)

	mux.HandleFunc("POST "+addMsgUrl, server.AddMsg)
	mux.HandleFunc("POST "+getMsgUrl, server.GetMsg)
	mux.HandleFunc("GET "+pingUrl, server.Ping)

	log.Println("start server port" + cfg.ListenPort)
	err := http.ListenAndServe(cfg.ListenPort, myMiddle(mux, cfg.Secrets))
	if err != nil {
		panic(err)
	}

}

func myMiddle(next http.Handler, secrets []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		if url.Path == pingUrl {
			next.ServeHTTP(w, r)
			return
		}
		secret := r.Header.Get("secret")
		for _, s := range secrets {
			if s == secret {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
