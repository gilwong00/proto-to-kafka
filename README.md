# Proto to Kafka Project

This repository contains two separate services demonstrating a Kafka event-driven system with Protobuf serialization:

- **Entity Producer** — a Go service that publishes Protobuf-encoded entity creation events to Kafka.
- **Entity Consumer** — a NestJS microservice that consumes these events, decodes the Protobuf messages, and processes entity workflows.

A shared `proto` directory contains the Protobuf `.proto` files managed by `buf` for schema consistency and code generation.

---

### Why Protobuf with Kafka?

Protobuf (Protocol Buffers) is used for message serialization to provide a compact, efficient, and strongly-typed binary format. This ensures that Kafka messages are transmitted with minimal overhead while maintaining strict schema validation. Schema validation helps catch errors early by enforcing message structure correctness both at production and consumption, improves compatibility between services, and supports backward/forward schema evolution — all essential for building reliable and maintainable event-driven systems.

## Architecture Overview

```
[Go Entity Producer]  -->  [Kafka Broker]  -->  [NestJS Entity Consumer]
    (protobuf)                (event bus)            (protobuf decode)
```

---

## Repository Structure

```
/
├── proto/                  # Shared Protobuf schema files and buf config
├── entity-producer/        # Go producer service
├── entity-consumer/        # NestJS consumer service
├── Makefile                # Common commands for Docker Compose and proto generation
├── docker-compose.yml      # Kafka, Kafka UI and Zookeeper services
└── README.md
```

---

## Prerequisites

- Go 1.20+ (for producer)
- Node.js 18+ and npm/yarn (for consumer)
- Docker (for Kafka, Zookeeper)
- buf CLI (v2) for Protobuf code generation
- Kafka cluster (local or remote)

---

## Setup

### 1. Generate Protobuf Code

The shared `.proto` files are in `proto/`. Use `buf` to generate Go and TypeScript code in the respective services:

```bash
# From project root
make generate
```

This runs `buf generate` in the `proto` directory and copies generated code to producer and consumer.

---

### 2. Start Kafka Locally

You can start Kafka and Zookeeper with Docker Compose via the Makefile:

```bash
make docker-up
```

Stop Kafka with:

```bash
make docker-down
```

---

## Entity Producer (Go)

### Setup & Run

```bash
cd entity-producer
go mod tidy
go run cmd/main.go
```

Configure via `.env` (example):

```env
KAFKA_BROKER=localhost:9092
PROTO_PATH=../proto/entity.proto
PORT=3000
```

The producer publishes Protobuf-encoded entity events to Kafka.

---

## Entity Consumer (NestJS)

### Setup & Run

```bash
cd entity-consumer
pnpm install
pnpm run start:dev
```

Configure `.env` (example):

```env
KAKFA_ENDPOINT=localhost:9092
```

The consumer listens to Kafka, decodes Protobuf messages, and processes entity events.

---

## Makefile Commands

| Command               | Description                                                  |
| --------------------- | ------------------------------------------------------------ |
| `make generate`       | Generate Protobuf code for producer and consumer using `buf` |
| `make docker-up`      | Start Kafka and Zookeeper via Docker Compose                 |
| `make docker-down`    | Stop Kafka and Zookeeper Docker containers                   |
| `make build-producer` | Builds the entity producer image                             |
| `make build-consumer` | Builds the entity consumer image                             |

---

## Notes

- Kafka topics must exist or be created dynamically.
- Protobuf ensures schema validation and compact messaging.
- Adjust `.env` files for your environment.
- `buf` manages Protobuf versioning and generation across services.
