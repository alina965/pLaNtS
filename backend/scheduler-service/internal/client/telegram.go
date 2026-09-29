package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TelegramClient struct {
	baseURL string
	client  *http.Client
	key     string
}

func NewTelegramClient(baseURL string, timeout time.Duration, key string) *TelegramClient {
	return &TelegramClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: timeout},
		key:     key,
	}
}

type notifyRequest struct {
	UserID string `json:"userId"`
	Text   string `json:"text"`
}

func (c *TelegramClient) Notify(ctx context.Context, userID uuid.UUID, text string) error {
	payload, err := json.Marshal(notifyRequest{
		UserID: userID.String(),
		Text:   text,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/notify", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("X-Internal-Key", c.key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
