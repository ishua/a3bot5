package botcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/ishua/a3bot5/tbot/internal/schema"
)

type publisher interface {
	Pub(ctx context.Context, channel, value string) error
}

type CmdRouter struct {
	publisher publisher
}

func NewCmdRouter(p publisher) *CmdRouter {
	return &CmdRouter{publisher: p}
}

type Message struct {
	UserName         string
	MsgId            int
	ReplyToMessageID int
	ChatId           int64
	Text             string
	Caption          string
	ReplyText        string
	FileUrl          string
}

type Command struct {
	text    string
	channel string
}

func (c *CmdRouter) Send(ctx context.Context, msg Message, allowCommands []string, myChannel string) error {

	commandText := msg.Text
	if len(commandText) == 0 {
		if len(msg.Caption) == 0 {
			return fmt.Errorf("message havn't text command")
		}
		commandText = msg.Caption
	}
	command, err := getCommand(commandText, allowCommands, myChannel)
	if err != nil {
		return err
	}
	if command.text == "/help" {
		msg.ReplyText = getHelpText()
	}
	q := schema.ChannelMsg{
		Command:          command.text,
		UserName:         msg.UserName,
		MsgId:            msg.MsgId,
		ReplyToMessageID: msg.ReplyToMessageID,
		ChatId:           msg.ChatId,
		Text:             msg.Text,
		ReplyText:        msg.ReplyText,
		Caption:          msg.Caption,
		FileUrl:          msg.FileUrl,
	}

	payload, err := q.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal err:" + err.Error())
	}

	return c.publisher.Pub(ctx, command.channel, string(payload))
}

func getCommand(str string, allowCommands []string, myChannel string) (Command, error) {
	var c Command
	s := strings.Split(str, " ")
	if len(s) < 1 {
		return c, fmt.Errorf("wrong command")
	}
	switch s[0] {
	case "/rate_usd", "usd":
		c.text = "/rate_usd"
		c.channel = "restjobs"
	case "/rate_eur", "eur":
		c.text = "/rate_eur"
		c.channel = "restjobs"
	case "/y2a", "y", "Y":
		c.text = "/y2a"
		c.channel = "ytd2feed"
	case "/torrent", "torrent", "t", "T":
		c.text = "/torrent"
		c.channel = "transmission"
	case "/note", "note", "n", "N", "Note":
		c.text = "/note"
		c.channel = "fsnotes"
	case "/help", "help", "h":
		c.text = "/help"
		c.channel = myChannel
	}

	if c.text == "" {
		return c, fmt.Errorf("command not found")
	}

	for _, command := range allowCommands {
		if command == c.text {
			return c, nil
		}
	}

	return Command{}, fmt.Errorf("deny command")
}

func getHelpText() string {
	return "my commands: \n/rate_usd\n/rate_eur\n/y2a\n/torrent\n/note"
}
