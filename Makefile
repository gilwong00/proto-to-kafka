.PHONY: docker-compose
docker-compose:
	docker compose up -d

.PHONY: generate
generate:
	buf generate

.PHONY: start-producer
start-producer:
	cd entity-producer && go run cmd/main.go