#!/bin/bash

# Script to setup directories and permissions on swarm nodes
# Run this on each node before deployment

set -e

GREEN='\033[0;32m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

# Create required directories
FRANZ_DATA_DIR="${FRANZ_DATA_DIR:-/var/lib/franz/data}"
FRANZ_LOG_DIR="${FRANZ_LOG_DIR:-/var/log/franz}"

log_info "Creating Franz directories..."

sudo mkdir -p "$FRANZ_DATA_DIR"
sudo mkdir -p "$FRANZ_LOG_DIR"

# Set ownership (adjust UID/GID as needed)
FRANZ_UID="${FRANZ_UID:-1001}"
FRANZ_GID="${FRANZ_GID:-1001}"

log_info "Setting ownership to $FRANZ_UID:$FRANZ_GID"

sudo chown -R "$FRANZ_UID:$FRANZ_GID" "$FRANZ_DATA_DIR"
sudo chown -R "$FRANZ_UID:$FRANZ_GID" "$FRANZ_LOG_DIR"

# Set permissions
sudo chmod 755 "$FRANZ_DATA_DIR"
sudo chmod 755 "$FRANZ_LOG_DIR"

log_info "Setup complete!"
log_info "Data directory: $FRANZ_DATA_DIR"
log_info "Log directory: $FRANZ_LOG_DIR"
