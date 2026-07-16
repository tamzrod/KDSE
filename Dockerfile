# KDSE MCP Server Dockerfile
# Multi-stage build for minimal image size
# Deployment: KDSE Runtime subsystem (github.com/kdse/runtime)
# Contains: CLI (cmd/kdse/) and MCP Server (cmd/mcp/) binaries
# Version: 2.0.0

# Build stage
FROM golang:1.22-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build both binaries
# Runtime CLI
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kdse ./cmd/kdse/
# MCP Server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kdse-mcp ./cmd/mcp/

# Runtime stage - minimal image
FROM alpine:3.19

# Install git for repository access (if needed at runtime)
# Install wget for healthcheck
RUN apk add --no-cache git ca-certificates wget

# Create non-root user
RUN adduser -D -u 1000 kdse

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/kdse .
COPY --from=builder /build/kdse-mcp .

# Set ownership
RUN chown -R kdse:kdse /app

# Environment variables (defaults)
ENV MCP_TRANSPORT=stdio
ENV MCP_HTTP_PORT=8080

# User setup
USER kdse

# Healthcheck for HTTP mode
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/health || exit 1

# Default command - MCP server
# Transport is selected via MCP_TRANSPORT environment variable
CMD ["./kdse-mcp"]
