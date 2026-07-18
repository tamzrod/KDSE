# KDSE MCP Server - Production Deployment

This directory contains the authoritative deployment configuration for the KDSE MCP Server.

**Target:** DigitalOcean with Nginx Proxy Manager  
**Network:** `nginx-proxy-manager_default` (external)

## Contents

- `docker-compose.yml` - MCP container configuration
- `.env.example` - Environment variables template
- `README.md` - This file

## Prerequisites

- Docker 20.10+ and Docker Compose v2
- Existing `nginx-proxy-manager_default` network (from Nginx Proxy Manager)
- Nginx Proxy Manager installation on the droplet

## Quick Start

```bash
# Navigate to deployment directory
cd deploy/mcp

# Copy and configure environment
cp .env.example .env
# Edit .env if needed

# Build and start
docker compose up -d

# Verify health
curl http://localhost:18181/health
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MCP_TRANSPORT` | `http` | Transport mode: `stdio` or `http` |
| `MCP_HTTP_PORT` | `18181` | HTTP server port |
| `KDSE_REPO_ROOT` | `/repo` | Repository root inside container |
| `KDSE_REPO_HOST_PATH` | `../../..` | Host path to mount |
| `IMAGE_TAG` | `production` | Docker image tag |
| `LOG_MAX_SIZE` | `10m` | Max log file size |
| `LOG_MAX_FILES` | `3` | Number of log files |

### Transport Modes

#### HTTP Mode (Production)
HTTP mode enables remote access via REST API through Nginx Proxy Manager.

```bash
MCP_TRANSPORT=http
MCP_HTTP_PORT=18181
```

#### STDIO Mode (Development)
STDIO mode is for local development and direct AI assistant integration.

```bash
MCP_TRANSPORT=stdio
```

## Network Configuration

The MCP deployment joins the external `nginx-proxy-manager_default` network for proxy integration.

```yaml
networks:
  proxy:
    external: true
    name: nginx-proxy-manager_default
```

## Volumes

The local repository is mounted into the container at `/repo`:

```yaml
volumes:
  - ${KDSE_REPO_HOST_PATH:-../../..}:/repo
```

## Verification

```bash
# View container status
docker compose ps

# View logs
docker compose logs -f kdse-mcp

# Health check
curl http://localhost:18181/health
```

## Troubleshooting

### Container fails to start

```bash
docker compose logs kdse-mcp
```

### Port already in use

Change the port in `.env`:
```
MCP_HTTP_PORT=8081
```

### Network not found

Ensure `nginx-proxy-manager_default` network exists (it should be created by Nginx Proxy Manager):

```bash
docker network ls | grep nginx-proxy-manager
```

If not found, create it:
```bash
docker network create nginx-proxy-manager_default
```

### Permission denied

```bash
sudo usermod -aG docker $USER
# Log out and back in
```

## Update

```bash
git pull
docker compose build
docker compose up -d
```

## Undeploy

```bash
docker compose down
```
