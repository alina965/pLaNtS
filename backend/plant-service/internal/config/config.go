package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAddress = ":8080"
	defaultTimeout = time.Second * 10
)

type Config struct {
	Addr        string
	Timeout     time.Duration
	PerenualKey string
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

	return &Config{
		Addr:        addr,
		Timeout:     timeout,
		PerenualKey: perenualKey,
	}, nil
}
