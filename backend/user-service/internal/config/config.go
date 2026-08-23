package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAddress         = ":8080"
	defaultAccessTokenTTL  = time.Minute
	defaultRefreshTokenTTL = 7 * 24 * time.Hour
	defaultBcryptCost      = 10
)

type Config struct {
	Addr            string
	DatabaseURL     string
	JwtSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	BcryptCost      int
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
	accessTokenTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		accessTokenTTL = defaultAccessTokenTTL
	}
	refreshTokenTTL, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		refreshTokenTTL = defaultRefreshTokenTTL
	}
	bcryptCost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		bcryptCost = defaultBcryptCost
	}

	return &Config{
		Addr:            addr,
		DatabaseURL:     databaseURL,
		JwtSecret:       jwtSecret,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		BcryptCost:      bcryptCost,
	}, nil
}
