package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAddress         = ":8080"
	defaultTimeout         = time.Second * 10
	defaultCacheTTLList    = time.Hour
	defaultCacheTTLDetails = time.Hour * 24
)

type Config struct {
	Addr            string
	Timeout         time.Duration
	PerenualKey     string
	RedisAddr       string
	CacheTTLList    time.Duration
	CacheTTLDetails time.Duration
}

func New(path string) (*Config, error) {
	_ = godotenv.Load(path)

	perenualKey := os.Getenv("PERENUAL_KEY")
	if perenualKey == "" {
		return nil, errors.New("PERENUAL_KEY environment variable not set")
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddress
	}

	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		timeout = defaultTimeout
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, errors.New("REDIS_ADDR environment variable not set")
	}

	cacheTTLList, err := time.ParseDuration(os.Getenv("CACHE_TTL_LIST"))
	if err != nil {
		cacheTTLList = defaultCacheTTLList
	}

	cacheTTLDetails, err := time.ParseDuration(os.Getenv("CACHE_TTL_DETAILS"))
	if err != nil {
		cacheTTLDetails = defaultCacheTTLDetails
	}

	return &Config{
		Addr:            addr,
		Timeout:         timeout,
		PerenualKey:     perenualKey,
		RedisAddr:       redisAddr,
		CacheTTLList:    cacheTTLList,
		CacheTTLDetails: cacheTTLDetails,
	}, nil
}
