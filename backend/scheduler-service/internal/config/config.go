package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAddress            = ":8083"
	defaultTimeout            = 30
	defaultLocation           = "Europe/Moscow"
	defaultKafkaBrokers       = "redpanda:9092"
	defaultKafkaTopicWatering = "watering.notify"
)

type Config struct {
	Addr               string
	UserPlantsURL      string
	TelegramURL        string
	Timeout            time.Duration
	InternalAPIKey     string
	Location           string
	KafkaBrokers       string
	KafkaTopicWatering string
}

func New(path string) (*Config, error) {
	_ = godotenv.Load(path)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = defaultAddress
	}

	userPlantsURL := os.Getenv("USER_PLANTS_URL")
	if userPlantsURL == "" {
		return nil, errors.New("USER_PLANTS_URL is required")
	}

	telegramURL := os.Getenv("TELEGRAM_URL")
	if telegramURL == "" {
		return nil, errors.New("TELEGRAM_URL is required")
	}

	var timeout int
	var err error

	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		timeout = defaultTimeout
	} else {
		timeout, err = strconv.Atoi(timeoutStr)
		if err != nil {
			return nil, err
		}
	}

	internalAPIKey := os.Getenv("INTERNAL_API_KEY")
	if internalAPIKey == "" {
		return nil, errors.New("INTERNAL_API_KEY environment variable not set")
	}

	location := os.Getenv("CRON_LOCATION")
	if location == "" {
		location = defaultLocation
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = defaultKafkaBrokers
	}

	kafkaTopicWatering := os.Getenv("KAFKA_TOPIC_WATERING")
	if kafkaTopicWatering == "" {
		kafkaTopicWatering = defaultKafkaTopicWatering
	}

	return &Config{
		Addr:               addr,
		UserPlantsURL:      userPlantsURL,
		TelegramURL:        telegramURL,
		Timeout:            time.Duration(timeout) * time.Second,
		InternalAPIKey:     internalAPIKey,
		Location:           location,
		KafkaBrokers:       kafkaBrokers,
		KafkaTopicWatering: kafkaTopicWatering,
	}, nil
}
