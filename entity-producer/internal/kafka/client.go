package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"

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
func (c *kafkaClient) Ping() (*kafkago.Conn, error) {
	if len(c.brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers configured")
	}
	conn, err := kafkago.Dial("tcp", c.brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka broker: %w", err)
	}
	return conn, nil
}

// Publish writes a message to the specified Kafka topic.
//
// This method is safe to call concurrently from multiple goroutines.
// If the topic does not exist, it attempts to create it and retries once.
//
// It includes a short delay after topic creation to allow Kafka to
// initialize topic metadata before retrying.
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
		Headers: []kafkago.Header{
			{Key: "eventType", Value: []byte(eventName)},
		},
	}

	if err := c.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("Failed to write to Kafka topic %q: %v", topic, err)

		// We check for `isUnknownTopicError` because Kafka returns this error
		// when attempting to write to a topic that does not yet exist.
		// In many production setups, topics must be created explicitly before use.
		// However, if auto-creation is enabled or you want to create topics on-demand,
		// this code attempts to create the topic dynamically.
		if isUnknownTopicError(err) {
			log.Printf("Topic %q not found, attempting to create it...", topic)
			if createErr := c.createTopic(topic, 1, 1); createErr != nil {
				return fmt.Errorf("failed to create missing topic %q: %w", topic, createErr)
			}

			// Kafka topic creation is asynchronous and metadata propagation
			// to all brokers and clients may take a short moment.
			// We wait briefly here to allow the topic to be fully initialized
			// before retrying the publish, reducing "Unknown Topic Or Partition" errors.
			time.Sleep(1 * time.Second)
			// Retry publishing the message once
			// For more robust handling, consider exponential backoff or other retry strategies.
			retryErr := c.writer.WriteMessages(ctx, msg)
			if retryErr != nil && isUnknownTopicError(retryErr) {
				// If still unknown topic, maybe log and fail gracefully
				log.Printf("Retry failed: topic %q still unknown after creation attempt", topic)
				return fmt.Errorf("failed to re-publish message after topic creation: %w", retryErr)
			}
			log.Println("Message published successfully after topic creation")
			return nil
		}
		return fmt.Errorf("failed to publish message: %w", err)
	}
	log.Println("Message published successfully")
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

// createTopic attempts to create a Kafka topic with the specified
// number of partitions and replication factor.
// It connects to the first broker to send the create topic request.
func (c *kafkaClient) createTopic(
	topic string,
	partitions int,
	replicationFactor int,
) error {
	conn, err := c.Ping()
	if err != nil {
		return err
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafkago.Conn
	controllerConn, err = kafkago.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()
	controllerConn.SetDeadline(time.Now().Add(10 * time.Second))
	topicConfigs := kafkago.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}
	err = controllerConn.CreateTopics(topicConfigs)
	if err != nil {
		log.Printf("Failed to create topic %q: %v", topic, err)
		return err
	}
	log.Printf("Topic %q created successfully", topic)
	return nil
}

// isUnknownTopicError checks if the given error indicates
// that the topic or partition does not exist on the broker.
func isUnknownTopicError(err error) bool {
	return errors.Is(err, kafkago.UnknownTopicOrPartition)
}
