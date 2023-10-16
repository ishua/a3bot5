package domain

import (
	"context"
	"encoding/json"
	"log"
)

type Msg struct {
	Command          string `json:"command"`
	UserName         string `json:"userName"`
	MsgId            int    `json:"msgId"`
	ReplyToMessageID int    `json:"replyToMessageID"`
	ChatId           int64  `json:"chatId"`
	Text             string `json:"text"`
	ReplyText        string `json:"replyText"`
}

type Model struct {
	getRate GetRate
}

func NewModel(getRate GetRate) *Model {
	return &Model{getRate: getRate}
}

type GetRate interface {
	GetRate(valute string) string
}

func (m *Model) DoJob(ctx context.Context, msg Msg) Msg {
	var replyText string
	var valute string
	switch msg.Command {
	case "/rate_usd":
		valute = "USD"
	case "/rate_eur":
		valute = "EUR"
	}

	if valute == "" {
		errText := "[restjobs] Wrong command to get rate"
		replyText = errText
		log.Println(errText)
	}

	replyText = m.getRate.GetRate(valute)
	nMsg := Msg{
		Command:          msg.Command,
		UserName:         msg.UserName,
		MsgId:            msg.MsgId,
		ReplyToMessageID: msg.ReplyToMessageID,
		ChatId:           msg.ChatId,
		Text:             msg.Text,
		ReplyText:        replyText,
	}

	return nMsg

}

const (
	cbr_url = "https://www.cbr-xml-daily.ru/daily_json.js"
)

type cbr_response struct {
	Date   string
	Valute valutes
}

type valutes struct {
	USD valute
	EUR valute
}

type valute struct {
	CharCode string
	Value    json.Number `json:"Value"`
}
