# KDSE MCP Server

**Version:** 0.1.0  
**Protocol:** Model Context Protocol (MCP) 2024-11-05

---

## Overview

The KDSE MCP Server provides a Model Context Protocol interface to Knowledge-Driven Software Engineering repository information. It enables AI assistants like OpenHands to query structured KDSE repository data through the standard MCP protocol.

## Design Principles

1. **Read-Only**: The server only reads repository data; it never modifies anything.
2. **Knowledge Separation**: The MCP server does not own repository knowledge. Repository knowledge remains the source of truth.
3. **Protocol Isolation**: MCP protocol handling is strictly separated from KDSE service logic.
4. **Structured Output**: All responses are structured JSON, suitable for programmatic consumption.

## Quick Start

### Using Docker Compose

```bash
cd mcp-server
docker compose up --build
```

### Local Development

```bash
cd mcp-server
go mod download
go run .
```

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
    "server": {
      "name": "kdse-mcp-server",
      "version": "0.1.0",
      "description": "Model Context Protocol server for Knowledge-Driven Software Engineering",
      "protocol": "2024-11-05"
    },
    "tools": [...]
  }
}
```

### initialize

Returns repository initialization information including module name, version, and supported features.

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
    "repository": {
      "root": "/workspace/project/KDSE",
      "exists": true
    },
    "module": "github.com/kdse/runtime",
    "version": "0.1.0",
    "goVersion": "1.22.5",
    "features": ["help", "initialize", "status"],
    "components": {
      "commands": ["kdse"],
      "packages": ["collect", "config", "context", "detection", "normalize", "report", "state", "types"]
    }
  }
}
```

### status

Returns current repository status information including git state, file counts, and KDSE compliance.

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
    "repository": {
      "root": "/workspace/project/KDSE",
      "exists": true
    },
    "git": {
      "available": true,
      "branch": "main",
      "commit": "abc12345",
      "has_changes": false
    },
    "files": {
      "total": 150,
      "by_type": {".go": 45, ".md": 80, ".yml": 10},
      "by_location": {"cmd": 5, "internal": 40, "runtime": 30, "docs": 75}
    },
    "kdse": {
      "compliant": true,
      "has_readme": true,
      "has_foundation": true,
      "has_runtime": true,
      "has_go_mod": true,
      "has_docker": false,
      "has_mcp_server": true
    }
  }
}
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

### Tool Calling

Tools are called using `tools/call`:

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

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `KDSE_REPO_ROOT` | `/workspace/project/KDSE` | Path to KDSE repository |
| `MCP_STDIO` | `true` | Use stdio transport (currently only mode) |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     MCP Client (OpenHands)                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ JSON-RPC 2.0 over stdio
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    KDSE MCP Server                           │
│  ┌─────────────────────────────────────────────────────────┐│
│  │              Protocol Handler (main.go)                 ││
│  │  • JSON-RPC 2.0 parsing                                  ││
│  │  • Request routing                                       ││
│  │  • Response formatting                                   ││
│  └─────────────────────────────────────────────────────────┘│
│                              │                               │
│                              ▼                               │
│  ┌─────────────────────────────────────────────────────────┐│
│  │              Tool Handler (tools/tools.go)               ││
│  │  • help - Tool information                              ││
│  │  • initialize - Repository metadata                     ││
│  │  • status - Repository state                            ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Read-only
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    KDSE Repository                           │
│                    (Source of Truth)                         │
└─────────────────────────────────────────────────────────────┘
```

## Future Extensions

This is v0.1 with minimal functionality. Planned additions:

- [ ] Experience tool for development experience queries
- [ ] Audit tool for compliance checks
- [ ] Architecture search capabilities
- [ ] AI reasoning integration
- [ ] HTTP transport mode

## License

Apache 2.0 - Same as KDSE runtime
