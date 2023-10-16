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
	rater     GetRate
	publisher Publisher
}

func NewModel(rater GetRate, publisher Publisher) *Model {
	return &Model{rater: rater, publisher: publisher}
}

type GetRate interface {
	GetRate(valute string) string
}

type Publisher interface {
	Pub(ctx context.Context, payload string)
}

func (m *Model) DoJob(ctx context.Context, msgStr string) {
	var msg Msg
	err := msg.unmarshalBinary([]byte(msgStr))
	if err != nil {
		log.Println("[restjobs] Wrong unmarshal msg " + err.Error())
		return
	}

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

	replyText = m.rater.GetRate(valute)
	nMsg := Msg{
		Command:          msg.Command,
		UserName:         msg.UserName,
		MsgId:            msg.MsgId,
		ReplyToMessageID: msg.ReplyToMessageID,
		ChatId:           msg.ChatId,
		Text:             msg.Text,
		ReplyText:        replyText,
	}

	payload, err := nMsg.marshalBinary()
	if err != nil {
		log.Println("[doJob] marshal msg" + err.Error())
		return
	}

	m.publisher.Pub(ctx, string(payload))
}

func (m *Msg) marshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Msg) unmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	return nil
}
