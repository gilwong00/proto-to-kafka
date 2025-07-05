Got it! Here’s a clean README draft for your Go entity producer with buf v2, without a contributions section:

---

# Entity Producer (Go)

This is a demo Go project that produces Protobuf-encoded messages to Kafka.

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/gilwong00/entity-producer.git
cd entity-producer
```

### 2. Install dependencies

Run the following to install all Go dependencies:

```bash
go mod tidy
```

### 3. Set up environment variables

Create a `.env` file at the project root for local development:

```env
KAFKA_BROKER=localhost:9092
PROTO_PATH=./proto/entity.proto
```

To specify multiple Kafka brokers, separate them by commas:

```env
KAFKA_BROKER=localhost:9092,localhost:9093
```

### 4. Generate protobuf code using buf v2

This project uses **[buf v2](https://buf.build/)** to generate Go protobuf code.

Install `buf` if you haven’t already:

```bash
brew install bufbuild/buf/buf  # macOS (Homebrew)
# Or follow official installation steps: https://docs.buf.build/installation
```

Generate Go code from `.proto` definitions:

```bash
buf generate
```

This command uses the `buf.yaml` and `buf.gen.yaml` files included in the repo.

### 5. Build and run

```bash
go build -o entity-producer ./cmd/producer
./entity-producer
```

Messages will be published to Kafka brokers specified in the `.env` file.
