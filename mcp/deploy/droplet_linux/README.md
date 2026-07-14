# KDSE MCP Server - Production Deployment for Linux Droplets

**Version:** 0.1.0  
**Target:** DigitalOcean Droplet / Ubuntu Linux VPS

---

## Overview

This is a production-ready deployment configuration for the KDSE MCP Server on Linux droplets. It builds the image locally from source and is completely isolated from the development deployment.

### Features

- Local build from source (no registry required)
- Auto-restart on failure or system reboot
- Log rotation to prevent disk space issues
- Healthcheck support
- Externalized configuration via `.env`
- Single-command deployment script

---

## Prerequisites

### System Requirements

- **OS:** Ubuntu 20.04+ / Debian 11+ (or any Linux with Docker support)
- **RAM:** 512MB minimum (1GB recommended)
- **Disk:** 5GB minimum
- **Network:** Internet connectivity for Docker base images

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
# Copy the example environment file
cp .env.example .env

# Edit the .env file with your settings (defaults work for most cases)
nano .env
```

**Configuration Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `KDSE_BUILD_CONTEXT` | `../../../mcp-server` | Path to mcp-server source |
| `KDSE_DOCKERFILE` | `Dockerfile` | Dockerfile name |
| `KDSE_IMAGE` | `kdse-mcp-server` | Docker image name |
| `KDSE_TAG` | `latest` | Image tag |
| `KDSE_CONTAINER_NAME` | `kdse-mcp-server` | Fixed container name |
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

### Quick Start (Single Command)

```bash
cd /opt/kdse/mcp/deploy/droplet_linux
docker compose up -d
```

This will:
1. Build the Docker image locally from `mcp-server/Dockerfile`
2. Create the network
3. Start the container

### Using the Deploy Script

```bash
cd /opt/kdse/mcp/deploy/droplet_linux

# Full deployment (build + start)
./deploy.sh deploy

# Or build separately first
./deploy.sh build
./deploy.sh start
```

### Available Commands

```bash
./deploy.sh deploy  - Build and start (full deployment)
./deploy.sh build   - Build Docker image from source
./deploy.sh start   - Start the container
./deploy.sh stop    - Stop the container
./deploy.sh restart - Restart the container
./deploy.sh status  - Show container status
./deploy.sh logs    - View logs (follow mode)
./deploy.sh health  - Run health verification
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

## Starting

### Automatic Start (Recommended)

The container is configured with `restart: unless-stopped`, which means:

- **On boot:** Container starts automatically
- **On failure:** Container restarts automatically
- **On manual stop:** Container stays stopped

```bash
# After reboot, verify container is running
./deploy.sh status
```

### Manual Start

```bash
./deploy.sh start
```

---

## Stopping

### Graceful Stop

```bash
./deploy.sh stop
```

### Remove Container

```bash
docker compose down
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

## Testing

Verify MCP communication:

```bash
echo '{"jsonrpc":"2.0","method":"initialize","params":{},"id":0}' | \
  docker exec -i kdse-mcp-server ./kdse-mcp-server
```

---

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker compose logs

# Verify .env exists
ls -la .env
```

### Build Fails

```bash
# Ensure you're in the correct directory
pwd  # Should be: /opt/kdse/mcp/deploy/droplet_linux

# Verify Dockerfile exists
ls -la ../../../mcp-server/Dockerfile
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
│  │              Docker Container                         │    │
│  │  ┌───────────────────────────────────────────────┐  │    │
│  │  │         KDSE MCP Server (stdio mode)          │  │    │
│  │  │         Binary only (built from source)         │  │    │
│  │  └───────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────┘    │
│                              │                              │
│  ┌───────────────────────────┴───────────────────────┐    │
│  │              /data/kdse (read-only volume)        │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  MCP Client ──────────────────► stdio                       │
└─────────────────────────────────────────────────────────────┘
```

---

## Security Considerations

1. **Data is read-only:** The KDSE data is mounted as read-only
2. **Non-root user:** The container runs as UID 1000
3. **Minimal image:** Based on Alpine Linux
4. **Network isolation:** Dedicated bridge network

---

## Support

For issues, please open an issue at:
https://github.com/tamzrod/KDSE/issues

---

## License

Apache 2.0 - Same as KDSE runtime
