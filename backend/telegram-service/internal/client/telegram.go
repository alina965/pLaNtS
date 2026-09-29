package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const telegramAPIBase = "https://api.telegram.org"

type BotClient struct {
	token      string
	httpClient *http.Client
}

func NewBotClient(token string, timeout time.Duration) *BotClient {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &BotClient{
		token:      token,
		httpClient: &http.Client{Timeout: timeout},
	}
}

type sendMessageRequest struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type telegramResponse struct {
	OK          bool            `json:"ok"`
	Description string          `json:"description"`
	Result      json.RawMessage `json:"result"`
}

func (c *BotClient) SendMessage(ctx context.Context, chatID int64, text string) error {
	payload, err := json.Marshal(sendMessageRequest{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramAPIBase, c.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var body telegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}
	if !body.OK {
		if body.Description != "" {
			return fmt.Errorf("telegram sendMessage failed: %s", body.Description)
		}
		return fmt.Errorf("telegram sendMessage failed with status %d", resp.StatusCode)
	}

	return nil
}

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Text      string `json:"text"`
	Chat      Chat   `json:"chat"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type getUpdatesRequest struct {
	Offset  int64 `json:"offset"`
	Timeout int   `json:"timeout"`
}

func (c *BotClient) GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]Update, error) {
	if timeoutSec <= 0 {
		timeoutSec = 30
	}

	payload, err := json.Marshal(getUpdatesRequest{
		Offset:  offset,
		Timeout: timeoutSec,
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/bot%s/getUpdates", telegramAPIBase, c.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body struct {
		OK          bool     `json:"ok"`
		Description string   `json:"description"`
		Result      []Update `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if !body.OK {
		if body.Description != "" {
			return nil, fmt.Errorf("telegram getUpdates failed: %s", body.Description)
		}
		return nil, fmt.Errorf("telegram getUpdates failed with status %d", resp.StatusCode)
	}

	return body.Result, nil
}
