package mcore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ishua/a3bot5/mcore/pkg/schema"
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

func (c *ClientMcore) AddMsg(ctx context.Context, msg schema.TelegramMsg) error {

	body, err := msg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("addmsg marshal err %w", err)
	}

	_, err = c.doPost(body) //TODO  сделать обработку полученного овтета, там может быть ошибка
	if err != nil {
		return fmt.Errorf("addmsg doPost: %w", err)
	}

	return nil
}

func (с *ClientMcore) GetMsg(ctx context.Context, queueName string) (schema.TelegramMsg, error) {
	return schema.TelegramMsg{}, nil
}

func (c *ClientMcore) doPost(body []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: c.timeout,
	}

	req, err := http.NewRequest("POST", c.addr, bytes.NewBuffer(body))
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
