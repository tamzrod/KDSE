# KDSE MCP Server Dockerfile
# Multi-stage build for minimal image size
# Version: 0.2.0

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

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kdse-mcp-server .

# Runtime stage - minimal image
FROM alpine:3.19

# Install git for repository access (if needed at runtime)
# Install wget for healthcheck
RUN apk add --no-cache git ca-certificates wget

# Create non-root user
RUN adduser -D -u 1000 kdse

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/kdse-mcp-server .

# Set ownership
RUN chown -R kdse:kdse /app

# Environment variables (defaults)
ENV MCP_TRANSPORT=stdio
ENV MCP_HTTP_PORT=8080

# User setup
USER kdse

# Default command - MCP server
# Transport is selected via MCP_TRANSPORT environment variable
CMD ["./kdse-mcp-server"]
