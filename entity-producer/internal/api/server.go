package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gilwong00/proto-to-kafka/gen"
	"github.com/gilwong00/proto-to-kafka/internal/kafka"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func (s ApiService) StartHttpServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /entity", s.createEntity)

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", s.port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// start the server
	go func() {
		fmt.Printf("Starting server on port %v\n", s.port)
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Error starting server: %s", err.Error())
			os.Exit(1)
		}
	}()
	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return server.Shutdown(ctx)
}

func (s *ApiService) createEntity(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error closing request body: %v", err)
		}
		if err := s.kafkaClient.Close(); err != nil {
			log.Printf("Error closing Kafka client: %v", err)
		}
	}()

	payload, err := ParseJSONBody[CreateEntityPayload](r)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	entityId, err := uuid.NewUUID()
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	entityPb := &gen.Entity{
		Id:   entityId.String(),
		Name: payload.name,
	}
	msgBytes, err := proto.Marshal(entityPb)
	if err != nil {
		log.Fatalf("Failed to marshal: %v", err)
	}
	kafkaKey := s.kafkaClient.GenerateKey("entity")
	if err := s.kafkaClient.Publish(
		ctx,
		"entity-created",
		kafka.EntityTopic,
		kafkaKey,
		msgBytes,
	); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	WriteSuccessResponse(w, "entity created")
}
