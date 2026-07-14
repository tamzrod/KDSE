# KDSE MCP Server - Development Deployment

**Version:** 0.2.0  
**Classification:** Development Only

---

## Overview

This deployment package is for **local development** only. It provides Docker Compose configurations for running the KDSE MCP Server in both STDIO and HTTP modes.

For production deployment, use `mcp/deploy/droplet_linux/`.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose v2

## Quick Start

### Build the Docker Image

```bash
cd mcp/deploy/development
docker compose build
```

### Run in STDIO Mode

```bash
docker compose --profile stdio up -d
docker exec -i kdse-mcp-stdio sh -c 'echo {"jsonrpc":"2.0","method":"help","params":{},"id":0}' | docker exec -i kdse-mcp-stdio ./kdse-mcp
```

### Run in HTTP Mode

```bash
docker compose --profile http up -d
curl http://localhost:8080/health
```

## Available Profiles

| Profile | Description |
|---------|-------------|
| `stdio` | STDIO transport for local development |
| `http` | HTTP transport for testing HTTP endpoints |
| `development` | Meta profile for all development modes |

## Stopping

```bash
docker compose down
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MCP_HTTP_PORT` | `8080` | HTTP server port |
| `KDSE_REPO_PATH` | `../../..` | Path to KDSE repository |

## Package Contents

```
mcp/deploy/development/
├── docker-compose.yml  # Docker Compose configuration
└── README.md           # This file
```
