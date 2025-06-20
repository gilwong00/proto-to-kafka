package main

import (
	"fmt"

	"github.com/gilwong00/proto-to-kafka/internal/config"
	"github.com/gilwong00/proto-to-kafka/internal/kafka"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		// log error
		panic(err)
	}
	// initialize kafka client
	kafkaClient := kafka.NewClient(config)
	fmt.Printf(">>> %+v\n", kafkaClient)
}
