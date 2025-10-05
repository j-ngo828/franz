# Multi-stage Dockerfile for Franz development and production

# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git protobuf-dev make

# Set working directory
WORKDIR /app

# Copy go mod files for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate protobuf code if needed
RUN if [ -d "api/proto" ]; then \
    protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/*.proto; \
    fi

# Build Franz components
RUN make build

# Development stage with hot reload using air
FROM golang:1.22-alpine AS dev

# Install development tools
RUN apk add --no-cache git protobuf-dev make && \
    go install github.com/cosmtrek/air@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Create non-root user for development
RUN addgroup -g 1001 -S franza && \
    adduser -u 1001 -S franz -G franza

# Create necessary directories
RUN mkdir -p /app/data /app/logs && \
    chown -R franz:franza /app

# Copy source code
WORKDIR /app
COPY . .

# Set permissions
RUN chown -R franz:franza /app

# Switch to non-root user
USER franz

# Expose ports
EXPOSE 9092 9093 9094

# Start air for hot reloading
CMD ["air", "-c", ".air.toml"]

# Production runtime stage
FROM alpine:latest AS runtime

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata protobuf-dev && \
    addgroup -g 1001 -S franza && \
    adduser -u 1001 -S franz -G franza

# Create necessary directories
RUN mkdir -p /app/data /app/logs /app/configs && \
    chown -R franz:franza /app

# Copy binaries from builder stage
COPY --from=builder /app/build/franz-broker* /app/
COPY --from=builder /app/build/franz-producer* /app/
COPY --from=builder /app/build/franz-consumer* /app/

# Copy configuration files
COPY configs/ /app/configs/

# Set permissions
RUN chown -R franz:franza /app && chmod +x /app/franz-*

# Switch to non-root user
USER franz

# Expose Franz ports
EXPOSE 9092 9093 9094

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD [ "/app/franz-broker", "--version" ] || exit 1

# Default command (can be overridden)
CMD ["/app/franz-broker"]

# Debug stage for troubleshooting
FROM runtime AS debug

# Install debugging tools
USER root
RUN apk add --no-cache busybox-extras net-tools lsof strace && \
    chown -R franz:franza /app

USER franz

# Keep container alive for debugging
CMD ["sh", "-c", "tail -f /dev/null"]
