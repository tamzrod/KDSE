# KDSE MCP Server

**Version:** 0.1.0  
**Protocol:** Model Context Protocol (MCP) 2024-11-05

---

## Overview

The KDSE MCP Server provides a Model Context Protocol interface for Knowledge-Driven Software Engineering. It enables AI assistants like OpenHands to communicate with KDSE through the standard MCP protocol.

This v0.1 release establishes the MCP communication foundation with static responses.

## Design Principles

1. **Static Responses**: v0.1 returns static data to prove MCP communication works
2. **Protocol Isolation**: MCP protocol handling is strictly separated from KDSE service logic
3. **Structured Output**: All responses are structured JSON, suitable for programmatic consumption
4. **Foundation Only**: Repository reading and advanced features are out of scope for v0.1

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
    "content": [{
      "type": "text",
      "text": "{\n  \"server\": {\n    \"name\": \"kdse-mcp-server\",\n    \"version\": \"0.1.0\",\n    ...\n  },\n  \"tools\": [...]\n}"
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
      "text": "{\n  \"repository\": {\"root\": \"/workspace/project/KDSE\", \"exists\": true},\n  \"module\": \"github.com/kdse/runtime\",\n  \"version\": \"0.1.0\",\n  \"goVersion\": \"1.22.5\",\n  \"features\": [\"help\", \"initialize\", \"status\"],\n  \"components\": {...}\n}"
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
│  │  • help - Tool information (static)                      ││
│  │  • initialize - Repository metadata (static)             ││
│  │  • status - Repository state (static)                    ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

## Future Extensions

This is v0.1 establishing the MCP foundation. Planned additions:

- [ ] Repository reading for dynamic responses
- [ ] Development Experience tool
- [ ] Audit tool for compliance checks
- [ ] Architecture search capabilities
- [ ] AI reasoning integration
- [ ] HTTP transport mode

## License

Apache 2.0 - Same as KDSE runtime
