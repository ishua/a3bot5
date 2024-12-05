package mytgclient

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgClient struct {
	bot *tgbotapi.BotAPI
}

func NewTgClient(token string, debug bool) (*TgClient, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = debug
	log.Printf("[INFO] Authorized: %s", bot.Self.UserName)

	// delete webhook
	dwh := tgbotapi.DeleteWebhookConfig(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true})
	_, err = bot.Request(dwh)

	if err != nil {
		return nil, err
	}

	return &TgClient{
		bot: bot,
	}, nil
}
