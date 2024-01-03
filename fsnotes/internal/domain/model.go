package domain

import (
	"context"
	"encoding/json"
	"log"
	"time"
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
	diary   Diarier
	publish Publisher
}

type Diarier interface {
	Add(ctx context.Context, now time.Time, theme string, text string) error
}

type Publisher interface {
	Pub(ctx context.Context, payload string)
}

func NewModel(aj Diarier, p Publisher) *Model {
	return &Model{
		diary:   aj,
		publish: p,
	}
}

func (m *Model) DoJob(ctx context.Context, msgStr string) {
	var msg Msg
	err := msg.unmarshalBinary([]byte(msgStr))
	if err != nil {
		log.Println("[fsnotes] Wrong unmarshal msg " + err.Error())
		return
	}

	var replyText string
	switch msg.Command {

	}

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

	m.publish.Pub(ctx, string(payload))
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
