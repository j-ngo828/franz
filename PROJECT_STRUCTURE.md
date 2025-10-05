# Franz Project Structure

```
franz/
├── cmd/                          # Command-line applications
│   ├── broker/                   # Main broker server
│   │   └── main.go              # Broker entry point
│   ├── producer/                 # Producer CLI tool
│   │   └── main.go              # Producer CLI entry point
│   ├── consumer/                 # Consumer CLI tool
│   │   └── main.go              # Consumer CLI entry point
│   └── cli/                      # Admin/management CLI
│       └── main.go              # Admin CLI entry point
│
├── internal/                     # Private application code
│   ├── broker/                   # Broker core logic
│   │   └── broker.go            # Broker service implementation
│   ├── storage/                  # Log-structured storage
│   │   └── storage.go           # Storage interface and implementation
│   ├── raft/                     # Raft consensus protocol
│   │   └── raft.go              # Raft node implementation
│   ├── protocol/                 # Wire protocol
│   │   └── protocol.go          # Binary protocol encoding/decoding
│   ├── config/                   # Configuration management
│   │   └── config.go            # Config structures and loading
│   └── metrics/                  # Metrics and monitoring
│       └── metrics.go           # Prometheus metrics
│
├── pkg/                          # Public library code
│   ├── client/                   # Client SDK
│   │   └── client.go            # Producer/Consumer client library
│   └── types/                    # Shared types
│       └── types.go             # Common data structures
│
├── tests/                        # Test suites
│   ├── unit/                     # Unit tests
│   ├── integration/              # Integration tests
│   └── benchmark/                # Benchmark tests
│
├── configs/                      # Configuration files
│   └── broker.example.yml       # Example broker configuration
│
├── scripts/                      # Deployment and utility scripts
│   ├── deploy.sh                # Docker Swarm deployment script
│   ├── setup-nodes.sh           # Node setup script
│   └── tests/                   # Test scripts
│       └── run-integration.sh   # Integration test runner
│
├── deployments/                  # Deployment configurations
│   └── README.md                # Deployment documentation
│
├── docs/                         # Documentation
│   ├── architecture.md          # Architecture overview
│   └── data-flow.mmd            # Data flow diagram
│
├── stack.yml                     # Docker Swarm stack file (production)
├── docker-compose.dev.yml        # Docker Compose for development
├── Dockerfile                    # Multi-stage Dockerfile
├── Makefile                      # Build automation
├── .air.toml                     # Hot reload configuration
├── .env.production.example       # Production environment variables
└── README.md                     # Project README
```

## Directory Guidelines

### `/cmd`
- Contains main applications for this project
- Each subdirectory has a `main.go` file
- Keep minimal code here, defer to `/internal` or `/pkg`
- Binary names should match directory names

### `/internal`
- Private application and library code
- Cannot be imported by external projects
- Core business logic goes here
- Structure mirrors the application architecture

### `/pkg`
- Library code that external applications can import
- Client SDK and shared types
- Must have a clear API and documentation

### `/tests`
- Additional test files and test data
- Organized by test type (unit, integration, benchmark)
- Can have helper packages for test utilities

### `/configs`
- Configuration file templates and examples
- Default configurations for different environments

### `/scripts`
- Scripts for deployment, automation, and development
- Should be documented and executable

### `/deployments`
- Infrastructure as Code (IaaS, PaaS, container orchestration)
- Kubernetes, Docker Swarm, cloud deployment configs

### `/docs`
- Design documents, architecture diagrams, API docs
- User and developer documentation

## Build Targets

- **Development**: Use `docker-compose.dev.yml` with hot reload
- **Production**: Build and deploy with `stack.yml` to Docker Swarm

## Naming Conventions

- **Go files**: `lowercase_with_underscores.go`
- **Go packages**: `lowercase` (no underscores)
- **Directories**: `lowercase` (kebab-case for multi-word)
- **Binaries**: `lowercase` (same as cmd directory)

## Import Paths

Assuming module name is `github.com/yourusername/franz`:

```go
import (
    "github.com/yourusername/franz/internal/broker"
    "github.com/yourusername/franz/pkg/client"
    "github.com/yourusername/franz/pkg/types"
)
```
