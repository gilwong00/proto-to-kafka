package kafka

import (
	"context"

	"github.com/gilwong00/proto-to-kafka/internal/config"
)

// Client represents a Kafka client that can publish messages to Kafka.
type Client interface {
	// PublishToKafka publishes messages or events to Kafka.
	// The context can be used to control cancellation and deadlines.
	Publish(ctx context.Context, topic string, key []byte, value []byte) error

	// Close shuts down the client and closes any open resources.
	Close() error
}

// NewClient creates a new instance of a Kafka Client using the provided configuration.
// The returned Client can be used to publish messages to Kafka.
// The config parameter should contain Kafka broker addresses, topic names, and other necessary settings.
func NewClient(config *config.Config) Client {
	cfg := &Config{
		Brokers: nil,
	}
	return newClient(cfg)
}
