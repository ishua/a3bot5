package botcmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ishua/a3bot5/mcore/pkg/schema"
)

type UserSettings struct {
	Name     string
	Commands []string
}

type CmdRouter struct {
	userSettings []UserSettings
	selfQueue    string
	queue        queue
	telegram     telegramer
}

type queue interface {
	AddMsg(ctx context.Context, msg schema.TelegramMsg) error
}

type telegramer interface {
	SendMsg(msgText string, chatId int64, replyId int)
}

func NewCmdRouter(us []UserSettings, selfQueue string) *CmdRouter {
	return &CmdRouter{userSettings: us, selfQueue: selfQueue}
}

func (c *CmdRouter) RegQueue(q queue) {
	c.queue = q
}

func (c *CmdRouter) RegTelegram(t telegramer) {
	c.telegram = t
}

type Command struct {
	text    string
	channel string
}

func (c *CmdRouter) Send2Telegram(ctx context.Context, msg schema.TelegramMsg) {
	if c.telegram == nil {
		log.Fatal("cmdRouter doesn't have telegram")
	}
	c.telegram.SendMsg(msg.ReplyText, msg.ChatId, msg.ReplyToMessageID)
}

func (c *CmdRouter) Add2Queue(ctx context.Context, msg schema.TelegramMsg) error {
	if c.queue == nil {
		log.Fatal("cmdRouter doesn't have queue")
	}
	commandText := msg.Text
	if len(commandText) == 0 {
		if len(msg.Caption) == 0 {
			return fmt.Errorf("message havn't text command")
		}
		commandText = msg.Caption
	}

	allowCommands := c.getUserCommands(msg.UserName)
	if len(allowCommands) == 0 {
		return fmt.Errorf("botrouter: user %s doesn't has allow commands", msg.UserName)
	}

	command, err := c.getCommand(commandText, allowCommands)
	if err != nil {
		return err
	}
	msg.Command = command.text
	msg.QueueName = command.channel

	if command.text == "/help" {
		msg.ReplyText = getHelpText()
	}

	return c.queue.AddMsg(ctx, msg)
}

func (c *CmdRouter) getCommand(str string, allowCommands []string) (Command, error) {
	var cmd Command
	s := strings.Split(str, " ")
	if len(s) < 1 {
		return cmd, fmt.Errorf("getCommand: wrong command")
	}
	switch s[0] {
	case "/rate_usd", "usd":
		cmd.text = "/rate_usd"
		cmd.channel = "restjobs"
	case "/rate_eur", "eur":
		cmd.text = "/rate_eur"
		cmd.channel = "restjobs"
	case "/y2a", "y", "Y":
		cmd.text = "/y2a"
		cmd.channel = "ytd2feed"
	case "/torrent", "torrent", "t", "T":
		cmd.text = "/torrent"
		cmd.channel = "transmission"
	case "/note", "note", "n", "N", "Note":
		cmd.text = "/note"
		cmd.channel = "fsnotes"
	case "/help", "help", "h":
		cmd.text = "/help"
		cmd.channel = c.selfQueue
	}

	if cmd.text == "" {
		return cmd, fmt.Errorf("getCommand: command not found")
	}

	for _, command := range allowCommands {
		if command == cmd.text {
			return cmd, nil
		}
	}

	return Command{}, fmt.Errorf("getCommand: deny command")
}

func getHelpText() string {
	return "my commands: \n/rate_usd\n/rate_eur\n/y2a help\n/torrent help\n/note help"
}

func (c *CmdRouter) getUserCommands(userName string) []string {
	for _, user := range c.userSettings {
		if user.Name == userName {
			return user.Commands
		}
	}
	return []string{}
}
