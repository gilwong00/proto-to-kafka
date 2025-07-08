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

	// Publish publishes a message or event to Kafka.
	// The context can be used to control cancellation and deadlines.
	//
	// If the target topic does not exist, Publish attempts to create it.
	// In case of transient errors or failures to write, the client
	// internally retries several times before enqueuing the message
	// for asynchronous retry in the background.
	//
	// Returns an error only if the message cannot be published immediately
	// or enqueued for retry.
	Publish(ctx context.Context, eventName string, topic string, key []byte, value []byte) error

	// Close gracefully shuts down the client, closes any open connections,
	// and stops any background retry workers.
	//
	// Close must be called exactly once during application shutdown to
	// avoid resource leaks.
	Close() error

	// GenerateKey creates a unique Kafka message key based on the given topic.
	// This key is typically used for partitioning and ensuring message ordering.
	//
	// Returns the generated key as a byte slice.
	GenerateKey(topic string) []byte

	// ListTopics returns a slice of all topic names currently available
	// in the Kafka cluster.
	//
	// It establishes a connection to a broker each time it is called.
	ListTopics() ([]string, error)
}

// NewClient creates a new instance of a Kafka Client using the provided configuration.
// The returned Client can be used to publish messages to Kafka.
// The config parameter should contain Kafka broker addresses, topic names, and other necessary settings.
func NewClient(config *config.Config) (Client, error) {
	cfg := &Config{
		Brokers: config.Brokers,
	}
	return newClient(cfg)
}
