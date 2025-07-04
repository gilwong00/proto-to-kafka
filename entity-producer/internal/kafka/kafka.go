package kafka

import (
	"context"

	"github.com/gilwong00/proto-to-kafka/internal/config"

	kafkago "github.com/segmentio/kafka-go"
)

// Client represents a Kafka client that can publish messages to Kafka.
type Client interface {
	// Ping verifies connectivity to the Kafka cluster by establishing
	// a TCP connection to at least one broker.
	//
	// Returns a live connection to the broker if successful,
	// or an error if the connection fails.
	//
	// The caller is responsible for closing the returned connection
	// to release resources.
	Ping() (*kafkago.Conn, error)
	// PublishToKafka publishes messages or events to Kafka.
	// The context can be used to control cancellation and deadlines.
	Publish(ctx context.Context, eventName string, topic string, key []byte, value []byte) error
	// Close shuts down the client and closes any open resources.
	Close() error
	// GenerateKey creates a unique Kafka message key based on the given topic.
	// This key is typically used for partitioning and ensuring message ordering.
	//
	// Returns the generated key as a byte slice.
	GenerateKey(topic string) []byte
}

// NewClient creates a new instance of a Kafka Client using the provided configuration.
// The returned Client can be used to publish messages to Kafka.
// The config parameter should contain Kafka broker addresses, topic names, and other necessary settings.
func NewClient(config *config.Config) Client {
	cfg := &Config{
		Brokers: config.Brokers,
	}
	return newClient(cfg)
}
