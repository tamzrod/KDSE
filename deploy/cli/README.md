# KDSE CLI - Isolated Deployment

This directory contains the authoritative deployment configuration for the KDSE CLI runtime.

## Contents

- `docker-compose.yml` - CLI container configuration
- `Dockerfile` - CLI image build
- `.env.example` - Environment variables template
- `README.md` - This file

## Quick Start

```bash
# Navigate to deployment directory
cd deploy/cli

# Build and run
docker compose up -d

# View status
docker compose run --rm kdse-runtime status
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `IMAGE_TAG` | `latest` | Docker image tag |
| `KDSE_REPO_PATH` | `../..` | Repository path (mounted into container) |

## Usage

### Run CLI commands

```bash
# Initialize workspace
docker compose run --rm kdse-runtime init

# Check status
docker compose run --rm kdse-runtime status

# View help
docker compose run --rm kdse-runtime help
```

## Build

```bash
# Build image
docker compose build

# Pull latest
docker compose pull
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
