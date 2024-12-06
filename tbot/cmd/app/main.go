package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/ishua/a3bot5/tbot/internal/botcmd"
	"github.com/ishua/a3bot5/tbot/internal/client/myredis"
	"github.com/ishua/a3bot5/tbot/internal/client/mytgclient"
)

type MyConfig struct {
	Token      string `required:"true" env:"TELEGRAMBOTTOKEN" usage:"token for your telegram bot"`
	Redis      string `default:"redis:6379" env:"REDIS" usage:"connect str to redis"`
	Debug      bool   `default:"false" usage:"turn on debug mode"`
	SubChannel string `default:"tbot" usage:"channel for subscribe jobs"`
	Users      []struct {
		User     string
		Commands []string
	} `usage:"allow users to use bot if empty then allows everybody"`
}

var cfg MyConfig

// init config
func init() {
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files: []string{"conf/tbot_config.yaml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}
	if len(cfg.Users) < 1 {
		panic("no cfg.Users in config")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// init glue
	userSettings := make([]botcmd.UserSettings, len(cfg.Users))
	for idx, cuser := range cfg.Users {
		userSettings[idx] = botcmd.UserSettings{
			Name:     cuser.User,
			Commands: cuser.Commands,
		}
	}
	botcmd := botcmd.NewCmdRouter(userSettings, cfg.SubChannel)

	//init queue client
	q := myredis.NewRedisClient(cfg.Redis, cfg.SubChannel, botcmd)
	botcmd.RegQueue(q)
	//init telegram client
	t, err := mytgclient.NewTgClient(cfg.Token, cfg.Debug, botcmd)
	if err != nil {
		log.Fatal("tg can't connect: " + err.Error())
	}
	botcmd.RegTelegram(t)

	//run listners
	q.ListeningQueue(ctx)
	t.ListeningTg(ctx)

	// stop service here
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// waiting signal for stop
	sig := <-sigChan
	log.Printf("Received signal: %s. Stopping...\n", sig)
	cancel()
	time.Sleep(1 * time.Second)
	log.Println("Program has stopped.")
}
