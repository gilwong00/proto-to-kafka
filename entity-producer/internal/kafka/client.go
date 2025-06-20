package kafka

import (
	"context"

	"github.com/gilwong00/proto-to-kafka/internal/config"
)

// kafkaClient is a concrete implementation of the Client interface
// that publishes messages to Kafka using the provided configuration.
type kafkaClient struct {
	config *config.Config
}

// Compile-time assertion to ensure kafkaClient implements Client interface.
var _ Client = (*kafkaClient)(nil)

// newClient creates a new kafkaClient instance configured with the given config.
// It is unexported to enforce creation via NewClient.
func newClient(config *config.Config) *kafkaClient {
	// TODO: create new kafka instance here
	return &kafkaClient{
		config: config,
	}
}

// PublishToKafka publishes messages or events to Kafka.
func (k *kafkaClient) PublishToKafka(ctx context.Context) {
	// TODO: implement Kafka publishing logic here
}
