package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	wikipediaUrl = "https://en.wikipedia.org/api/rest_v1/page/summary/"
)

type WikipediaClient struct {
	client *http.Client
}

func NewWikipediaClient(timeout time.Duration) *WikipediaClient {
	return &WikipediaClient{client: &http.Client{Timeout: timeout}}
}

func (c *WikipediaClient) GetDescription(name string) *PlantDescriptionResponse {
	normalizedName := url.PathEscape(strings.ReplaceAll(strings.TrimSpace(name), " ", "_"))

	fullUrl := wikipediaUrl + normalizedName

	request, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil
	}
	request.Header.Set("User-Agent", "pLaNtS/0.1 https://github.com/alina965/pLaNtS")

	resp, err := c.client.Do(request)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var plantDescription PlantDescriptionResponse
	err = json.NewDecoder(resp.Body).Decode(&plantDescription)
	if err != nil {
		return nil
	}

	return &plantDescription
}
