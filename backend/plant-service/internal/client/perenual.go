package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/alina965/pLaNtS/plant-service/internal/cache"
)

const (
	perenualBaseUrl         = "https://perenual.com"
	perenualSpeciesListPath = "/api/v2/species-list"
	perenualDetailsPath     = "/api/v2/species/details/%d"
)

type PerenualClient struct {
	client          *http.Client
	key             string
	cache           *cache.RedisCache
	cacheTTLList    time.Duration
	cacheTTLDetails time.Duration
}

func NewPerenualClient(timeout time.Duration, key string, cache *cache.RedisCache, cacheTTLList time.Duration, cacheTTLDetails time.Duration) *PerenualClient {
	return &PerenualClient{client: &http.Client{Timeout: timeout}, key: key, cache: cache, cacheTTLList: cacheTTLList, cacheTTLDetails: cacheTTLDetails}
}

func (c *PerenualClient) GetPlants(ctx context.Context, page int, q string) (*SpeciesListResponse, error) {
	if c.cache != nil {
		key := fmt.Sprintf("perenual:list:page=%d:q=%s", page, q)
		res, err := c.cache.Get(ctx, key)
		if err != nil {
			log.Print("error during get redis operation: ", err)
		} else if res != nil {
			var plants SpeciesListResponse
			if err := json.Unmarshal(res, &plants); err == nil {
				return &plants, nil
			}
		}
	}

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

	if c.cache != nil {
		key := fmt.Sprintf("perenual:list:page=%d:q=%s", page, q)
		value, err := json.Marshal(plants)
		if err != nil {
			log.Print("redis marshal: ", err)
		} else if err := c.cache.Set(ctx, key, value, c.cacheTTLList); err != nil {
			log.Print("redis set: ", err)
		}
	}

	return &plants, nil
}

func (c *PerenualClient) GetPlantDetails(ctx context.Context, id int) (*SpeciesDetailResponse, error) {
	if c.cache != nil {
		key := fmt.Sprintf("perenual:details:id=%d", id)
		res, err := c.cache.Get(ctx, key)
		if err != nil {
			log.Print("error during get redis operation: ", err)
		} else if res != nil {
			var plants SpeciesDetailResponse
			if err := json.Unmarshal(res, &plants); err == nil {
				return &plants, nil
			}
		}
	}
	fullURL := perenualBaseUrl + fmt.Sprintf(perenualDetailsPath, id) + "?key=" + c.key

	var plantDetail SpeciesDetailResponse
	err := c.getRequest(ctx, fullURL, &plantDetail)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		key := fmt.Sprintf("perenual:details:id=%d", id)
		value, err := json.Marshal(plantDetail)
		if err != nil {
			log.Print("redis marshal: ", err)
		} else if err := c.cache.Set(ctx, key, value, c.cacheTTLDetails); err != nil {
			log.Print("redis set: ", err)
		}
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
