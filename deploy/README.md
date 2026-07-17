# KDSE Deployment Architecture

This directory contains the authoritative deployment configurations for KDSE.

## Deployment Targets

| Target | Directory | Purpose |
|--------|-----------|---------|
| **MCP** | `deploy/mcp/` | MCP Server for AI integration |
| **CLI** | `deploy/cli/` | CLI runtime for command-line operations |

## Quick Reference

### MCP Server Deployment

```bash
cd deploy/mcp
cp .env.example .env
docker compose up -d
```

### CLI Deployment

```bash
cd deploy/cli
cp .env.example .env
docker compose up -d
```

## Architecture

```
deploy/
├── README.md          # This file
├── mcp/               # MCP Server deployment
│   ├── docker-compose.yml
│   ├── Dockerfile
│   ├── .env.example
│   └── README.md
└── cli/               # CLI deployment
    ├── docker-compose.yml
    ├── Dockerfile
    ├── .env.example
    └── README.md
```

## Requirements

- Docker 20.10+
- Docker Compose v2
- For MCP: `nginx-proxy-manager_default` network

## Network Configuration

The MCP deployment joins the external `nginx-proxy-manager_default` network for integration with Nginx Proxy Manager.

The CLI deployment uses an isolated bridge network (`kdse-cli-network`).
