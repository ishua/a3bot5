package mytgclient

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ishua/a3bot5/mcore/pkg/schema"
)

type TgClient struct {
	bot       *tgbotapi.BotAPI
	botRouter botRouter
}

type botRouter interface {
	Add2Queue(ctx context.Context, msg schema.TelegramMsg) error
}

func NewTgClient(token string, debug bool, botRouter botRouter) (*TgClient, error) {
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
		bot:       bot,
		botRouter: botRouter,
	}, nil
}

func (tg *TgClient) ListeningTg(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tg.bot.GetUpdatesChan(u)

	go func() {
		for {
			select {
			case update := <-updates:
				{
					if update.Message == nil {
						continue
					}
					msg, err := tg.createMessage(update)
					if err != nil {
						tg.SendMsg(err.Error(), update.Message.Chat.ID, getReplyId(update))
						continue
					}

					err = tg.botRouter.Add2Queue(ctx, msg)
					if err != nil {
						tg.SendMsg(err.Error(), update.Message.Chat.ID, getReplyId(update))
					}
				}

			case <-ctx.Done():
				{
					log.Println("stopping listen telegram")
					return
				}

			}
		}
	}()
}

func (tg *TgClient) createMessage(update tgbotapi.Update) (schema.TelegramMsg, error) {

	var fileUrl string
	var err error
	if update.Message.Document != nil {
		fileUrl, err = tg.bot.GetFileDirectURL(update.Message.Document.FileID)
		if err != nil {
			return schema.TelegramMsg{}, fmt.Errorf("tg createMessage can't get file url: %w", err)
		}
	}

	return schema.TelegramMsg{
		UserName:         update.Message.Chat.UserName,
		MsgId:            update.Message.MessageID,
		ReplyToMessageID: getReplyId(update),
		ChatId:           update.Message.Chat.ID,
		Text:             update.Message.Text,
		Caption:          update.Message.Caption,
		ReplyText:        "",
		FileUrl:          fileUrl,
	}, nil

}

func (tg *TgClient) SendMsg(msgText string, chatId int64, replyId int) {
	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.ParseMode = "html"
	msg.ReplyToMessageID = replyId

	newMsg, err := tg.bot.Send(msg)
	if err != nil {
		log.Printf("%d %s", newMsg.MessageID, err.Error())
	}
}

func getReplyId(update tgbotapi.Update) int {
	var replyMsgId int
	if update.Message != nil && update.Message.ReplyToMessage != nil {
		replyMsgId = update.Message.ReplyToMessage.MessageID
	}
	return replyMsgId
}
