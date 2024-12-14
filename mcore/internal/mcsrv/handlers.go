package mcsrv

import (
	"context"
	"encoding/json"
	"fmt"
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

	var amr schema.AddMsgReq
	err := json.NewDecoder(req.Body).Decode(&amr)
	var resp []byte
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		resp = getErrResp(fmt.Errorf("body AddMsg decode err: %w", err))
		w.Write(resp)
		return
	}

	ctx := context.Background()
	err = s.md.AddMessageToQueue(ctx, amr.TelegramMsg)

	if err != nil {
		resp = getErrResp(fmt.Errorf("add to queue err: %w", err))
		w.Write(resp)
		return
	}
	w.Write(getOkResp())
}

func (s *Handlers) GetMsg(w http.ResponseWriter, req *http.Request) {

	var gmr schema.GetMsgReq
	err := json.NewDecoder(req.Body).Decode(&gmr)
	var resp []byte
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		resp = getErrResp(fmt.Errorf("body GetMsg decode err: %w", err))
		w.Write(resp)
		return
	}

	ctx := context.Background()
	msg, err := s.md.GetMessageFromQueue(ctx, gmr.QueueName)
	if err != nil {
		if err.Error() == "redis: nil" {
			js, err := json.Marshal(schema.GetMsgRes{Empty: true, Status: "OK"})
			if err != nil {
				resp = getErrResp(fmt.Errorf("marshal msg err: %w", err))
				w.Write(resp)
				return
			}
			w.Write(js)
			return
		}
		resp = getErrResp(fmt.Errorf("get from queue: %s err: %w", gmr.QueueName, err))
		w.Write(resp)
		return
	}

	js, err := json.Marshal(schema.GetMsgRes{Data: msg, Status: "OK"})
	if err != nil {
		resp = getErrResp(fmt.Errorf("marshal msg err: %w", err))
		w.Write(resp)
		return
	}
	w.Write(js)
}

func getOkResp() []byte {
	ok, _ := json.Marshal(schema.Resp{
		Status: "OK",
	})
	return ok
}

func getErrResp(err error) []byte {
	log.Println("handler: " + err.Error())
	ok, _ := json.Marshal(schema.Resp{
		Error:  err.Error(),
		Status: "error",
	})
	return ok
}
