package client

import (
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

func (c *PerenualClient) GetPlants(page int, q string) (*SpeciesListResponse, error) {
	params := url.Values{}
	params.Set("key", c.key)
	params.Set("indoor", "1")
	params.Set("page", strconv.Itoa(page))
	if q != "" {
		params.Set("q", q)
	}
	fullURL := perenualBaseUrl + perenualSpeciesListPath + "?" + params.Encode()

	var plants SpeciesListResponse
	err := c.getRequest(fullURL, &plants)
	if err != nil {
		return nil, err
	}

	return &plants, nil
}

func (c *PerenualClient) GetPlantDetails(id int) (*SpeciesDetailResponse, error) {
	fullURL := perenualBaseUrl + fmt.Sprintf(perenualDetailsPath, id) + "?key=" + c.key

	var plantDetail SpeciesDetailResponse
	err := c.getRequest(fullURL, &plantDetail)
	if err != nil {
		return nil, err
	}

	return &plantDetail, nil
}

func (c *PerenualClient) getRequest(url string, response any) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return err
	}

	return nil
}
