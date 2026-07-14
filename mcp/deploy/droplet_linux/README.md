# KDSE MCP Server - Production Deployment for Linux Droplets

**Version:** 0.1.0  
**Target:** DigitalOcean Droplet / Ubuntu Linux VPS

---

## Overview

This is a production-ready deployment configuration for the KDSE MCP Server on Linux droplets. It is completely isolated from the development deployment.

### Features

- Production-grade Docker Compose configuration
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
- **Network:** Internet connectivity for Docker images

### Required Software

```bash
# Install Docker
curl -fsSL https://get.docker.com | sh

# Install Docker Compose (if not included)
sudo apt-get update
sudo apt-get install -y docker-compose

# Add your user to docker group (optional, for non-root usage)
sudo usermod -aG docker $USER
```

### Firewall Requirements

```bash
# Allow SSH
sudo ufw allow 22

# If using HTTP transport in future (port 8080)
# sudo ufw allow 8080
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

# Edit the .env file with your settings
nano .env
```

**Important Configuration:**

| Variable | Default | Description |
|----------|---------|-------------|
| `KDSE_IMAGE` | `kdse-mcp-server` | Docker image name |
| `KDSE_TAG` | `latest` | Image tag (use `latest` or version like `0.1.0`) |
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

### Step 4: Build or Pull Docker Image

**Option A: Pull Pre-built Image**
```bash
# Edit .env to set KDSE_TAG=latest
./deploy.sh pull
```

**Option B: Build Local Image**
```bash
# Build from local source
cd /opt/kdse/mcp-server
docker build -t kdse-mcp-server:local ..

# Edit .env to set KDSE_TAG=local
cd /opt/kdse/mcp/deploy/droplet_linux
nano .env  # Change KDSE_TAG=local
```

---

## Deployment

### Full Deployment

```bash
cd /opt/kdse/mcp/deploy/droplet_linux
./deploy.sh deploy
```

This script will:
1. Pull the latest Docker image (or use local if configured)
2. Stop any existing container
3. Start the new container
4. Verify deployment

### Individual Commands

```bash
# Start the service
./deploy.sh start

# Stop the service
./deploy.sh stop

# Restart the service
./deploy.sh restart

# Check status
./deploy.sh status

# View logs
./deploy.sh logs

# Health check
./deploy.sh health
```

---

## Updating

### Update to Latest Release

```bash
cd /opt/kdse/mcp/deploy/droplet_linux

# Pull latest image
./deploy.sh pull

# Redeploy
./deploy.sh deploy
```

### Update KDSE Repository Data

```bash
# Update from git
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

This stops the container gracefully (SIGTERM).

### Force Stop

```bash
docker stop kdse-mcp-server
```

### Remove Container

```bash
docker compose down
```

**Note:** This removes the container but not volumes or images.

---

## Viewing Logs

### Follow Mode (Real-time)

```bash
./deploy.sh logs
```

### Last N Lines

```bash
docker compose logs --tail=100
```

### Since Timestamp

```bash
docker compose logs --since "2024-01-01T00:00:00"
```

### Filter by Level

```bash
docker compose logs --level error
```

### Export Logs

```bash
# Export to file
docker compose logs > kdse-mcp-logs.txt

# Export with timestamps
docker compose logs -t > kdse-mcp-logs.txt
```

---

## Troubleshooting

### Container Won't Start

```bash
# Check container logs
docker compose logs

# Check Docker daemon status
systemctl status docker

# Verify .env file exists
ls -la .env
```

### Port Already in Use

```bash
# Check what's using the port
sudo netstat -tlnp | grep 8080

# Or if using stdio mode, check if container is already running
docker ps -a | grep kdse
```

### Permission Denied

```bash
# Fix data directory permissions
sudo chown -R 1000:1000 /data/kdse

# Or rebuild image without user switching
# Edit Dockerfile and remove 'USER kdse' line
```

### Image Pull Fails

```bash
# Check Docker Hub connectivity
docker pull hello-world

# Retry with explicit tag
docker pull kdse-mcp-server:latest
```

### Container Exits Immediately

```bash
# Check exit code
docker compose ps

# Inspect container
docker inspect kdse-mcp-server

# Run interactively to see errors
docker run -it --rm kdse-mcp-server:latest
```

### Health Check Fails

```bash
# Health checks are informational for stdio mode
# If container is running, it should work fine

# Verify MCP communication
echo '{"jsonrpc":"2.0","method":"initialize","params":{},"id":0}' | \
  docker exec -i kdse-mcp-server ./kdse-mcp-server
```

### Disk Space Issues

```bash
# Check disk usage
df -h

# Clean Docker resources
docker system prune -a

# Reduce log retention (already configured)
# max-size: "10m", max-file: "3"
```

### Network Issues

```bash
# Check Docker network
docker network ls
docker network inspect kdse-mcp-network

# Recreate network
docker network rm kdse-mcp-network
./deploy.sh deploy
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
│  │  └───────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────┘    │
│                              │                              │
│  ┌───────────────────────────┴───────────────────────┐    │
│  │              /data/kdse (read-only mount)          │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  MCP Client ──────────────────► stdio                       │
└─────────────────────────────────────────────────────────────┘
```

---

## Environment Variables Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `KDSE_IMAGE` | `kdse-mcp-server` | Docker image name |
| `KDSE_TAG` | `latest` | Image tag |
| `KDSE_CONTAINER_NAME` | `kdse-mcp-server` | Fixed container name |
| `KDSE_NETWORK` | `kdse-mcp-network` | Docker network name |
| `KDSE_DATA_PATH` | `/data/kdse` | Host data path |
| `KDSE_REPO_ROOT` | `/data/kdse` | Container data path |
| `MCP_STDIO` | `true` | Use stdio transport |

---

## Service Management with systemd (Optional)

For automatic startup without Docker Compose:

```bash
# Create systemd service
sudo nano /etc/systemd/system/kdse-mcp.service
```

```ini
[Unit]
Description=KDSE MCP Server
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/kdse/mcp/deploy/droplet_linux
ExecStart=/usr/local/bin/docker-compose up -d
ExecStop=/usr/local/bin/docker-compose stop
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start
sudo systemctl enable kdse-mcp
sudo systemctl start kdse-mcp

# Check status
sudo systemctl status kdse-mcp
```

---

## Security Considerations

1. **Data is read-only:** The KDSE data is mounted as read-only in the container
2. **Non-root user:** The container runs as a non-root user (UID 1000)
3. **No secrets stored:** No sensitive data in the container
4. **Minimal image:** Based on Alpine Linux for minimal attack surface
5. **Network isolation:** Container uses a dedicated bridge network

---

## Support

For issues with the KDSE MCP Server, please open an issue at:
https://github.com/tamzrod/KDSE/issues

---

## License

Apache 2.0 - Same as KDSE runtime
