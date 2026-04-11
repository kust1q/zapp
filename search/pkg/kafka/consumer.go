package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kust1q/Zapp/search/internal/config"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Handler func(ctx context.Context, topic string, data []byte) error

type eventConsumer struct {
	readers map[string]*kafka.Reader
}

func NewEventConsumer(cfg *config.KafkaConfig) *eventConsumer {
	c := &eventConsumer{
		readers: make(map[string]*kafka.Reader),
	}

	for _, topic := range cfg.Topics {
		c.readers[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:          cfg.Brokers,
			GroupID:          cfg.Consumer.GroupID,
			Topic:            topic,
			MinBytes:         10 << 10,
			MaxBytes:         10 << 20,
			MaxWait:          500 * time.Millisecond,
			ReadBatchTimeout: 1 * time.Second,
		})
	}

	return c
}

func (c *eventConsumer) Run(ctx context.Context, handler Handler) error {
	var wg sync.WaitGroup

	for topic := range c.readers {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			c.runReader(ctx, t, handler)
		}(topic)
	}

	<-ctx.Done()
	wg.Wait()
	return nil
}

func (c *eventConsumer) runReader(ctx context.Context, topic string, handler Handler) {
	reader := c.readers[topic]

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logrus.WithFields(logrus.Fields{
				"topic": topic,
				"error": err,
			}).Error("kafka read error")
			time.Sleep(time.Second)
			continue
		}

		if err := handler(ctx, topic, msg.Value); err != nil {
			logrus.WithFields(logrus.Fields{
				"topic": topic,
				"error": err,
			}).Error("handler error")
		}
	}
}

func (c *eventConsumer) Close() error {
	var errs []error
	for topic, reader := range c.readers {
		if err := reader.Close(); err != nil {
			errs = append(errs, fmt.Errorf("topic %s: %w", topic, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}
	return nil
}
