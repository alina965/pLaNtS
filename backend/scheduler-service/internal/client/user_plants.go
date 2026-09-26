package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alina965/pLaNtS/scheduler-service/internal/domain"
	"github.com/google/uuid"
)

type UserPlantsClient struct {
	baseURL string
	client  *http.Client
	key     string
}

func NewUserPlantsClient(baseURL string, timeout time.Duration, key string) *UserPlantsClient {
	return &UserPlantsClient{baseURL: baseURL, client: &http.Client{Timeout: timeout}, key: key}
}

func (c *UserPlantsClient) ListNeedingWater(ctx context.Context) ([]*domain.UserPlant, error) {
	path := c.baseURL + "/internal/user-plants/needing-water"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Internal-Key", c.key)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var listNeedingWater []*domain.UserPlant

	err = json.NewDecoder(resp.Body).Decode(&listNeedingWater)
	if err != nil {
		return nil, err
	}

	return listNeedingWater, nil
}

func (c *UserPlantsClient) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	path := c.baseURL + "/internal/user-plants/" + id.String() + "/status"

	payload, err := json.Marshal(domain.UpdateStatusRequest{Status: &status})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, path, bytes.NewReader(payload))
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
