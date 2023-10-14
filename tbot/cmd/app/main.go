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

	"github.com/cristalhq/aconfig"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MyConfig struct {
	Token string `required:"true" env:"TELEGRAMBOTTOKEN" usage:"token for your telegram bot"`
	Debug bool   `default:"false" usage:"turn on debug mode"`
}

var cfg MyConfig
var (
	shutdownTimeout = 5 * time.Second
	// telegram
	Bot     *tgbotapi.BotAPI
	Updates <-chan tgbotapi.Update
)

// init config
func init() {
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		Files: []string{"config.json"},
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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "msgText")
			msg.ParseMode = "html"
			newMsg, err := Bot.Send(msg)
			if err != nil {
				log.Printf("%d %s", newMsg.MessageID, err.Error())
			}
		}
	}()

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
