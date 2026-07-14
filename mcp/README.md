# KDSE MCP Server

**Version:** 0.2.0  
**Protocol:** Model Context Protocol (MCP) 2024-11-05

---

## Overview

The KDSE MCP Server provides a Model Context Protocol interface for Knowledge-Driven Software Engineering. It supports **dual transport modes**:

1. **STDIO transport** - For local development
2. **HTTP transport** - For remote deployment

This v0.2 release adds HTTP transport for remote MCP client connections while maintaining STDIO support for local development.

## Design Principles

1. **Dual Transport Architecture**: Same server binary supports STDIO (local) and HTTP (remote)
2. **Protocol Isolation**: MCP protocol handling is strictly separated from KDSE service logic
3. **Structured Output**: All responses are structured JSON, suitable for programmatic consumption
4. **Foundation Only**: Repository reading and advanced features are out of scope for v0.1

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
│  │  • initialize()  • help()  • status()                  │ │
│  └────────────────────────────────────────────────────────┘ │
│                            │                                 │
│                            ▼                                 │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              Tool Implementations                       │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

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
cd mcp

# STDIO mode (local development)
docker compose --profile stdio up -d

# HTTP mode (development server)
docker compose --profile http up -d
```

**Note:** For production deployment, use the configuration in `mcp/deploy/droplet_linux/`

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

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [{
      "type": "text",
      "text": "{\n  \"server\": {\n    \"name\": \"kdse-mcp-server\",\n    \"version\": \"0.2.0\",\n   ...\n  },\n  \"tools\": [...]\n}"
    }]
  }
}
```

### initialize

Returns repository initialization information (static).

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

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [{
      "type": "text",
      "text": "{\n  \"repository\": {\"root\": \"/workspace/project/KDSE\", \"exists\": true},\n  \"module\": \"github.com/kdse/runtime\",\n  \"version\": \"0.2.0\",\n  \"goVersion\": \"1.22.5\",\n  \"features\": [\"help\", \"initialize\", \"status\"],\n  \"components\": {...}\n}"
    }]
  }
}
```

### status

Returns repository status information (static).

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

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [{
      "type": "text",
      "text": "{\n  \"repository\": {\"root\": \"/workspace/project/KDSE\", \"exists\": true},\n  \"git\": {\"available\": true, \"branch\": \"main\", ...},\n  \"kdse\": {\"compliant\": true, ...}\n}"
    }]
  }
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

## Future Extensions

This is v0.2 establishing dual transport support. Planned additions:

- [ ] Repository reading for dynamic responses
- [ ] Development Experience tool
- [ ] Audit tool for compliance checks
- [ ] Architecture search capabilities
- [ ] AI reasoning integration
- [ ] Streaming responses for long-running operations

## License

Apache 2.0 - Same as KDSE runtime
