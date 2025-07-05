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

.PHONY: build-producer
build-producer:
# ideally use a git hash to tag
	docker build -t entity-producer:$(shell date +%Y%m%d%H%M%S) -f Dockerfile.entity-producer .

.PHONY: build-consumer
build-consumer:
# ideally use a git hash to tag
	docker build -t entity-consumer:$(shell date +%Y%m%d%H%M%S) -f Dockerfile.entity-consumer .

.PHONY: build
build:
	@$(MAKE) build-producer &
	@$(MAKE) build-consumer &
	wait