package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	kafkago "github.com/segmentio/kafka-go"
)

// kafkaClient is a concrete implementation of the Client interface
// that publishes messages to Kafka using the provided configuration.
//
// It manages a shared kafka-go Writer for efficient reuse across multiple
// publish calls.
type kafkaClient struct {
	brokers []string
	writer  *kafkago.Writer
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
		brokers: config.Brokers,
		writer:  writer,
	}
}

// Ping attempts to open a TCP connection to the first Kafka broker
// to verify that it is reachable and accepting connections.
// It returns an error if the connection fails.
func (c *kafkaClient) Ping() error {
	if len(c.brokers) == 0 {
		return fmt.Errorf("no Kafka brokers configured")
	}
	conn, err := kafka.Dial("tcp", c.brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka broker: %w", err)
	}
	defer conn.Close()
	// Optional: can add a timeout or metadata check here
	return nil
}

// Publish writes a message to the specified Kafka topic.
//
// This method is safe to call concurrently from multiple goroutines.
func (c *kafkaClient) Publish(
	ctx context.Context,
	eventName string,
	topic string,
	key []byte,
	value []byte,
) error {
	msg := kafkago.Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Time:  time.Now(),
		Headers: []kafka.Header{
			{Key: "eventType", Value: []byte(eventName)},
		},
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

// GenerateKey generates a unique Kafka message key by concatenating
// the given topic string with a newly generated UUID.
// The returned key is a byte slice suitable for use as a Kafka message key.
//
// Example output: "topic-f47ac10b-58cc-4372-a567-0e02b2c3d479"
func (c *kafkaClient) GenerateKey(topic string) []byte {
	return []byte(fmt.Sprintf("%s-%s", topic, uuid.NewString()))
}
