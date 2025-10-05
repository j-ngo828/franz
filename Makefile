.PHONY: all build clean test test-integration bench deps proto fmt lint \
        docker-build docker-run docker-dev docker-test docker-push \
        dev deploy deploy-status deploy-remove setup tools install-air

# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOLINT=golangci-lint
DOCKER=docker
DOCKER_COMPOSE=docker-compose
BUILD_DIR=build
BINARY_PREFIX=franz
BROKER_BINARY=$(BUILD_DIR)/$(BINARY_PREFIX)-broker
PRODUCER_BINARY=$(BUILD_DIR)/$(BINARY_PREFIX)-producer
CONSUMER_BINARY=$(BUILD_DIR)/$(BINARY_PREFIX)-consumer

# Go parameters
GO111MODULE=on
CGO_ENABLED=0
GOFLAGS=-trimpath
LDFLAGS=-s -w -X main.version=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev") \
        -X main.commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
        -X main.buildTime=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# Build targets
all: fmt lint test build

# === BUILD TARGETS ===

# Install air for hot reloading
install-air:
	$(GOCMD) install github.com/cosmtrek/air@latest

# Setup development environment
setup: deps install-air
	@if [ ! -f ".air.toml" ]; then \
		curl -o .air.toml https://raw.githubusercontent.com/cosmtrek/air/master/air.toml; \
	fi

tools: install-air
	$(GOCMD) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOCMD) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install github.com/cosmtrek/air@latest

deps:
	$(GOMOD) download
	$(GOMOD) tidy

build: $(BUILD_DIR) build-broker build-producer build-consumer

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

build-broker:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GOBUILD) \
		-ldflags "$(LDFLAGS)" \
		-o $(BROKER_BINARY) \
		-v ./cmd/broker

build-producer:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GOBUILD) \
		-ldflags "$(LDFLAGS)" \
		-o $(PRODUCER_BINARY) \
		-v ./cmd/producer

build-consumer:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GOBUILD) \
		-ldflags "$(LDFLAGS)" \
		-o $(CONSUMER_BINARY) \
		-v ./cmd/consumer

# === TESTING TARGETS ===

test:
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-integration:
	@if [ -d "tests/integration" ]; then \
		$(GOTEST) -v -tags=integration ./tests/integration/...; \
	else \
		@echo "Integration tests not found. Create tests/integration/ directory."; \
	fi

bench:
	$(GOTEST) -bench=. -benchmem ./...

# === CODE QUALITY TARGETS ===

fmt:
	$(GOFMT) -s -w ./...
	@echo "Go code formatted"

lint:
	$(GOLINT) run ./...
	@echo "Linting completed"

# === DOCKER TARGETS ===

docker-build:
	$(DOCKER) build --target runtime -t franz:$(shell git rev-parse --short HEAD) .
	$(DOCKER) tag franz:$(shell git rev-parse --short HEAD) franz:latest

docker-build-dev:
	$(DOCKER) build --target dev -t franz:dev .
	$(DOCKER) build --target debug -t franz:debug .

docker-run:
	$(DOCKER) run --rm -p 9092:9092 -v $(shell pwd)/data:/app/data franz:latest

docker-run-dev:
	$(DOCKER) run --rm \
		-p 9092:9092 \
		-p 9093:9093 \
		-v $(shell pwd):/app \
		-v $(shell pwd)/data:/app/data \
		franz:dev

docker-dev: docker-build-dev
	@if [ ! -f ".air.toml" ]; then \
		$(MAKE) install-air; \
		curl -o .air.toml https://raw.githubusercontent.com/cosmtrek/air/master/.air.toml; \
	fi
	$(DOCKER) run --rm \
		-p 9092:9092 \
		-p 9093:9093 \
		-v $(shell pwd):/app \
		-v $(shell pwd)/data:/app/data \
		franz:dev

docker-test:
	$(DOCKER) run --rm \
		-v $(shell pwd):/app \
		-w /app \
		golang:1.22-alpine sh -c "make test"

docker-push: docker-build
	$(DOCKER) push franz:latest

# === DOCKER COMPOSE TARGETS (Development Only) ===

dev:
	$(DOCKER_COMPOSE) -f docker-compose.dev.yml up --build

dev-down:
	$(DOCKER_COMPOSE) -f docker-compose.dev.yml down -v

dev-logs:
	$(DOCKER_COMPOSE) -f docker-compose.dev.yml logs -f

# === DOCKER SWARM DEPLOYMENT (Production) ===

deploy:
	@bash scripts/deploy.sh deploy

deploy-status:
	@bash scripts/deploy.sh status

deploy-remove:
	@bash scripts/deploy.sh remove

setup-nodes:
	@bash scripts/setup-nodes.sh

# === DEVELOPMENT WORKFLOW ===

run-broker:
	@if [ -f "$(BROKER_BINARY)" ]; then \
		./$(BROKER_BINARY) --config configs/broker.yaml; \
	else \
		@echo "Broker binary not found. Run 'make build-broker' first."; \
		exit 1; \
	fi

run-producer:
	@if [ -f "$(PRODUCER_BINARY)" ]; then \
		./$(PRODUCER_BINARY) --help; \
	else \
		@echo "Producer binary not found. Run 'make build-producer' first."; \
		exit 1; \
	fi

run-consumer:
	@if [ -f "$(CONSUMER_BINARY)" ]; then \
		./$(CONSUMER_BINARY) --help; \
	else \
		@echo "Consumer binary not found. Run 'make build-consumer' first."; \
		exit 1; \
	fi

# === PROTOCOL BUFFERS ===

proto:
	@if [ -d "api/proto" ]; then \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			api/proto/*.proto; \
		echo "Protocol buffers generated"; \
	else \
		@echo "No proto directory found"; \
	fi

# === CLEANUP ===

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.txt

clean-docker:
	$(DOCKER) system prune -f
	$(DOCKER) image rm franz:latest franz:dev franz:debug 2>/dev/null || true

clean-dist:
	rm -rf dist/
	rm -rf build/

clean-all: clean clean-docker clean-dist
	@echo "All artifacts cleaned"

# === UTILITY TARGETS ===

version:
	@echo "Franz version: $(shell git describe --tags --abbrev=0 2>/dev/null || echo 'dev')"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build time: $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"

help:
	@echo "Available targets:"
	@echo ""
	@echo "Build & Setup:"
	@echo "  setup            - Setup development environment with all dependencies"
	@echo "  build            - Build all Franz components"
	@echo "  deps             - Download and tidy Go dependencies"
	@echo ""
	@echo "Testing:"
	@echo "  test             - Run unit tests with coverage"
	@echo "  test-integration - Run integration tests"
	@echo "  bench            - Run performance benchmarks"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt              - Format Go code"
	@echo "  lint             - Run GolangCI-Lint"
	@echo ""
	@echo "Development (Local with Hot Reload):"
	@echo "  dev              - Start development environment with docker-compose"
	@echo "  dev-down         - Stop development environment"
	@echo "  dev-logs         - View development logs"
	@echo ""
	@echo "Production Deployment (Docker Swarm):"
	@echo "  deploy           - Deploy Franz to Docker Swarm"
	@echo "  deploy-status    - Check deployment status"
	@echo "  deploy-remove    - Remove deployment from Swarm"
	@echo "  setup-nodes      - Setup directories on Swarm nodes"
	@echo ""
	@echo "Docker Images:"
	@echo "  docker-build     - Build production Docker image"
	@echo "  docker-dev       - Run development container with hot reload"
	@echo "  docker-push      - Push image to registry"
	@echo ""
	@echo "Utilities:"
	@echo "  proto            - Generate protocol buffer code"
	@echo "  clean            - Remove build artifacts"
	@echo "  version          - Show version information"
	@echo "  help             - Show this help"

# Default target
.DEFAULT_GOAL := help

# Include local makefile if it exists
include Makefile.local 2>/dev/null || true
