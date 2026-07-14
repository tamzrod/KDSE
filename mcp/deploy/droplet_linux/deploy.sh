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

# Container name from env or default
CONTAINER_NAME="${KDSE_CONTAINER_NAME:-kdse-mcp-server}"

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
    
    # Pull latest image or build from local
    if [ "${KDSE_TAG:-latest}" != "local" ]; then
        log_info "Pulling Docker image: ${KDSE_IMAGE:-kdse-mcp-server}:${KDSE_TAG:-latest}"
        docker compose pull
    else
        log_warning "Using local image (KDSE_TAG=local). Run 'docker build' first."
    fi
    
    # Create network if it doesn't exist
    log_info "Ensuring network exists..."
    docker network create "${KDSE_NETWORK:-kdse-mcp-network}" 2>/dev/null || true
    
    # Stop existing container if running
    if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_info "Stopping existing container..."
        docker compose stop
    fi
    
    # Remove old container if exists
    if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_info "Removing old container..."
        docker compose rm -f
    fi
    
    # Start the service
    log_info "Starting KDSE MCP Server..."
    docker compose up -d
    
    log_success "Deployment complete!"
}

# =============================================================================
# STATUS CHECK
# =============================================================================

status() {
    echo ""
    echo "=== Container Status ==="
    docker compose ps
    echo ""
    
    if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        log_success "Container '${CONTAINER_NAME}' is running"
        
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
    
    # Check container health status
    HEALTH=$(docker inspect --format='{{.State.Health.Status}}' "${CONTAINER_NAME}" 2>/dev/null || echo "none")
    
    if [ "$HEALTH" = "healthy" ]; then
        log_success "Container health: ${HEALTH}"
    elif [ "$HEALTH" = "none" ]; then
        log_warning "No healthcheck configured (this is normal for stdio mode)"
        log_info "Verifying container is responsive..."
        
        # For stdio containers, just verify it's running
        if docker ps --format '{{.Status}}' | grep -q "Up"; then
            log_success "Container is up and running"
        fi
    else
        log_warning "Container health: ${HEALTH}"
    fi
    
    # Test MCP communication if possible
    log_info "Testing MCP initialization..."
    if echo '{"jsonrpc":"2.0","method":"initialize","params":{},"id":0}' | timeout 5 docker compose exec -T "${CONTAINER_NAME}" ./kdse-mcp-server 2>/dev/null | grep -q "protocolVersion"; then
        log_success "MCP server is responding correctly"
    else
        log_warning "Could not verify MCP response (this may be expected for stdio mode)"
    fi
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
        docker compose up -d
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
    *)
        echo "Usage: $0 {deploy|start|stop|restart|status|logs|health|pull}"
        echo ""
        echo "Commands:"
        echo "  deploy  - Full deployment (pull, stop old, start new)"
        echo "  start   - Start the container"
        echo "  stop    - Stop the container"
        echo "  restart - Restart the container"
        echo "  status  - Show container status"
        echo "  logs    - View logs (follow mode)"
        echo "  health  - Run health verification"
        echo "  pull    - Pull latest images"
        exit 1
        ;;
esac
