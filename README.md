# Franz - A High-Performance Message Queue in Go

Franz is a distributed message queue system inspired by Apache Kafka, built from scratch in Go. It's designed to achieve high throughput (1M+ messages/sec) with lower latency than traditional Java-based systems.

## Features

### Core Features (Phase 1)
- [x] Topic-based publish/subscribe messaging
- [x] Partitioned topics for parallelism
- [x] Consumer groups with automatic partition assignment
- [x] At-least-once delivery guarantees
- [x] Persistent storage with log segments
- [x] Leader/follower replication

### Advanced Features (Phase 2)
- [ ] Exactly-once semantics
- [ ] Transactions
- [ ] Stream processing SQL
- [ ] Schema registry
- [ ] Multi-tenancy
- [ ] Time-travel queries

## Architecture

Franz uses a distributed architecture with multiple brokers, where each topic is divided into partitions. Each partition has a leader broker and configurable replicas for fault tolerance.

## Quick Start

### Building Franz

```bash
# Install dependencies
make deps

# Build all components
make build

# Run tests
make test
```

### Running a Local Cluster

```bash
# Start a 3-broker cluster
make cluster-start

# Create a topic
./franz-admin create-topic --topic events --partitions 10 --replication-factor 3

# Produce messages
./franz-producer --broker localhost:9092 --topic events < messages.txt

# Consume messages
./franz-consumer --broker localhost:9092 --topic events --group my-app
```

## Project Structure

```
franz/
├── cmd/                 # Executable commands
│   ├── broker/         # Broker server
│   ├── producer/       # CLI producer
│   └── consumer/       # CLI consumer
├── pkg/                # Core packages
│   ├── broker/         # Broker implementation
│   ├── storage/        # Storage engine
│   ├── protocol/       # Wire protocol
│   ├── client/         # Client library
│   ├── consensus/      # Raft consensus
│   └── common/         # Shared utilities
├── api/proto/          # Protocol buffer definitions
├── configs/            # Configuration files
├── scripts/            # Utility scripts
└── tests/              # Test suites
```

## Development

### Prerequisites

- Go 1.22 or later
- Protocol Buffers compiler
- Make

### Building from Source

```bash
git clone https://github.com/jang3435/franz
cd franz
make build
```

### Running Tests

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Benchmarks
make bench
```

## Performance

Franz is designed for high throughput and low latency:

- **Target**: 1M+ messages/second per broker
- **Latency**: Sub-millisecond p99
- **Storage**: Efficient segment-based storage with compression
- **Network**: Zero-copy transfers, batching, and pipelining

## Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests to our repository.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
