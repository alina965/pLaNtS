package kafka

import (
	"context"
	"encoding/json"

	"github.com/alina965/pLaNtS/scheduler-service/internal/domain"
	"github.com/segmentio/kafka-go"
)

type WateringProducer struct {
	writer *kafka.Writer
}

func NewWateringProducer(brokers string, topic string) *WateringProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &WateringProducer{writer}
}

func (p *WateringProducer) Publish(ctx context.Context, event domain.WateringNotifyEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{Value: data})
}

func (p *WateringProducer) Close() error {
	return p.writer.Close()
}
