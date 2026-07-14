# KDSE MCP Server Dockerfile
# Multi-stage build for minimal image size

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
RUN apk add --no-cache git ca-certificates

# Create non-root user
RUN adduser -D -u 1000 kdse

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/kdse-mcp-server .

# Set ownership
RUN chown -R kdse:kdse /app

# Environment variables
ENV MCP_STDIO=true

# User setup
USER kdse

# Default command - MCP server over stdio
CMD ["./kdse-mcp-server"]
