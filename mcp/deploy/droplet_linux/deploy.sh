#!/bin/bash
# KDSE MCP Server - Deployment Script for Linux Droplet
# Usage: ./deploy.sh [--env-file <file>]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --env-file)
            ENV_FILE="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [--env-file <file>]"
            echo ""
            echo "Options:"
            echo "  --env-file <file>    Use custom environment file (default: .env)"
            echo "  --help               Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

echo "========================================"
echo "KDSE MCP Server - Linux Droplet Deploy"
echo "========================================"

# Check prerequisites
echo "[1/5] Checking prerequisites..."

if ! command -v docker &> /dev/null; then
    echo "ERROR: Docker is not installed"
    echo "Install Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! docker compose version &> /dev/null; then
    echo "ERROR: Docker Compose v2 is not installed"
    echo "Install Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "✓ Docker and Docker Compose are available"

# Load environment file if it exists
if [ -f "$ENV_FILE" ]; then
    echo "[2/5] Loading environment from $ENV_FILE..."
    set -a
    source "$ENV_FILE"
    set +a
else
    echo "[2/5] No .env file found, using defaults..."
    MCP_TRANSPORT="${MCP_TRANSPORT:-stdio}"
    MCP_HTTP_PORT="${MCP_HTTP_PORT:-8080}"
    KDSE_REPO_ROOT="${KDSE_REPO_ROOT:-/app}"
    IMAGE_TAG="${IMAGE_TAG:-latest}"
fi

# Build Docker image
echo "[3/5] Building Docker image..."

# Go to repository root
cd "$SCRIPT_DIR/../../../"

# Build the image
docker build -t "kdse-mcp:${IMAGE_TAG}" .

echo "✓ Docker image built successfully"

# Start services
echo "[4/5] Starting KDSE MCP Server..."

cd "$SCRIPT_DIR"

# Select profile based on transport
case "$MCP_TRANSPORT" in
    http)
        PROFILE="http"
        SERVICE="kdse-mcp-http"
        ;;
    *)
        PROFILE="stdio"
        SERVICE="kdse-mcp-stdio"
        ;;
esac

docker compose --profile "$PROFILE" up -d

# Wait for container to be healthy (if HTTP mode)
if [ "$MCP_TRANSPORT" = "http" ]; then
    echo "[5/5] Waiting for service to be healthy..."
    MAX_WAIT=30
    WAIT_COUNT=0
    while [ $WAIT_COUNT -lt $MAX_WAIT ]; do
        if curl -s http://localhost:${MCP_HTTP_PORT}/health > /dev/null 2>&1; then
            echo "✓ Service is healthy!"
            break
        fi
        sleep 1
        WAIT_COUNT=$((WAIT_COUNT + 1))
    done
    
    if [ $WAIT_COUNT -eq $MAX_WAIT ]; then
        echo "⚠ Service may not be fully healthy yet, check logs with:"
        echo "   docker compose logs -f $SERVICE"
    fi
else
    echo "[5/5] Service started in STDIO mode"
fi

echo ""
echo "========================================"
echo "Deployment Complete!"
echo "========================================"
echo ""
echo "Mode: $MCP_TRANSPORT"
if [ "$MCP_TRANSPORT" = "http" ]; then
    echo "URL: http://localhost:$MCP_HTTP_PORT"
    echo "Health: http://localhost:$MCP_HTTP_PORT/health"
fi
echo ""
echo "View logs: docker compose -f $SCRIPT_DIR/docker-compose.yml logs -f"
echo "Stop: docker compose -f $SCRIPT_DIR/docker-compose.yml down"
echo ""
