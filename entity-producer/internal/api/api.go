package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gilwong00/proto-to-kafka/internal/kafka"
)

type ApiService struct {
	kafkaClient kafka.Client
	port        int32
}

func NewApiService(kafkaClient kafka.Client, port int32) *ApiService {
	return &ApiService{
		kafkaClient: kafkaClient,
		port:        port,
	}
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, errMessage string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(errMessage))
}

func ParseJSONBody[T any](r *http.Request) (T, error) {
	var parsed T
	err := json.NewDecoder(r.Body).Decode(&parsed)
	if err != nil {
		return parsed, fmt.Errorf("unable to parse JSON: `%s`", err)
	}
	return parsed, nil
}

func WriteSuccessResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
