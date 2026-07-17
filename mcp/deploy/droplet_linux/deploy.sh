#!/bin/bash
#===============================================================================
# KDSE MCP Server Deployment Script
# Deploys the KDSE MCP Server to a Linux droplet
#===============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo_step() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

echo_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

echo_error() {
    echo -e "${RED}✗ $1${NC}"
}

echo_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Check prerequisites
echo_step "1. Checking Prerequisites"

# Check Docker
if command -v docker &> /dev/null; then
    echo_success "Docker installed: $(docker --version)"
else
    echo_error "Docker is not installed"
    echo "Install Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check Docker Compose
if docker compose version &> /dev/null; then
    echo_success "Docker Compose installed: $(docker compose version --short)"
elif command -v docker-compose &> /dev/null; then
    echo_warning "Using docker-compose (v1)"
else
    echo_error "Docker Compose is not installed"
    exit 1
fi

# Load environment file if exists
if [ -f ".env" ]; then
    echo_success "Loading environment from .env"
    export $(cat .env | grep -v '^#' | xargs)
fi

# Set defaults
export MCP_TRANSPORT=${MCP_TRANSPORT:-http}
export MCP_HTTP_PORT=${MCP_HTTP_PORT:-18181}
export IMAGE_TAG=${IMAGE_TAG:-latest}

echo_step "2. Building Docker Images"

echo "Building kdse-mcp image..."
docker compose build kdse-mcp

echo_step "3. Stopping Existing Containers"

echo "Stopping existing containers..."
docker compose down 2>/dev/null || true

echo_step "4. Starting Containers"

echo "Starting KDSE MCP Server in ${MCP_TRANSPORT} mode..."
docker compose up -d

echo_step "5. Verifying Deployment"

# Wait for container to start
sleep 3

# Check container status
if docker compose ps | grep -q "Up"; then
    echo_success "Container is running"
else
    echo_error "Container failed to start"
    echo "Check logs: docker compose logs"
    exit 1
fi

# Health check (for HTTP mode)
if [ "$MCP_TRANSPORT" = "http" ]; then
    echo "Checking health endpoint..."
    for i in {1..10}; do
        if curl -s http://localhost:${MCP_HTTP_PORT}/health &>/dev/null; then
            echo_success "Health check passed"
            break
        fi
        if [ $i -eq 10 ]; then
            echo_warning "Health check timeout - container may still be starting"
        fi
        sleep 1
    done
fi

echo_step "6. Deployment Complete"

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║          KDSE MCP Server Deployed Successfully                     ║${NC}"
echo -e "${GREEN}╠════════════════════════════════════════════════════════════════════════╣${NC}"
echo -e "${GREEN}║ Transport: ${MCP_TRANSPORT}                                                   ║${NC}"
if [ "$MCP_TRANSPORT" = "http" ]; then
echo -e "${GREEN}║ Port: ${MCP_HTTP_PORT}                                                        ║${NC}"
fi
echo -e "${GREEN}╚════════════════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo "Commands:"
echo "  Status:   docker compose ps"
echo "  Logs:     docker compose logs -f"
echo "  Stop:     docker compose down"
echo "  Restart:  docker compose restart"
echo ""
