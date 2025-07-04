package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds the application configuration values loaded from environment variables.
type Config struct {
	// KafkaBroker contains one or more Kafka broker addresses.
	// Example: []string{"localhost:9092", "kafka2:9092"}
	Brokers []string

	// ProtoPath is the path to the .proto file used for message serialization.
	ProtoPath string

	// Port is the port where our Http server will be served from
	Port int32
}

// NewConfig loads environment variables into a Config struct.
// It optionally loads a .env file in local development environments,
// then uses Viper to retrieve typed configuration values.
func NewConfig() (*Config, error) {
	// Load .env in dev if present
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Bind env vars (this step is optional because AutomaticEnv picks them up)
	viper.BindEnv("KAFKA_BROKER")
	viper.BindEnv("PROTO_PATH")
	viper.BindEnv("PORT")

	// Viper does not split comma-separated values by default for GetStringSlice.
	// So we handle that ourselves:
	kafkaBrokers := viper.GetString("KAFKA_BROKER")
	var brokers []string
	if kafkaBrokers != "" {
		// Split by comma and trim spaces
		parts := strings.Split(kafkaBrokers, ",")
		for _, p := range parts {
			brokers = append(brokers, strings.TrimSpace(p))
		}
	}

	return &Config{
		Brokers:   brokers,
		ProtoPath: viper.GetString("PROTO_PATH"),
		Port:      viper.GetInt32("PORT"),
	}, nil
}
