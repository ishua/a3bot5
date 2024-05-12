package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ishua/a3bot5/libs/closer"
	"github.com/ishua/a3bot5/tbot/internal/botcmd"
	"github.com/ishua/a3bot5/tbot/internal/schema"
	"github.com/redis/go-redis/v9"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

type pubsub struct {
	*redis.Client
}

func (r pubsub) Pub(ctx context.Context, channel, value string) error {
	return r.Publish(ctx, channel, value).Err()
}

var cfg MyConfig
var (
	shutdownTimeout = 5 * time.Second
	// telegram
	Bot     *tgbotapi.BotAPI
	Updates <-chan tgbotapi.Update
	rdb     pubsub
	botCmd  botcmd.CmdRouter
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
	if len(cfg.Users) < 1 {
		panic("no cfg.Users in config")
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
		log.Fatalf("[ERROR] Telegram: %s\n", err)
	}
}

// init redis
func init() {
	log.Println("init redis: " + cfg.Redis)
	rdb = pubsub{redis.NewClient(&redis.Options{
		Addr:     cfg.Redis,
		Password: "", // no password set
		DB:       0,  // use default DB
	})}
}

// init botCmd
func init() {
	botCmd = *botcmd.NewCmdRouter(rdb)
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

			var allowCommands []string
			for _, user := range cfg.Users {
				if user.User == update.Message.Chat.UserName {
					allowCommands = user.Commands
				}
			}

			if len(allowCommands) < 1 {
				errStr := "user haven't allow commands: " + update.Message.Chat.UserName
				botSendErr(errStr, update.Message.Chat.ID, update.Message.MessageID)
				continue
			}

			var fileUrl string
			var err error
			if update.Message.Document != nil {
				fileUrl, err = Bot.GetFileDirectURL(update.Message.Document.FileID)
				if err != nil {
					errStr := "can't get file url " + err.Error()
					botSendErr(errStr, update.Message.Chat.ID, update.Message.MessageID)
					continue
				}
			}
			var replyMsgId int
			if update.Message.ReplyToMessage != nil {
				replyMsgId = update.Message.ReplyToMessage.MessageID
			}
			bmsg := botcmd.Message{
				UserName:         update.Message.Chat.UserName,
				MsgId:            update.Message.MessageID,
				ReplyToMessageID: replyMsgId,
				ChatId:           update.Message.Chat.ID,
				Text:             update.Message.Text,
				Caption:          update.Message.Caption,
				ReplyText:        "",
				FileUrl:          fileUrl,
			}
			err = botCmd.Send(context.Background(), bmsg, allowCommands, cfg.SubChannel)

			if err != nil {
				errStr := "wrong route message " + err.Error()
				botSendErr(errStr, update.Message.Chat.ID, update.Message.MessageID)
			}
		}
	}()

	pubsub := rdb.Subscribe(ctx, cfg.SubChannel)
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var m schema.ChannelMsg
			err := m.UnmarshalBinary([]byte(msg.Payload))
			if err != nil {
				log.Println("wrong unmarshal msg " + err.Error())
				continue
			}
			msg := tgbotapi.NewMessage(m.ChatId, m.ReplyText)
			msg.ParseMode = "html"
			newMsg, err := Bot.Send(msg)
			if m.MsgId != 0 {
				msg.ReplyToMessageID = m.MsgId
			}
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

func botSendErr(errStr string, chatId int64, replyId int) {
	log.Println(errStr)
	msg := tgbotapi.NewMessage(chatId, errStr)
	msg.ParseMode = "html"
	newMsg, err := Bot.Send(msg)
	msg.ReplyToMessageID = replyId
	if err != nil {
		log.Printf("%d %s", newMsg.MessageID, err.Error())
	}
}
