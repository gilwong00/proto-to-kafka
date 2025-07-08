package main

import (
	"log"

	"github.com/gilwong00/proto-to-kafka/internal/api"
	"github.com/gilwong00/proto-to-kafka/internal/config"
	"github.com/gilwong00/proto-to-kafka/internal/kafka"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Printf("failed to initialized config %v", err)
		panic(err)
	}
	// initialize kafka client
	kafkaClient, err := kafka.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating kafka client %v", err)
	}
	// test kafka connection
	conn, err := kafkaClient.Ping()
	if err != nil {
		log.Fatalf("Error connecting to kafka %v", err)
	}
	defer conn.Close()
	// initialize api server
	apiServer := api.NewApiService(kafkaClient, config.Port)
	if err := apiServer.StartHttpServer(); err != nil {
		log.Fatalf("Error while starting server: %+v", err)
	}
}
