# KDSE MCP Server - Linux Droplet Deployment

This directory contains deployment scripts and configuration for deploying the KDSE MCP Server on a Linux droplet (VPS).

## Contents

- `deploy.sh` - Main deployment script
- `docker-compose.yml` - Production Docker Compose configuration
- `.env.example` - Environment variables template
- `setup-service.sh` - Systemd service setup script
- `README.md` - This file

## Quick Start

```bash
# Navigate to deployment directory
cd mcp/deploy/droplet_linux

# Run deployment
./deploy.sh
```

## Prerequisites

- Ubuntu 20.04+ or Debian 11+ (recommended)
- Docker 20.10+ and Docker Compose v2
- 1GB RAM minimum (2GB recommended)
- 10GB disk space

## Deployment Options

### Option 1: Quick Deploy (Recommended)

```bash
./deploy.sh
```

### Option 2: Custom Environment

```bash
cp .env.example .env
# Edit .env with your configuration
./deploy.sh --env-file .env
```

### Option 3: Docker Compose Direct

```bash
docker compose -f docker-compose.yml up -d
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MCP_TRANSPORT` | `stdio` | Transport mode: `stdio` or `http` |
| `MCP_HTTP_PORT` | `8080` | HTTP server port (if using http transport) |
| `KDSE_REPO_ROOT` | `/app` | Repository root inside container |
| `IMAGE_TAG` | `latest` | Docker image tag |

### STDIO Mode (Default)

STDIO mode is designed for local development and AI assistant integration:

```bash
docker compose --profile stdio up -d
```

### HTTP Mode (Production)

HTTP mode enables remote access via REST API:

```bash
docker compose --profile http up -d
```

## Systemd Service

For production deployments, set up a systemd service:

```bash
sudo ./setup-service.sh
sudo systemctl enable kdse-mcp
sudo systemctl start kdse-mcp
```

## Verification

Check deployment status:

```bash
# View container status
docker compose ps

# View logs
docker compose logs -f

# Health check (HTTP mode)
curl http://localhost:8080/health
```

## Troubleshooting

### Container fails to start

```bash
docker compose logs kdse-mcp-stdio
```

### Port already in use

Change the port in `.env`:
```
MCP_HTTP_PORT=8081
```

### Permission denied

Ensure Docker has proper permissions:
```bash
sudo usermod -aG docker $USER
# Log out and back in
```

## Update

To update to the latest version:

```bash
git pull
docker compose pull
docker compose up -d
```

## Undeploy

```bash
docker compose down
# If using systemd:
sudo systemctl stop kdse-mcp
sudo systemctl disable kdse-mcp
```
