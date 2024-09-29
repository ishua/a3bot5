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

	var amr AddMsgReq
	err := json.NewDecoder(req.Body).Decode(&amr)
	if err != nil {
		errResp(w, "body AddMsg decode err: "+err.Error())
		return
	}

	ctx := context.Background()
	err = s.md.AddMessageToQueue(ctx, amr.TelegramMsg)
	if err != nil {
		errResp(w, "add to queue err: "+err.Error())
		return
	}
}

func (s *Handlers) GetMsg(w http.ResponseWriter, req *http.Request) {

	type GetMsgReq struct {
		QueueName string `json:"queueName"`
	}

	type GetMsgRes struct {
		Data  schema.TelegramMsg `json:"telegramMsg"`
		Emty  bool               `json:"emty"`
		Error string             `json:"error"`
	}

	var gmr GetMsgReq
	err := json.NewDecoder(req.Body).Decode(&gmr)
	if err != nil {
		errResp(w, "body GetMsg decode err: "+err.Error())
		return
	}

	ctx := context.Background()
	msg, err := s.md.GetMessageFromQueue(ctx, gmr.QueueName)
	if err != nil {
		if err.Error() == "redis: nil" {
			js, err := json.Marshal(GetMsgRes{Emty: true})
			if err != nil {
				errResp(w, "marshal msg err: "+err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
		errResp(w, "get from queue:"+gmr.QueueName+" err: "+err.Error())
		return
	}

	js, err := json.Marshal(GetMsgRes{Data: msg})
	if err != nil {
		errResp(w, "marshal msg err: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func errResp(w http.ResponseWriter, err string) {
	log.Println("handler: " + err)
	type errResp struct {
		Error string `json:"error"`
	}
	w.Header().Set("Content-Type", "application/json")

	js, _ := json.Marshal(errResp{
		Error: err,
	})
	w.Write(js)
}
