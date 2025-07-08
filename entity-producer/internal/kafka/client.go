package kafka

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
)

const (
	maxRetryAttempts    = 3
	defaultRetryBackoff = 200 * time.Millisecond
)

// pendingMessage holds the details for a Kafka message to retry.
type pendingMessage struct {
	ctx       context.Context
	eventName string
	topic     string
	key       []byte
	value     []byte
}

// kafkaClient is a concrete implementation of the Client interface
// that publishes messages to Kafka using the provided configuration.
//
// It manages a shared kafka-go Writer for efficient reuse across multiple
// publish calls.
type kafkaClient struct {
	brokers    []string
	writer     *kafkago.Writer
	retryQueue chan pendingMessage
}

// Compile-time assertion to ensure kafkaClient implements the Client interface.
var _ Client = (*kafkaClient)(nil)

// newClient creates a new kafkaClient instance configured with the given Config.
// It lazily establishes connections and reuses a shared kafka-go Writer.
func newClient(config *Config) (*kafkaClient, error) {
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers configured")
	}
	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(config.Brokers...),
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
		BatchTimeout: 10 * time.Millisecond,
	}
	return &kafkaClient{
		brokers:    config.Brokers,
		writer:     writer,
		retryQueue: make(chan pendingMessage, 1000),
	}, nil
}

// Ping attempts to open a TCP connection to the first Kafka broker
// to verify that it is reachable and accepting connections.
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
// It checks for topic existence and creates it if necessary.
func (c *kafkaClient) Publish(
	ctx context.Context,
	eventName string,
	topic string,
	key []byte,
	value []byte,
) error {
	topics, err := c.ListTopics()
	if err != nil {
		return err
	}
	doesExist := false
	for _, topicName := range topics {
		if topic == topicName {
			doesExist = true
			break
		}
	}
	if !doesExist {
		if err := c.createTopic(topic, 1, 1); err != nil {
			return fmt.Errorf("failed to create topic: %w", err)
		}
		// Give Kafka time to propagate topic metadata
		time.Sleep(500 * time.Millisecond)
	}
	return c.queueMessage(ctx, eventName, topic, key, value, 0)
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

// ListTopics retrieves all topic names from Kafka.
func (c *kafkaClient) ListTopics() ([]string, error) {
	conn, err := kafkago.Dial("tcp", c.brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka broker: %w", err)
	}
	defer conn.Close()
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions: %w", err)
	}
	topicMap := make(map[string]struct{})
	for _, p := range partitions {
		topicMap[p.Topic] = struct{}{}
	}
	var topics []string
	for topic := range topicMap {
		topics = append(topics, topic)
	}
	return topics, nil
}

// createTopic attempts to create a Kafka topic with the specified
// number of partitions and replication factor.
// It connects to the first broker to send the create topic request.
func (c *kafkaClient) createTopic(
	topic string,
	partitions int,
	replicationFactor int,
) error {
	conn, err := kafkago.Dial("tcp", c.brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka broker: %w", err)
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get Kafka controller: %w", err)
	}
	controllerConn, err := kafkago.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return fmt.Errorf("failed to connect to controller broker: %w", err)
	}
	defer controllerConn.Close()
	controllerConn.SetDeadline(time.Now().Add(10 * time.Second))
	topicConfigs := kafkago.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}
	if err := controllerConn.CreateTopics(topicConfigs); err != nil {
		log.Printf("Failed to create topic %q: %v", topic, err)
		return err
	}
	log.Printf("Topic %q created successfully", topic)
	return nil
}

// queueMessage attempts to publish a message and retries on failure.
// If all attempts fail, the message is enqueued for background retry.
func (c *kafkaClient) queueMessage(
	ctx context.Context,
	eventName string,
	topic string,
	key []byte,
	value []byte,
	retryBackoff time.Duration,
) error {
	if retryBackoff == 0 {
		retryBackoff = defaultRetryBackoff
	}
	msg := kafkago.Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Time:  time.Now(),
		Headers: []kafkago.Header{
			{Key: "eventType", Value: []byte(eventName)},
		},
	}
	var lastErr error
	for attempt := 1; attempt <= maxRetryAttempts; attempt++ {
		err := c.writer.WriteMessages(ctx, msg)
		if err == nil {
			return nil
		}
		lastErr = err
		log.Printf("Publish attempt %d failed: %v", attempt, err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt) * retryBackoff):
		}
	}
	log.Printf("All %d publish attempts failed, enqueuing for retry", maxRetryAttempts)
	select {
	case c.retryQueue <- pendingMessage{ctx, eventName, topic, key, value}:
		log.Println("Message enqueued for retry")
	default:
		log.Println("Retry queue full, dropping message")
		return fmt.Errorf("retry queue full, dropping message: %w", lastErr)
	}
	return lastErr
}
