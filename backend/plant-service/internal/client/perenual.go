package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	perenualBaseUrl         = "https://perenual.com"
	perenualSpeciesListPath = "/api/v2/species-list"
	perenualDetailsPath     = "/api/v2/species/details/%d"
)

type PerenualClient struct {
	client *http.Client
	key    string
}

func NewPerenualClient(timeout time.Duration, key string) *PerenualClient {
	return &PerenualClient{client: &http.Client{Timeout: timeout}, key: key}
}

func (c *PerenualClient) GetPlants(ctx context.Context, page int, q string) (*SpeciesListResponse, error) {
	params := url.Values{}
	params.Set("key", c.key)
	params.Set("indoor", "1")
	params.Set("page", strconv.Itoa(page))
	if q != "" {
		params.Set("q", q)
	}
	fullURL := perenualBaseUrl + perenualSpeciesListPath + "?" + params.Encode()

	var plants SpeciesListResponse
	err := c.getRequest(ctx, fullURL, &plants)
	if err != nil {
		return nil, err
	}

	return &plants, nil
}

func (c *PerenualClient) GetPlantDetails(ctx context.Context, id int) (*SpeciesDetailResponse, error) {
	fullURL := perenualBaseUrl + fmt.Sprintf(perenualDetailsPath, id) + "?key=" + c.key

	var plantDetail SpeciesDetailResponse
	err := c.getRequest(ctx, fullURL, &plantDetail)
	if err != nil {
		return nil, err
	}

	return &plantDetail, nil
}

func (c *PerenualClient) getRequest(ctx context.Context, rawURL string, response any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}
