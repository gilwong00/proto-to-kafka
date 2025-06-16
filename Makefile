.PHONY: docker-compose
docker-compose:
	docker compose up -d

.PHONY: generate
generate:
	buf generate