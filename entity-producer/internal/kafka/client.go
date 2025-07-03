package kafka

import (
	"context"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

// kafkaClient is a concrete implementation of the Client interface
// that publishes messages to Kafka using the provided configuration.
//
// It manages a shared kafka-go Writer for efficient reuse across multiple
// publish calls.
type kafkaClient struct {
	writer *kafkago.Writer
}

// Compile-time assertion to ensure kafkaClient implements the Client interface.
var _ Client = (*kafkaClient)(nil)

// newClient creates a new kafkaClient instance configured with the given Config.
// This constructor is unexported to enforce creation via NewClient.
//
// The returned client maintains a single underlying kafka-go Writer,
// which establishes connections lazily upon first use and reuses them
// for subsequent messages.
func newClient(config *Config) *kafkaClient {
	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(config.Brokers...),
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
		BatchTimeout: 10 * time.Millisecond,
	}

	return &kafkaClient{
		writer: writer,
	}
}

// Publish writes a message to the specified Kafka topic.
//
// This method is safe to call concurrently from multiple goroutines.
func (c *kafkaClient) Publish(
	ctx context.Context,
	topic string,
	key []byte,
	value []byte,
) error {
	msg := kafkago.Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	if err := c.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

// Close gracefully shuts down the Kafka writer and releases any resources.
//
// It should be called exactly once when the application is shutting down.
func (c *kafkaClient) Close() error {
	return c.writer.Close()
}
