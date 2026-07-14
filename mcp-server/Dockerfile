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

# Install git for repository access
RUN apk add --no-cache git ca-certificates

# Create non-root user
RUN adduser -D -u 1000 kdse

WORKDIR /home/kdse

# Copy binary from builder
COPY --from=builder /build/kdse-mcp-server .

# Copy repository root for read access
COPY . ./

# Set ownership
RUN chown -R kdse:kdse /home/kdse

# Environment variables
ENV KDSE_REPO_ROOT=/home/kdse
ENV MCP_STDIO=true

# User setup
USER kdse

# Default command - MCP server over stdio
CMD ["./kdse-mcp-server"]

# Alternative: MCP server over HTTP (uncomment if needed)
# EXPOSE 8080
# CMD ["./kdse-mcp-server", "--http"]
