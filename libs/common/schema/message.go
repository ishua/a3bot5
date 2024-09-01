package schema

import "encoding/json"

type ChannelMsg struct {
	Command          string `json:"command"`
	UserName         string `json:"userName"`
	MsgId            int    `json:"msgId"`
	ReplyToMessageID int    `json:"replyToMessageID"`
	ChatId           int64  `json:"chatId"`
	Text             string `json:"text"`
	ReplyText        string `json:"replyText"`
	Caption          string `json:"caption"`
	FileUrl          string `json:"fileUrl"`
}

func (m *ChannelMsg) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *ChannelMsg) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	return nil
}
