package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kust1q/Zapp/main/internal/config"
	"github.com/segmentio/kafka-go"
)

type EventProducer struct {
	writers map[string]*kafka.Writer
}

func NewEventProducer(cfg *config.KafkaConfig) *EventProducer {
	producer := &EventProducer{
		writers: make(map[string]*kafka.Writer),
	}
	for _, topic := range cfg.Topics {
		producer.writers[topic] = &kafka.Writer{
			Addr:        kafka.TCP(cfg.Brokers...),
			Topic:       topic,
			Balancer:    &kafka.LeastBytes{},
			MaxAttempts: cfg.Producer.MaxRetries,
		}
	}
	return producer
}

func (p *EventProducer) Publish(ctx context.Context, topic string, event any) error {
	writer, exists := p.writers[topic]
	if !exists {
		return fmt.Errorf("topic %s not configured", topic)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	return writer.WriteMessages(ctx,
		kafka.Message{
			Value: data,
		},
	)
}

func (p *EventProducer) Close() error {
	for topic, writer := range p.writers {
		if err := writer.Close(); err != nil {
			return fmt.Errorf("close writer %s failed: %w", topic, err)
		}
	}
	return nil
}
