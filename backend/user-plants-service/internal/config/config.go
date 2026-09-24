package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

const defaultAddress = ":8082"

type Config struct {
	Addr        string
	DatabaseURL string
	JwtSecret   string
}

func New(path string) (*Config, error) {
	_ = godotenv.Load(path)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddress
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable not set")
	}

	return &Config{
		Addr:        addr,
		DatabaseURL: databaseURL,
		JwtSecret:   jwtSecret,
	}, nil
}
