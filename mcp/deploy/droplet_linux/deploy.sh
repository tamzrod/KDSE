#!/bin/bash
# =============================================================================
# KDSE MCP Server - Production Deployment Script for Linux Droplets
# =============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Load environment variables
if [ -f .env ]; then
    echo -e "${BLUE}Loading environment from .env...${NC}"
    set -a
    source .env
    set +a
else
    echo -e "${YELLOW}Warning: .env file not found. Using defaults from docker-compose.yml${NC}"
fi

# Transport mode: stdio or http
TRANSPORT_MODE="${MCP_TRANSPORT:-http}"

# Container name (always kdse-mcp for default HTTP service)
CONTAINER_NAME="${KDSE_HTTP_CONTAINER_NAME:-kdse-mcp}"
HTTP_PORT="${MCP_HTTP_PORT:-8080}"

# =============================================================================
# HELPER FUNCTIONS
# =============================================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log_error "Docker daemon is not running. Please start Docker."
        exit 1
    fi
    
    log_success "Docker is installed and running"
}

# =============================================================================
# MAIN DEPLOYMENT
# =============================================================================

deploy() {
    log_info "Starting KDSE MCP Server deployment..."
    log_info "Transport mode: ${TRANSPORT_MODE}"
    
    # Build image locally from source
    log_info "Building Docker image from source..."
    docker compose build --no-cache
    
    # Create network if it doesn't exist
    log_info "Ensuring network exists..."
    docker network create "${KDSE_NETWORK:-kdse-mcp-network}" 2>/dev/null || true
    
    # Stop existing containers
    log_info "Stopping any existing containers..."
    docker compose down 2>/dev/null || true
    
    # Start the service based on transport mode
    log_info "Starting KDSE MCP Server in ${TRANSPORT_MODE} mode..."
    
    if [ "$TRANSPORT_MODE" = "http" ]; then
        docker compose up -d
        log_info "HTTP server will be available on port ${HTTP_PORT}"
    else
        docker compose --profile stdio up -d
        log_info "STDIO server ready for local connections"
    fi
    
    log_success "Deployment complete!"
}

# =============================================================================
# STATUS CHECK
# =============================================================================

status() {
    echo ""
    echo "=== KDSE MCP Server Status ==="
    echo "Transport mode: ${TRANSPORT_MODE}"
    echo ""
    docker compose ps
    echo ""
    
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_success "Container '${CONTAINER_NAME}' is running"
        
        # Show transport-specific info
        if [ "$TRANSPORT_MODE" = "http" ]; then
            echo ""
            echo "=== HTTP Endpoint ==="
            echo "Health: http://localhost:${HTTP_PORT}/health"
            echo "MCP:    http://localhost:${HTTP_PORT}/mcp"
        fi
        
        # Show recent logs
        echo ""
        echo "=== Recent Logs (last 10 lines) ==="
        docker compose logs --tail=10
    else
        log_error "Container '${CONTAINER_NAME}' is NOT running"
    fi
}

# =============================================================================
# HEALTH CHECK
# =============================================================================

health_check() {
    log_info "Performing health verification..."
    
    # Check if container is running
    if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_error "Container is not running"
        return 1
    fi
    
    if [ "$TRANSPORT_MODE" = "http" ]; then
        # Test HTTP health endpoint
        log_info "Testing HTTP health endpoint..."
        if curl -sf "http://localhost:${HTTP_PORT}/health" > /dev/null 2>&1; then
            log_success "HTTP health check passed"
        else
            log_error "HTTP health check failed"
            return 1
        fi
        
        # Test MCP initialization over HTTP
        log_info "Testing MCP initialization..."
        RESPONSE=$(curl -sf -X POST "http://localhost:${HTTP_PORT}/mcp" \
            -H "Content-Type: application/json" \
            -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":0}' 2>/dev/null)
        
        if echo "$RESPONSE" | grep -q "protocolVersion"; then
            log_success "MCP server is responding correctly"
        else
            log_warning "Could not verify MCP response"
        fi
    else
        # STDIO mode - just verify it's running
        HEALTH=$(docker inspect --format='{{.State.Health.Status}}' "${CONTAINER_NAME}" 2>/dev/null || echo "none")
        if [ "$HEALTH" = "healthy" ] || [ "$HEALTH" = "none" ]; then
            log_success "Container is running in STDIO mode"
        else
            log_warning "Container health: ${HEALTH}"
        fi
    fi
}

# =============================================================================
# SWITCH TRANSPORT MODE
# =============================================================================

switch_mode() {
    local new_mode="$1"
    
    if [ "$new_mode" != "stdio" ] && [ "$new_mode" != "http" ]; then
        log_error "Invalid mode. Use 'stdio' or 'http'"
        exit 1
    fi
    
    log_info "Switching to ${new_mode} mode..."
    
    # Stop current containers
    docker compose down
    
    # Update environment
    export MCP_TRANSPORT="$new_mode"
    
    # Restart with new mode
    if [ "$new_mode" = "http" ]; then
        docker compose up -d
    else
        docker compose --profile stdio up -d
    fi
    
    log_success "Switched to ${new_mode} mode"
}

# =============================================================================
# MAIN COMMAND HANDLER
# =============================================================================

case "${1:-deploy}" in
    deploy)
        check_docker
        deploy
        status
        health_check
        ;;
    start)
        check_docker
        if [ "$TRANSPORT_MODE" = "http" ]; then
            docker compose up -d
        else
            docker compose --profile stdio up -d
        fi
        log_success "Started"
        ;;
    stop)
        docker compose stop
        log_success "Stopped"
        ;;
    restart)
        docker compose restart
        log_success "Restarted"
        ;;
    status)
        status
        ;;
    logs)
        shift
        docker compose logs -f "$@"
        ;;
    health|check)
        health_check
        ;;
    pull)
        check_docker
        docker compose pull
        log_success "Images pulled"
        ;;
    build)
        check_docker
        log_info "Building Docker image from source..."
        docker compose build
        log_success "Image built successfully"
        ;;
    switch)
        switch_mode "${2:-http}"
        ;;
    http)
        switch_mode "http"
        ;;
    stdio)
        switch_mode "stdio"
        ;;
    *)
        echo "Usage: $0 {deploy|start|stop|restart|status|logs|health|build|switch}"
        echo ""
        echo "Commands:"
        echo "  deploy  - Full deployment (build, stop old, start new)"
        echo "  start   - Start the container"
        echo "  stop    - Stop the container"
        echo "  restart - Restart the container"
        echo "  status  - Show container status"
        echo "  logs    - View logs (follow mode)"
        echo "  health  - Run health verification"
        echo "  build   - Build Docker image from source"
        echo "  switch  - Switch transport mode (stdio|http)"
        echo ""
        echo "Environment Variables:"
        echo "  MCP_TRANSPORT     - Transport mode: 'stdio' or 'http' (default: http)"
        echo "  MCP_HTTP_PORT     - HTTP port for remote mode (default: 8080)"
        echo ""
        echo "Current mode: ${TRANSPORT_MODE}"
        exit 1
        ;;
esac
