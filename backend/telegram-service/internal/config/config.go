package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAddress            = ":8084"
	defaultTimeout            = 40 * time.Second
	defaultKafkaBrokers       = "redpanda:9092"
	defaultKafkaTopicWatering = "watering.notify"
	defaultKafkaGroupID       = "telegram-notifier"
)

type Config struct {
	Addr               string
	Token              string
	BotUsername        string
	DatabaseURL        string
	JwtSecret          string
	InternalAPIKey     string
	Timeout            time.Duration
	KafkaBrokers       string
	KafkaTopicWatering string
	KafkaGroupID       string
}

func New(path string) (*Config, error) {
	_ = godotenv.Load(path)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddress
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	botUsername := os.Getenv("TELEGRAM_BOT_USERNAME")
	if botUsername == "" {
		return nil, errors.New("TELEGRAM_BOT_USERNAME environment variable not set")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable not set")
	}

	internalAPIKey := os.Getenv("INTERNAL_API_KEY")
	if internalAPIKey == "" {
		return nil, errors.New("INTERNAL_API_KEY environment variable not set")
	}

	timeout := defaultTimeout
	if v := os.Getenv("TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, err
		}
		timeout = d
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = defaultKafkaBrokers
	}

	kafkaTopicWatering := os.Getenv("KAFKA_TOPIC_WATERING")
	if kafkaTopicWatering == "" {
		kafkaTopicWatering = defaultKafkaTopicWatering
	}

	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = defaultKafkaGroupID
	}

	return &Config{
		Addr:               addr,
		Token:              token,
		BotUsername:        botUsername,
		DatabaseURL:        databaseURL,
		JwtSecret:          jwtSecret,
		InternalAPIKey:     internalAPIKey,
		Timeout:            timeout,
		KafkaBrokers:       kafkaBrokers,
		KafkaTopicWatering: kafkaTopicWatering,
		KafkaGroupID:       kafkaGroupID,
	}, nil
}
