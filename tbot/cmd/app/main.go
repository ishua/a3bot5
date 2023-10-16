package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ishua/a3bot5/libs/closer"
	"github.com/redis/go-redis/v9"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MyConfig struct {
	Token           string `required:"true" env:"TELEGRAMBOTTOKEN" usage:"token for your telegram bot"`
	Redis           string `default:"redis:6379" env:"REDIS" usage:"connect str to redis"`
	Debug           bool   `default:"false" usage:"turn on debug mode"`
	SubChannel      string `default:"tbot" usage:"channel for subscribe jobs"`
	RestJobsChannel string `default:"restjobs" usage:"channel for jobs for restjobs"`
}

var cfg MyConfig
var (
	shutdownTimeout = 5 * time.Second
	// telegram
	Bot     *tgbotapi.BotAPI
	Updates <-chan tgbotapi.Update
	rdb     *redis.Client
)

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
}

// init telegram
func init() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(cfg.Token)
	Bot.Debug = cfg.Debug

	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Telegram: %s\n", err)
		os.Exit(1)
	}
	log.Printf("[INFO] Authorized: %s", Bot.Self.UserName)
	// delete webhook
	dwh := tgbotapi.DeleteWebhookConfig(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true})
	_, err = Bot.Request(dwh)

	if err != nil {
		log.Fatal(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	Updates = Bot.GetUpdatesChan(u)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Telegram: %s\n", err)
		os.Exit(1)
	}
}

// init redis
func init() {
	log.Println("init redis: " + cfg.Redis)
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}

}

func run(ctx context.Context) error {
	closer := new(closer.Closer)

	go func() {
		for update := range Updates {
			// ignore edited messages
			if update.Message == nil {
				continue
			}

			q := qmsg{
				Command:          "/rate_usd", //TODO
				UserName:         update.Message.Chat.UserName,
				MsgId:            update.Message.MessageID,
				ReplyToMessageID: 0,
				ChatId:           update.Message.Chat.ID,
				Text:             update.Message.Text,
			}

			payload, err := q.marshalBinary()
			if err != nil {
				log.Println("marshal err:" + err.Error())
				continue
			}

			err = rdb.Publish(ctx, cfg.RestJobsChannel, //TODO
				string(payload)).Err()

			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				msg.ParseMode = "html"
				newMsg, err := Bot.Send(msg)
				if err != nil {
					log.Printf("%d %s", newMsg.MessageID, err.Error())
				}
			}

		}
	}()

	pubsub := rdb.Subscribe(ctx, cfg.SubChannel)
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var m qmsg
			err := m.unmarshalBinary([]byte(msg.Payload))
			if err != nil {
				log.Println("wrong unmarshal msg " + err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(m.ChatId, m.ReplyText)
			msg.ParseMode = "html"
			newMsg, err := Bot.Send(msg)
			if err != nil {
				log.Printf("%d %s", newMsg.MessageID, err.Error())
			}

		}
	}()

	closer.Add(func(ctx context.Context) error {
		pubsub.Close()
		return nil
	})

	//waiting
	<-ctx.Done() // block

	log.Println("shutting down server gracefully")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := closer.Close(shutdownCtx); err != nil {
		return fmt.Errorf("closer: %v", err)
	}

	return nil
}

type qmsg struct {
	Command          string `json:"command"`
	UserName         string `json:"userName"`
	MsgId            int    `json:"msgId"`
	ReplyToMessageID int    `json:"replyToMessageID"`
	ChatId           int64  `json:"chatId"`
	Text             string `json:"text"`
	ReplyText        string `json:"replyText"`
}

func (m *qmsg) marshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *qmsg) unmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	return nil
}
