#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
STACK_NAME="${STACK_NAME:-franz}"
VERSION="${VERSION:-latest}"
REGISTRY="${REGISTRY:-}"
ENV_FILE="${ENV_FILE:-.env.production}"

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running in swarm mode
check_swarm() {
    if ! docker info --format '{{.Swarm.LocalNodeState}}' | grep -q "active"; then
        log_error "Docker is not running in swarm mode"
        log_info "Initialize swarm with: docker swarm init"
        exit 1
    fi
    log_info "Swarm mode active"
}

# Label nodes for placement
label_nodes() {
    log_info "Labeling nodes for Franz broker placement..."

    # Get all worker nodes
    NODES=$(docker node ls --filter role=worker -q)

    if [ -z "$NODES" ]; then
        log_warn "No worker nodes found, using manager nodes"
        NODES=$(docker node ls -q | head -3)
    fi

    COUNT=0
    for NODE in $NODES; do
        if [ $COUNT -lt 3 ]; then
            docker node update --label-add franz.broker=true "$NODE"
            log_info "Labeled node $(docker node inspect --format '{{.Description.Hostname}}' $NODE)"
            COUNT=$((COUNT + 1))
        fi
    done

    if [ $COUNT -lt 3 ]; then
        log_warn "Only $COUNT nodes labeled (minimum 3 recommended for raft consensus)"
    fi
}

# Build and tag image
build_image() {
    log_info "Building Franz image version: $VERSION"

    docker build \
        --target runtime \
        --tag franz:${VERSION} \
        --tag franz:latest \
        .

    if [ -n "$REGISTRY" ]; then
        log_info "Tagging for registry: $REGISTRY"
        docker tag franz:${VERSION} ${REGISTRY}/franz:${VERSION}
        docker tag franz:latest ${REGISTRY}/franz:latest
    fi
}

# Push to registry
push_image() {
    if [ -n "$REGISTRY" ]; then
        log_info "Pushing images to registry..."
        docker push ${REGISTRY}/franz:${VERSION}
        docker push ${REGISTRY}/franz:latest
    else
        log_warn "No registry specified, skipping push"
    fi
}

# Create required directories on nodes
setup_directories() {
    log_info "Setting up data directories..."

    # TODO: Use docker exec or ssh to create directories on each node
    # For now, assume directories exist or will be created by volumes

    log_info "Ensure the following directories exist on all nodes:"
    log_info "  - /var/lib/franz/data"
    log_info "  - /var/log/franz"
}

# Deploy stack
deploy_stack() {
    log_info "Deploying stack: $STACK_NAME"

    # Export variables for docker stack
    export VERSION

    # Load environment file if exists
    if [ -f "$ENV_FILE" ]; then
        log_info "Loading environment from: $ENV_FILE"
        set -a
        source "$ENV_FILE"
        set +a
    fi

    docker stack deploy \
        --compose-file stack.yml \
        --with-registry-auth \
        "$STACK_NAME"

    log_info "Stack deployed successfully"
}

# Show stack status
show_status() {
    log_info "Stack services:"
    docker stack services "$STACK_NAME"

    echo ""
    log_info "Running tasks:"
    docker stack ps "$STACK_NAME" --no-trunc
}

# Main deployment flow
main() {
    log_info "Starting Franz deployment..."

    check_swarm
    label_nodes
    build_image
    push_image
    setup_directories
    deploy_stack

    echo ""
    log_info "Deployment complete!"
    echo ""

    show_status

    echo ""
    log_info "To check logs: docker service logs -f ${STACK_NAME}_franz-broker"
    log_info "To scale: docker service scale ${STACK_NAME}_franz-broker=5"
    log_info "To remove: docker stack rm ${STACK_NAME}"
}

# Parse command line arguments
case "${1:-deploy}" in
    deploy)
        main
        ;;
    status)
        show_status
        ;;
    remove)
        log_info "Removing stack: $STACK_NAME"
        docker stack rm "$STACK_NAME"
        ;;
    *)
        echo "Usage: $0 {deploy|status|remove}"
        exit 1
        ;;
esac
