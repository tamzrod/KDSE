# KDSE MCP Server

**Version:** 0.3.0  
**Protocol:** Model Context Protocol (MCP) 2024-11-05

---

## Overview

The KDSE MCP Server provides a Model Context Protocol interface for Knowledge-Driven Software Engineering. It supports **dual transport modes**:

1. **STDIO transport** - For local development
2. **HTTP transport** - For remote deployment

This v0.3 release introduces the **.kdse/ workspace architecture** - a core architectural principle that KDSE owns only its own workspace and never pollutes the user's repository.

## Architectural Principle

**The repository is the user's library. The `.kdse` directory is the librarian's workspace.**

KDSE may observe, analyze, index, and document the repository. KDSE must not reorganize or pollute the repository outside its own `.kdse` workspace unless explicitly instructed by the user.

## Design Principles

1. **Workspace Isolation**: All KDSE artifacts are stored under `.kdse/` to avoid polluting the repository
2. **Dual Transport Architecture**: Same server binary supports STDIO (local) and HTTP (remote)
3. **Protocol Isolation**: MCP protocol handling is strictly separated from KDSE service logic
4. **Structured Output**: All responses are structured JSON, suitable for programmatic consumption
5. **Legacy Migration**: Automatic detection and migration of legacy KDSE directories

## Architecture

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
│  │  • initialize()  • help()  • status()  • collect()     │ │
│  │  • foundation()  • audit()  • migrate()                │ │
│  └────────────────────────────────────────────────────────┘ │
│                            │                                 │
│                            ▼                                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              Workspace Manager (.kdse/)                 │ │
│  │  • foundation/  • knowledge/  • context/               │ │
│  │  • artifacts/   • runtime/   • sessions/               │ │
│  │  • reports/     • cache/     • confidence/             │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## .kdse Workspace Structure

All KDSE-managed artifacts reside under the `.kdse/` directory:

```
.kdse/
├── foundation/      # Foundation documentation
├── knowledge/       # Derived engineering knowledge
├── context/         # Context handoff files
├── artifacts/       # Collected artifacts inventory
├── runtime/         # Runtime state and configuration
├── sessions/        # Session management
├── confidence/      # Confidence metrics
├── operational/     # Operational knowledge
├── developmental/   # Development experience
├── reports/         # Audit and analysis reports
├── cache/           # Cached analysis results
└── normalized/      # Normalized documentation
```

**Key Principle**: The repository root should remain untouched except for files explicitly requested by the user.

## Transport Configuration

The transport is selected via the `MCP_TRANSPORT` environment variable:

| Transport | Environment | Use Case |
|-----------|-------------|----------|
| STDIO | `MCP_TRANSPORT=stdio` | Local development |
| HTTP | `MCP_TRANSPORT=http` | Remote deployment |

**HTTP Port:** Configure via `MCP_HTTP_PORT` (default: 8080)

## Quick Start

### Local Development (STDIO)

```bash
cd mcp
go mod download
go run .
```

### Remote Deployment (HTTP)

```bash
# Set transport and port
export MCP_TRANSPORT=http
export MCP_HTTP_PORT=8080

# Run
go run .
```

### Using Docker Compose (Development)

```bash
cd mcp/deploy/development

# STDIO mode (local development)
docker compose --profile stdio up -d

# HTTP mode (development server)
docker compose --profile http up -d
```

**Note:** For production deployment, use `mcp/deploy/droplet_linux/`

## Available Tools

### help

Returns information about all available KDSE MCP tools.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "help"
  },
  "id": 1
}
```

### initialize

Initializes the KDSE `.kdse/` workspace. Creates the workspace directory if it doesn't exist. Returns workspace information including migration status for legacy directories.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "initialize"
  },
  "id": 2
}
```

### status

Returns current repository status including workspace state and compliance indicators.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "status"
  },
  "id": 3
}
```

### collect

Collects and catalogs engineering artifacts into `.kdse/artifacts/`.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "collect"
  },
  "id": 4
}
```

### foundation

Returns or creates foundation documentation under `.kdse/foundation/`.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "foundation"
  },
  "id": 5
}
```

### audit

Generates audit reports under `.kdse/reports/`.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "audit"
  },
  "id": 6
}
```

### migrate

Migrates any legacy KDSE directories from repository root to `.kdse/`.

**Request:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "migrate"
  },
  "id": 7
}
```

## HTTP Transport Endpoints

When running in HTTP mode, the server exposes:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/mcp` | POST | MCP JSON-RPC endpoint |

### Testing HTTP Endpoint

```bash
# Check health
curl http://localhost:8080/health

# Test MCP initialization
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":0}'
```

## MCP Protocol

### Initialization

The MCP client must send an `initialize` request before using any tools:

```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "openhands",
      "version": "1.0.0"
    }
  },
  "id": 0
}
```

### Tool Listing

List available tools with `tools/list`:

```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "params": {},
  "id": 1
}
```

## Legacy Migration

If your repository has legacy KDSE directories in the root:

```
foundation/
knowledge/
context/
artifacts/
```

Run the `migrate` tool to move them to `.kdse/`:

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "migrate"
  },
  "id": 8
}
```

This ensures all KDSE artifacts reside under `.kdse/` per the architectural principle.

## License

Apache 2.0 - Same as KDSE runtime
