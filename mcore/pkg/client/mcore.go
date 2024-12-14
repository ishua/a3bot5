package mcoreclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ishua/a3bot5/mcore/pkg/schema"
)

const (
	addMsgUrl = "/add-msg"
	getMsgUrl = "/get-msg"
	timeOut   = 60
)

type ClientMcore struct {
	addr    string
	secret  string
	timeout time.Duration
}

func NewClienMcore(addr, secret string) *ClientMcore {
	return &ClientMcore{
		addr:    addr,
		secret:  secret,
		timeout: 10 * time.Second,
	}
}

type Telegramer interface {
	Send2Telegram(ctx context.Context, msg schema.TelegramMsg)
}

func (c *ClientMcore) AddMsg(ctx context.Context, msg schema.TelegramMsg) error {

	body, err := msg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("addmsg marshal err %w", err)
	}

	respByte, err := c.doPost(addMsgUrl, body)
	if err != nil {
		return fmt.Errorf("addmsg doPost: %w", err)
	}

	var resp schema.Resp
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		return fmt.Errorf("addmsg unmarshal resp: %w", err)
	}

	if resp.Status != "OK" {
		return fmt.Errorf("addmsg err: %s", resp.Error)
	}

	return nil
}

func (с *ClientMcore) GetMsg(ctx context.Context, queueName string) (schema.TelegramMsg, error) {
	var ret schema.TelegramMsg
	getMsgReq := schema.GetMsgReq{
		QueueName: queueName,
	}
	body, err := json.Marshal(getMsgReq)
	if err != nil {
		return ret, fmt.Errorf("getMsg marshal body: %w", err)
	}

	respByte, err := с.doPost(getMsgUrl, body)
	if err != nil {
		return ret, fmt.Errorf("getMsg doPost: %w", err)
	}

	var resp schema.GetMsgRes
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		return ret, fmt.Errorf("getMsg Unmarshal: %w", err)
	}

	if resp.Error != "" {
		return ret, fmt.Errorf("getMsg error status: %s", resp.Error)
	}

	return resp.Data, nil
}

func (c *ClientMcore) doPost(addr string, body []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: c.timeout,
	}

	req, err := http.NewRequest("POST", c.addr+addr, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("doPost NewRequest err %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("secret", c.secret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doPost http request %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("doPost some error status %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doPost read body %w", err)
	}
	return respBody, nil
}

func (c *ClientMcore) ListenGetMsg(ctx context.Context, t Telegramer, queueName string) error {
	timeout := time.Duration(timeOut * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				{
					log.Println("stopping listen telegram")
					return
				}
			default:
				{
					msg, err := c.GetMsg(ctx, queueName)
					if err != nil {
						log.Printf("listen %s err: %s", queueName, err.Error())
						continue
					}
					t.Send2Telegram(ctx, msg)
					time.Sleep(timeout)
				}
			}
		}
	}()

	return nil
}
