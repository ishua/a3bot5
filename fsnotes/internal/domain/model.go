package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

const helpText = `This is help for /note command
 - /note diary add 5bx <text>
 - /note diary add entry <text>
 - /note diary help
 - synonyms: diary = d, add=a, entry=e, 5bx=5`

type Note struct {
	Theme string
	Text  string
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

	note, err := msg.ParseCommand()
	if err != nil {
		replyText = err.Error()
	} else {
		err = m.diary.Add(ctx, time.Now(), note.Theme, note.Text)
		if err != nil {
			replyText = "add notes err: " + err.Error()
		} else {
			replyText = "notes added"
		}
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

func (m *Msg) ParseCommand() (Note, error) {
	// note diary add 5bx 2
	// note diary add entry
	// note diary list week/month/day
	// note help
	texts := strings.Split(m.Text, " ")
	if len(texts) <= 2 {
		return Note{}, fmt.Errorf("Wrong fsnotes command")
	}
	if texts[1] == "help" {
		return Note{}, fmt.Errorf(helpText)
	}
	if texts[1] == "diary" || texts[1] == "d" {
		if texts[2] == "add" || texts[2] == "a" {
			if len(texts) < 5 {
				return Note{}, fmt.Errorf("Wrong fsnotes command")
			}
			switch texts[3] {
			case "5bx", "5":
				return Note{
					Theme: "5bx",
					Text:  TextClean(texts),
				}, nil
			case "entry", "e":
				return Note{
					Theme: "entry",
					Text:  TextClean(texts),
				}, nil
			}

		}
		if texts[2] == "list" || texts[2] == "l" {
			return Note{}, fmt.Errorf("Not working")
		}
	}
	return Note{}, fmt.Errorf("fsnotes doesn't have this command")
}

func TextClean(texts []string) string {
	var ret string
	for i, s := range texts[4:] {
		s = strings.ReplaceAll(s, "\n", " ")
		ret = ret + s
		if len(texts) != i+5 {
			ret = ret + " "
		}
	}
	return ret
}
