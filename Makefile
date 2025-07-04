.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down -v --remove-orphans

.PHONY: generate
generate:
	buf generate

.PHONY: start-producer
start-producer:
	cd entity-producer && go run cmd/main.go