package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	// TODO: update this to an array of strings to support multiple brookers
	KafkaBroker string
	ProtoPath   string
}

// NewConfig loads environment variables into a Config struct.
// It optionally loads a .env file in local development environments,
// then uses Viper to retrieve typed configuration values.
func NewConfig() (*Config, error) {
	// Load .env in dev if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("Error loading .env file: %v", err)
		}
	}

	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Bind specific environment variables to logical keys
	viper.BindEnv("KAFKA_BROKER")
	// might not need this
	viper.BindEnv("PROTO_PATH")

	return &Config{
		KafkaBroker: viper.GetString("KAFKA_BROKER"),
		ProtoPath:   viper.GetString("PROTO_PATH"),
	}, nil
}
