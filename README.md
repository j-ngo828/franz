# Franz - A High-Performance Message Queue in Go

Franz is a distributed message queue system inspired by Apache Kafka, built from scratch in Go. It's designed to achieve high throughput (1M+ messages/sec) with lower latency than traditional Java-based systems.

## Features

### Core Features (Phase 1)
- [ ] Topic-based publish/subscribe messaging
- [ ] Partitioned topics for parallelism
- [ ] Consumer groups with automatic partition assignment
- [ ] At-least-once delivery guarantees
- [ ] Persistent storage with log segments
- [ ] Leader/follower replication

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
├── cmd/                    # Command-line applications
│   ├── broker/            # Main broker server
│   ├── producer/          # Producer CLI tool
│   ├── consumer/          # Consumer CLI tool
│   └── cli/               # Admin/management CLI
├── internal/              # Private application code
│   ├── broker/            # Broker core logic
│   ├── storage/           # Log-structured storage
│   ├── raft/              # Raft consensus protocol
│   ├── protocol/          # Wire protocol
│   ├── config/            # Configuration management
│   └── metrics/           # Metrics and monitoring
├── pkg/                   # Public library code
│   ├── client/            # Client SDK
│   └── types/             # Shared types
├── configs/               # Configuration files
├── scripts/               # Deployment and utility scripts
│   ├── deploy.sh          # Docker Swarm deployment
│   └── setup-nodes.sh     # Node preparation
├── deployments/           # Deployment configurations
├── docs/                  # Documentation
└── tests/                 # Test suites
    ├── unit/
    ├── integration/
    └── benchmark/
```

See [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) for detailed documentation.

## Development

### Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- Protocol Buffers compiler (optional)
- Make

### Development with Docker Compose

```bash
# Start development environment with hot reload
make dev

# View logs
make dev-logs

# Stop development environment
make dev-down
```

### Building from Source

```bash
git clone https://github.com/jang3435/franz
cd franz
make build
```

## Deployment

### Docker Swarm (Production)

Franz uses Docker Swarm for production deployments with Raft consensus across multiple nodes.

#### Prerequisites

- Docker Swarm initialized cluster
- At least 3 nodes (recommended for Raft consensus)
- Network connectivity between nodes

#### Deployment Steps

1. **Prepare environment**:
   ```bash
   cp .env.production.example .env.production
   # Edit .env.production with your configuration
   ```

2. **Setup nodes** (run on each node):
   ```bash
   make setup-nodes
   ```

3. **Deploy to Swarm**:
   ```bash
   make deploy
   ```

4. **Check deployment status**:
   ```bash
   make deploy-status
   ```

#### Management Commands

```bash
# View service logs
docker service logs -f franz_franz-broker

# Scale brokers
docker service scale franz_franz-broker=5

# Update service
VERSION=v1.2.0 make deploy

# Remove deployment
make deploy-remove
```

### Docker Compose (Development Only)

For local development with hot reload:

```bash
# Start single broker with hot reload
make dev

# View logs
make dev-logs

# Stop
make dev-down
```

### Running Tests

```bash
# Unit tests
make test

# Integration tests
make test-integration

# Benchmarks
make bench

# Run tests in Docker
make docker-test
```

### Port Configuration

Development (single broker):
- **Franz Protocol**: `localhost:9092`
- **Raft Consensus**: `localhost:9093`
- **Metrics**: `localhost:9094`

Production (Docker Swarm):
- Services use overlay networking
- Published ports: 9092 (Franz), 9093 (Raft), 9094 (Metrics)

## Troubleshooting

### Development Issues

**Container won't start:**
```bash
# Check logs
make dev-logs

# Rebuild without cache
docker-compose -f docker-compose.dev.yml build --no-cache
```

**Port conflicts:**
```bash
# Check what's using ports
lsof -i :9092

# Modify docker-compose.dev.yml to use different ports
```

### Deployment Issues

**Swarm not initialized:**
```bash
docker swarm init
```

**Service not starting:**
```bash
# Check service status
docker service ps franz_franz-broker --no-trunc

# View service logs
docker service logs franz_franz-broker
```

**Node labels missing:**
```bash
# Manually label nodes
docker node update --label-add franz.broker=true <node-id>
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
