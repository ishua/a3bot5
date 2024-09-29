package domain

import (
	"context"

	"github.com/ishua/a3bot5/mcore/pkg/schema"
)

type AddQueue interface {
	AddTelegramMsg(ctx context.Context, msg schema.TelegramMsg) error
}

type GetQueue interface {
	GetTelegramMsg(ctx context.Context, queueName string) (schema.TelegramMsg, error)
}

type MyDomain struct {
	qadder  AddQueue
	qgetter GetQueue
}

func NewMyDomain(qadder AddQueue, qgetter GetQueue) *MyDomain {
	return &MyDomain{
		qadder:  qadder,
		qgetter: qgetter,
	}
}

func (md *MyDomain) AddMessageToQueue(ctx context.Context, msg schema.TelegramMsg) error {
	return md.qadder.AddTelegramMsg(ctx, msg)
}

func (md *MyDomain) GetMessageFromQueue(ctx context.Context, queueName string) (schema.TelegramMsg, error) {
	return md.qgetter.GetTelegramMsg(ctx, queueName)
}
