# KDSE MCP Server - Production Deployment for Linux Droplets

**Version:** 0.2.0  
**Target:** DigitalOcean Droplet / Ubuntu Linux VPS

---

## Overview

This is a production-ready deployment configuration for the KDSE MCP Server on Linux droplets. It supports **dual transport modes**:

1. **STDIO mode** - For local development and debugging
2. **HTTP mode** - For remote MCP client connections over the network

### Features

- Local build from source (no registry required)
- Dual transport support (STDIO + HTTP)
- Auto-restart on failure or system reboot
- Log rotation to prevent disk space issues
- Healthcheck support for HTTP mode
- Externalized configuration via `.env`
- Single-command deployment script
- Switch between transport modes without rebuild

---

## Prerequisites

### System Requirements

- **OS:** Ubuntu 20.04+ / Debian 11+ (or any Linux with Docker support)
- **RAM:** 512MB minimum (1GB recommended)
- **Disk:** 5GB minimum
- **Network:** Internet connectivity for Docker base images
- **Port:** 18181 (or custom) open for HTTP transport

### Required Software

```bash
# Install Docker
curl -fsSL https://get.docker.com | sh

# Install Docker Compose v2 (if not included)
sudo apt-get update
sudo apt-get install -y docker-compose-v2

# Add your user to docker group (optional, for non-root usage)
sudo usermod -aG docker $USER
```

---

## Installation

### Step 1: Clone the Repository

```bash
# SSH to your droplet
ssh root@your-droplet-ip

# Clone the repository
git clone https://github.com/tamzrod/KDSE.git /opt/kdse
cd /opt/kdse/mcp/deploy/droplet_linux
```

### Step 2: Configure Environment

```bash
# Create environment file
cat > .env << 'EOF'
# Transport mode: 'stdio' or 'http' (default: http for remote access)
MCP_TRANSPORT=http

# HTTP server port (default: 18181)
MCP_HTTP_PORT=18181

# KDSE data paths
KDSE_DATA_PATH=/data/kdse
KDSE_REPO_ROOT=/data/kdse
EOF

# Or edit manually
nano .env
```

**Configuration Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `MCP_TRANSPORT` | `http` | Transport mode: `stdio` or `http` |
| `MCP_HTTP_PORT` | `18181` | HTTP port for remote mode |
| `KDSE_BUILD_CONTEXT` | `../..` | Path to MCP source |
| `KDSE_DOCKERFILE` | `Dockerfile` | Dockerfile name |
| `KDSE_IMAGE` | `kdse-mcp` | Docker image name |
| `KDSE_TAG` | `latest` | Image tag |
| `KDSE_HTTP_CONTAINER_NAME` | `kdse-mcp` | HTTP container name |
| `KDSE_STDIO_CONTAINER_NAME` | `kdse-mcp-stdio` | STDIO container name |
| `KDSE_DATA_PATH` | `/data/kdse` | Host path for KDSE data |
| `KDSE_REPO_ROOT` | `/data/kdse` | Container path for KDSE data |

### Step 3: Prepare Data Directory

```bash
# Create the data directory on the host
sudo mkdir -p /data/kdse

# Clone KDSE repository to data directory
sudo git clone https://github.com/tamzrod/KDSE.git /data/kdse

# Set proper ownership
sudo chown -R 1000:1000 /data/kdse
```

---

## Deployment

### Quick Start (HTTP Mode - Remote Access)

```bash
cd /opt/kdse/mcp/deploy/droplet_linux

# Deploy with HTTP transport (default)
./deploy.sh deploy
```

This will:
1. Build the Docker image locally from `mcp/Dockerfile`
2. Create the network
3. Start the HTTP server container on port 18181

### Using the Deploy Script

```bash
cd /opt/kdse/mcp/deploy/droplet_linux

# Full deployment (build + start)
./deploy.sh deploy

# Build separately first
./deploy.sh build
./deploy.sh start

# Switch transport modes
./deploy.sh switch http   # Switch to HTTP mode
./deploy.sh switch stdio  # Switch to STDIO mode
```

### Available Commands

```bash
./deploy.sh deploy  - Build and start (full deployment)
./deploy.sh build   - Build Docker image from source
./deploy.sh start   - Start the container (uses MCP_TRANSPORT setting)
./deploy.sh stop    - Stop the container
./deploy.sh restart - Restart the container
./deploy.sh status  - Show container status and endpoints
./deploy.sh logs    - View logs (follow mode)
./deploy.sh health  - Run health verification
./deploy.sh switch  - Switch transport mode (stdio|http)
```

---

## HTTP Transport (Remote Access)

### Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check endpoint |
| `/mcp` | POST | MCP JSON-RPC endpoint |

### Testing HTTP Endpoint

```bash
# Check health
curl http://localhost:18181/health

# Test MCP initialization
curl -X POST http://localhost:18181/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":0}'
```

### Expected Response

```json
{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "protocolVersion": "2024-11-05",
    "serverInfo": {
      "name": "kdse-mcp",
      "version": "0.2.0"
    },
    "capabilities": {
      "tools": {
        "listChanged": false
      }
    }
  }
}
```

---

## Connecting Remote MCP Clients

### Claude Desktop

Add to your Claude Desktop configuration file (`~/.claude/settings.json`):

```json
{
  "mcpServers": {
    "kdse": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-http",
        "http://YOUR_DROPLET_IP:18181/mcp"
      ]
    }
  }
}
```

**Note:** Claude Desktop supports MCP servers via the official SDK. For production use, install the `@modelcontextprotocol/server-http` package:

```bash
npm install -g @modelcontextprotocol/server-http
```

Then configure Claude Desktop:

```json
{
  "mcpServers": {
    "kdse": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-http", "http://YOUR_DROPLET_IP:18181/mcp"]
    }
  }
}
```

### OpenHands

Configure OpenHands to use the HTTP endpoint:

```bash
# Set the MCP endpoint via environment variable
export OPENHANDS_MCP_SERVER=http://YOUR_DROPLET_IP:18181/mcp

# Or in your OpenHands configuration file (~/.openhands/config.toml)
cat >> ~/.openhands/config.toml << 'EOF'
[mcp]
servers = ["kdse"]

[mcp.servers.kdse]
type = "http"
url = "http://YOUR_DROPLET_IP:18181/mcp"
EOF
```

### Using curl as a Simple Test Client

```bash
#!/bin/bash
# Simple MCP client using curl

ENDPOINT="http://YOUR_DROPLET_IP:18181/mcp"

# Initialize
echo "=== Initialize ==="
curl -s -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":0}' | jq

# List tools
echo "=== List Tools ==="
curl -s -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","params":{},"id":1}' | jq

# Call help tool
echo "=== Help Tool ==="
curl -s -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"help"},"id":2}' | jq
```

### Python MCP Client Example

```python
import requests
import json

class KDSEMCPClient:
    def __init__(self, endpoint):
        self.endpoint = endpoint
        self.session_id = None
    
    def send_request(self, method, params=None, request_id=1):
        response = requests.post(
            self.endpoint,
            json={
                "jsonrpc": "2.0",
                "method": method,
                "params": params or {},
                "id": request_id
            },
            headers={"Content-Type": "application/json"}
        )
        return response.json()
    
    def initialize(self):
        return self.send_request("initialize", {}, 0)
    
    def list_tools(self):
        return self.send_request("tools/list", {}, 1)
    
    def call_tool(self, name, arguments=None):
        return self.send_request("tools/call", {
            "name": name,
            "arguments": arguments or {}
        }, 2)

# Usage
client = KDSEMCPClient("http://YOUR_DROPLET_IP:18181/mcp")
print(client.initialize())
print(client.list_tools())
```

---

## STDIO Mode (Local Development)

For local development with STDIO transport:

```bash
# Deploy in STDIO mode
MCP_TRANSPORT=stdio ./deploy.sh deploy

# Or use docker compose directly
docker compose --profile stdio up -d

# Test with docker exec
docker exec -i kdse-mcp-stdio sh -c 'echo {"jsonrpc":"2.0","method":"initialize","params":{},"id":0}' | ./kdse-mcp
```

---

## Switching Transport Modes

```bash
# Switch to HTTP mode
./deploy.sh switch http

# Switch to STDIO mode
./deploy.sh switch stdio

# Or set environment and restart
export MCP_TRANSPORT=http
docker compose --profile http up -d
```

---

## Updating

### Update from Git

```bash
cd /opt/kdse

# Pull latest changes
git pull origin main

# Rebuild and restart
cd mcp/deploy/droplet_linux
./deploy.sh deploy
```

### Update KDSE Repository Data

```bash
cd /data/kdse
sudo git pull origin main

# Restart to pick up changes
./deploy.sh restart
```

---

## Health Checks

### HTTP Mode Health Check

```bash
# Using the deploy script
./deploy.sh health

# Manual check
curl http://localhost:18181/health
# Expected: {"status":"healthy"}
```

### Container Health Status

```bash
# Check Docker health status
docker inspect --format='{{.State.Health.Status}}' kdse-mcp
```

---

## Viewing Logs

```bash
# Follow mode (real-time)
./deploy.sh logs

# Last 100 lines
docker compose logs --tail=100

# Export to file
docker compose logs > kdse-logs.txt
```

---

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker compose logs

# Verify .env exists
ls -la .env

# Check port availability
netstat -tlnp | grep 18181
```

### Build Fails

```bash
# Ensure you're in the correct directory
pwd  # Should be: /opt/kdse/mcp/deploy/droplet_linux

# Verify Dockerfile exists
ls -la ../../Dockerfile
```

### HTTP Client Connection Refused

```bash
# Check if container is running
docker ps | grep kdse-mcp

# Check firewall rules
sudo ufw allow 18181/tcp

# Verify port binding
docker port kdse-mcp
```

### Permission Denied

```bash
sudo chown -R 1000:1000 /data/kdse
```

### Disk Space Issues

```bash
# Clean unused Docker resources
docker system prune -a

# Check disk usage
df -h
```

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Linux Droplet                              │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              Docker Container (kdse-mcp)         │    │
│  │  ┌───────────────────────────────────────────────┐  │    │
│  │  │         KDSE MCP Server (HTTP mode)            │  │    │
│  │  │         • /health - Health check              │  │    │
│  │  │         • /mcp   - MCP JSON-RPC endpoint       │  │    │
│  │  └───────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────┘    │
│                              │                              │
│  ┌───────────────────────────┴───────────────────────┐    │
│  │              /data/kdse (read-only volume)        │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  Remote MCP Clients ──────────────────► Port 18181          │
│  (Claude Desktop, OpenHands, etc.)                         │
└─────────────────────────────────────────────────────────────┘
```

### Dual Transport Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     KDSE MCP Server                          │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                  Transport Layer                        │ │
│  │  ┌──────────────────┐    ┌──────────────────────────┐  │ │
│  │  │   STDIO Transport │    │   HTTP Transport          │  │ │
│  │  │   (Local Dev)     │    │   (Remote Deployment)     │  │ │
│  │  │                   │    │   • /health endpoint      │  │ │
│  │  │   stdin/stdout    │    │   • /mcp endpoint        │  │ │
│  │  │   JSON-RPC        │    │   CORS enabled           │  │ │
│  │  └──────────────────┘    └──────────────────────────┘  │ │
│  └────────────────────────────────────────────────────────┘ │
│                            │                                 │
│                            ▼                                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │               KDSE Service Layer                        │ │
│  │  • initialize()  • help()  • status()                  │ │
│  └────────────────────────────────────────────────────────┘ │
│                            │                                 │
│                            ▼                                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              Tool Implementations                      │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## Security Considerations

1. **Data is read-only:** The KDSE data is mounted as read-only
2. **Non-root user:** The container runs as UID 1000
3. **Minimal image:** Based on Alpine Linux
4. **Network isolation:** Dedicated bridge network
5. **CORS enabled:** For cross-origin requests from web clients
6. **No authentication:** Currently no auth (add reverse proxy/AuthN for production)

### Production Hardening (Recommended)

For production deployments without authentication on the MCP server itself, consider:

- **Reverse proxy with authentication:** Put Nginx or Caddy in front with Basic Auth or JWT
- **Firewall:** Restrict port 18181 to specific IPs via UFW
- **TLS:** Add HTTPS via reverse proxy (let's encrypt)
- **Rate limiting:** Configure at reverse proxy level

Example with Caddy:

```bash
# Caddyfile
kdse.example.com {
    reverse_proxy localhost:18181 {
        header_up X-Real-IP {remote_host}
    }
    
    basicauth /* {
        your-username JHash123...
    }
}
```

---

## Support

For issues, please open an issue at:
https://github.com/tamzrod/KDSE/issues

---

## License

Apache 2.0 - Same as KDSE runtime
