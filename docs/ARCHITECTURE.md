# Franz Architecture Documentation

## Overview

Franz is a high-performance distributed message queue system inspired by Apache Kafka, built from scratch in Go. This document describes the architecture and implementation strategy.

## Project Structure

```
franz/
├── cmd/                    # Executable commands
│   ├── broker/            # Broker server
│   ├── producer/          # CLI producer tool
│   └── consumer/          # CLI consumer tool
├── pkg/                   # Core packages
│   ├── broker/            # Broker implementation
│   ├── storage/           # Storage engine
│   ├── protocol/          # Wire protocol
│   ├── client/            # Client library
│   ├── consensus/         # Consensus (Raft/etcd)
│   └── common/            # Shared types and utilities
├── api/proto/             # Protocol buffer definitions
├── configs/               # Configuration files
├── scripts/               # Utility scripts
├── docs/                  # Documentation
└── tests/                 # Test suites
```

## Core Components

### 1. Broker (`pkg/broker`)
- **Purpose**: Main server component that handles client requests
- **Key Interfaces**:
  - `Broker`: Main broker service
  - `PartitionManager`: Manages topic partitions
  - `ReplicationManager`: Handles data replication
  - `ClusterManager`: Manages cluster coordination

### 2. Storage Engine (`pkg/storage`)
- **Purpose**: Persistent storage for messages
- **Key Interfaces**:
  - `Storage`: Storage engine interface
  - `Log`: Partition log management
  - `Segment`: Log segment files
  - `Index`: Offset indexing
- **Design**: Segment-based storage with memory-mapped files for performance

### 3. Protocol (`pkg/protocol`)
- **Purpose**: Wire protocol for client-broker communication
- **Key Components**:
  - Binary protocol for efficiency
  - gRPC for modern API
  - Request/Response types
  - Error codes

### 4. Client Library (`pkg/client`)
- **Purpose**: Go client for producing and consuming messages
- **Key Interfaces**:
  - `Client`: Main client interface
  - `Producer`: Message producer
  - `Consumer`: Message consumer
  - `Partitioner`: Partitioning strategy

### 5. Consensus (`pkg/consensus`)
- **Purpose**: Distributed consensus for metadata and coordination
- **Options**:
  - Raft implementation (built-in)
  - etcd integration (external)
- **Used For**:
  - Leader election
  - Metadata storage
  - Configuration management

### 6. Common (`pkg/common`)
- **Purpose**: Shared types and utilities
- **Contains**:
  - `Message` struct
  - Error definitions
  - Compression types

## Implementation Phases

### Phase 1: Core Functionality (MVP)
1. **Basic Storage Engine**
   - Implement log segments
   - Basic append/read operations
   - Simple indexing

2. **Simple Broker**
   - TCP server
   - Handle produce/fetch requests
   - Single partition support

3. **Basic Client**
   - Producer implementation
   - Consumer implementation
   - Simple partitioner

4. **Minimal Protocol**
   - Produce request/response
   - Fetch request/response
   - Metadata request/response

### Phase 2: Distribution & Reliability
1. **Partitioning**
   - Multiple partitions per topic
   - Partition assignment
   - Load balancing

2. **Replication**
   - Leader/follower model
   - In-sync replicas (ISR)
   - Leader election

3. **Consumer Groups**
   - Group coordination
   - Partition assignment
   - Offset management

4. **Consensus Integration**
   - Raft implementation
   - Metadata management
   - Cluster coordination

### Phase 3: Performance & Features
1. **Performance Optimizations**
   - Zero-copy transfers
   - Batch processing
   - Compression support
   - Memory-mapped files

2. **Advanced Features**
   - Transactions
   - Exactly-once semantics
   - Log compaction
   - Quotas

### Phase 4: Enhancements Beyond Kafka
1. **Native Stream Processing**
   - Built-in SQL support
   - Stream joins and aggregations

2. **Time-Travel Queries**
   - Query historical data by timestamp
   - Efficient time-based indexing

3. **Multi-Tenancy**
   - Namespace isolation
   - Resource quotas per tenant

4. **Adaptive Partitioning**
   - Dynamic partition rebalancing
   - Load-based partition splitting

## Development Guidelines

### Code Organization
- Each package should have clear interfaces
- Implementation details should be private
- Use dependency injection for testability
- Write comprehensive tests for each component

### Performance Considerations
- Minimize allocations in hot paths
- Use object pools where appropriate
- Leverage Go's concurrency primitives
- Profile and benchmark regularly

### Testing Strategy
- Unit tests for all packages
- Integration tests for end-to-end flows
- Benchmark tests for performance
- Chaos testing for reliability

## Next Steps

1. **Start with Storage Engine**
   - Implement basic log structure
   - Add segment management
   - Create simple index

2. **Build Basic Broker**
   - TCP server setup
   - Request handling
   - Storage integration

3. **Create Simple Client**
   - Producer implementation
   - Basic consumer
   - Test end-to-end flow

4. **Add Tests**
   - Unit tests for each component
   - Integration tests
   - Benchmarks

## Building and Running

```bash
# Download dependencies
make deps

# Build all components
make build

# Run tests
make test

# Start a broker
./build/franz-broker --config configs/broker.yaml

# Produce messages
echo "Hello Franz!" | ./build/franz-producer --topic test

# Consume messages
./build/franz-consumer --topic test --group my-app
```

## Contributing

When implementing new features:
1. Start with the interface definition
2. Write tests first (TDD approach)
3. Implement the functionality
4. Ensure all tests pass
5. Add benchmarks if applicable
6. Update documentation
