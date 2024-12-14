package schema

import (
	"encoding/json"
	"strconv"
)

type TelegramMsg struct {
	Command          string `json:"command"`
	UserName         string `json:"userName"`
	MsgId            int    `json:"msgId"`
	ReplyToMessageID int    `json:"replyToMessageID"`
	ChatId           int64  `json:"chatId"`
	Text             string `json:"text"`
	ReplyText        string `json:"replyText"`
	Caption          string `json:"caption"`
	FileUrl          string `json:"fileUrl"`
	QueueName        string `json:"queueName"`
}

func (m *TelegramMsg) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *TelegramMsg) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	return nil
}

func (m *TelegramMsg) String() string {
	return "MsgId: " + strconv.Itoa(m.MsgId)
}

type GetMsgReq struct {
	QueueName string `json:"queueName"`
}

type GetMsgRes struct {
	Data   TelegramMsg `json:"telegramMsg"`
	Empty  bool        `json:"empty"`
	Error  string      `json:"error"`
	Status string      `json:"status"`
}

type Resp struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}

type AddMsgReq struct {
	TelegramMsg
}
