package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/alina965/pLaNtS/telegram-service/internal/domain"
	"github.com/alina965/pLaNtS/telegram-service/internal/telegram"
	"github.com/segmentio/kafka-go"
)

type WateringConsumer struct {
	reader  *kafka.Reader
	service *telegram.Service
}

func NewWateringConsumer(brokers, topic, groupID string, service *telegram.Service) *WateringConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokers},
		Topic:   topic,
		GroupID: groupID,
	})
	return &WateringConsumer{reader: reader, service: service}
}

func (c *WateringConsumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			log.Printf("kafka fetch: %v", err)
			continue
		}

		var event domain.WateringNotifyEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("kafka unmarshal: %v", err)
			if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
				log.Printf("kafka commit: %v", commitErr)
			}
			continue
		}

		if err := c.service.Notify(ctx, event.UserID, event.Text); err != nil {
			log.Printf("kafka notify: %v", err)
			
			if errors.Is(err, telegram.ErrNotLinked) {
				if commitErr := c.reader.CommitMessages(ctx, msg); commitErr != nil {
					log.Printf("kafka commit: %v", commitErr)
				}
				continue
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("kafka commit: %v", err)
		}
	}
}

func (c *WateringConsumer) Close() error {
	return c.reader.Close()
}
