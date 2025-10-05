# Deployment Configurations

This directory contains deployment configurations for different environments.

## Docker Swarm Deployment

### Prerequisites

- Docker Swarm initialized on your cluster
- At least 3 nodes for Raft consensus (recommended)
- Proper network connectivity between nodes

### Deployment Steps

1. **Setup nodes** - Run on each node:
   ```bash
   ./scripts/setup-nodes.sh
   ```

2. **Configure environment** - Copy and edit environment file:
   ```bash
   cp .env.production.example .env.production
   ```

3. **Deploy the stack**:
   ```bash
   ./scripts/deploy.sh
   ```

### Management Commands

Check status:
```bash
./scripts/deploy.sh status
```

Remove stack:
```bash
./scripts/deploy.sh remove
```

Scale brokers:
```bash
docker service scale franz_franz-broker=5
```

View logs:
```bash
docker service logs -f franz_franz-broker
```

## Kubernetes Deployment

TODO: Add Kubernetes manifests and Helm charts
