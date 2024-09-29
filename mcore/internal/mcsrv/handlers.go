package mcsrv

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ishua/a3bot5/mcore/internal/domain"
	"github.com/ishua/a3bot5/mcore/pkg/schema"
)

type Handlers struct {
	md *domain.MyDomain
}

func NewSrvHandlers(md *domain.MyDomain) *Handlers {
	return &Handlers{md: md}
}

func (s *Handlers) AddMsg(w http.ResponseWriter, req *http.Request) {

	type AddMsgReq struct {
		schema.TelegramMsg
	}

	type AddMsgRes struct {
		Error string `json:"error"`
	}

	var msg AddMsgReq
	err := json.NewDecoder(req.Body).Decode(&msg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	ctx := context.Background()
	err = s.md.AddMessageToQueue(ctx, msg.TelegramMsg)
	w.Header().Set("Content-Type", "application/json")
	var res AddMsgRes
	if err != nil {
		res.Error = err.Error()
		js, _ := json.Marshal(res)
		w.Write(js)
	}
	js, _ := json.Marshal(res)
	w.Write(js)
}

// func (s *Handlers) GetMsg(w http.ResponseWriter, req *http.Request) {

// 	type GetMsgReq struct {
// 		QueueName string `json:"queueName"`
// 	}

// 	type GetMsgRes struct {
// 		Data  schema.TelegramMsg `json:"telegramMsg"`
// 		Error string             `json:"error"`
// 	}
// }
