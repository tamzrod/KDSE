# MCP Runtime

**Document Version:** 1.0  
**Type:** Normative  
**Effective Date:** 2026-07-17

---

## Purpose

This document defines the MCP (Model Context Protocol) Runtime for KDSE. The MCP Runtime is a **thin adapter** that translates MCP tool requests into Workspace Engine calls.

---

## Core Principle: Thin Adapter

The MCP Runtime MUST be a thin adapter. It contains NO business logic.

### Responsibilities

| Responsibility | Description | Required |
|----------------|-------------|----------|
| Receive tool requests | Receive MCP tool invocations | YES |
| Call Workspace Engine | Invoke engine methods | YES |
| Return structured responses | Return JSON-RPC responses | YES |
| Handle errors | Convert errors to structured responses | YES |

### Non-Responsibilities

The MCP Runtime does NOT:
- Implement engineering rules
- Validate artifacts directly
- Manage state directly
- Enforce phases directly
- Execute shell commands
- Call CLI internally

---

## Architecture

### Position in Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                      AI Assistant (Client)                          │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ MCP Protocol
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        MCP Runtime                                   │
│                       (Thin Adapter)                                 │
│                                                                     │
│  • Receives: tool requests via MCP                                   │
│  • Calls: Workspace Engine methods                                   │
│  • Returns: structured JSON-RPC responses                             │
│  • NO business logic                                                 │
│  • NO shell commands                                                 │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Engine Interface
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Workspace Engine                               │
│                        (State Owner)                                  │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        KDSE Methodology                               │
│                        (Engineering Rules)                            │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Tool Definitions

### Tool Registry

```go
// Tool represents an MCP tool
type Tool struct {
    Name        string
    Description string
    InputSchema ToolInputSchema
    Handler     ToolHandler
}

// ToolHandler is the function signature for tool handlers
type ToolHandler func(ctx context.Context, input map[string]interface{}) (*ToolResult, error)
```

### Tool List

```go
var tools = []Tool{
    {
        Name:        "kdse_init",
        Description: "Initialize KDSE runtime in the current workspace",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]Property{
                "template": {
                    Type:        "string",
                    Description: "Template to use for initialization",
                    Default:     "default",
                },
                "runtime_type": {
                    Type:        "string",
                    Description: "Runtime type (cli or mcp)",
                    Default:     "mcp",
                },
            },
        },
        Handler: handleInit,
    },
    {
        Name:        "kdse_verify",
        Description: "Verify KDSE runtime state",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]Property{},
        },
        Handler: handleVerify,
    },
    {
        Name:        "kdse_phase",
        Description: "Get current phase or advance to next phase",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]Property{
                "action": {
                    Type:        "string",
                    Description: "Action to perform: 'show' or 'advance'",
                    Enum:        []string{"show", "advance"},
                },
            },
            Required: []string{"action"},
        },
        Handler: handlePhase,
    },
    {
        Name:        "kdse_artifacts",
        Description: "List or verify artifacts for a phase",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]Property{
                "phase": {
                    Type:        "string",
                    Description: "Phase to list artifacts for",
                    Enum:        []string{"knowledge", "architecture", "implementation", "verification", "reports"},
                },
                "verify": {
                    Type:        "boolean",
                    Description: "Whether to verify artifacts",
                    Default:     false,
                },
            },
            Required: []string{"phase"},
        },
        Handler: handleArtifacts,
    },
    {
        Name:        "kdse_report",
        Description: "Generate a report",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]Property{
                "type": {
                    Type:        "string",
                    Description: "Type of report to generate",
                    Enum:        []string{"phase", "verification", "summary", "progress"},
                },
            },
            Required: []string{"type"},
        },
        Handler: handleReport,
    },
    {
        Name:        "kdse_knowledge",
        Description: "Collect or promote knowledge artifacts",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]Property{
                "action": {
                    Type:        "string",
                    Description: "Action: 'collect', 'promote', or 'list'",
                    Enum:        []string{"collect", "promote", "list"},
                },
                "path": {
                    Type:        "string",
                    Description: "Path to knowledge artifact (for promote)",
                },
            },
            Required: []string{"action"},
        },
        Handler: handleKnowledge,
    },
}
```

---

## Tool Handlers

### Init Handler

```go
func handleInit(ctx context.Context, input map[string]interface{}) (*ToolResult, error) {
    // Parse input
    opts := workspace.InitOptions{
        Path:     getWorkingDir(),
        Type:     parseRuntimeType(input["runtime_type"]),
        Template: input["template"].(string),
    }
    
    // Call Workspace Engine (NO local logic)
    engine := workspace.NewEngine()
    ws, err := engine.InitializeWorkspace(ctx, opts)
    if err != nil {
        return formatError(err)
    }
    
    // Return structured response
    return &ToolResult{
        Success: true,
        Data: map[string]interface{}{
            "path":           ws.Path,
            "runtime_type":   ws.Runtime.Type,
            "runtime_version": ws.Runtime.Version,
            "phase":          ws.State.CurrentPhase,
        },
        Message: "KDSE runtime initialized successfully",
    }, nil
}
```

### Verify Handler

```go
func handleVerify(ctx context.Context, input map[string]interface{}) (*ToolResult, error) {
    // Call Workspace Engine
    engine := workspace.NewEngine()
    result, err := engine.VerifyWorkspace(ctx)
    if err != nil {
        return formatError(err)
    }
    
    if !result.Valid {
        return &ToolResult{
            Success: false,
            Data: map[string]interface{}{
                "valid":   false,
                "errors":  formatErrors(result.Errors),
                "phase":   result.Phase,
            },
            Message: "Workspace verification failed",
        }, nil
    }
    
    return &ToolResult{
        Success: true,
        Data: map[string]interface{}{
            "valid":           true,
            "phase":           result.Phase,
            "runtime_type":    result.RuntimeInfo.Type,
            "runtime_version": result.RuntimeInfo.Version,
            "timestamp":       result.Timestamp,
        },
        Message: "Workspace verified successfully",
    }, nil
}
```

### Phase Handler

```go
func handlePhase(ctx context.Context, input map[string]interface{}) (*ToolResult, error) {
    engine := workspace.NewEngine()
    action := input["action"].(string)
    
    switch action {
    case "show":
        phase, err := engine.GetPhase(ctx)
        if err != nil {
            return formatError(err)
        }
        return &ToolResult{
            Success: true,
            Data: map[string]interface{}{
                "name":        phase.Name,
                "description": phase.Description,
                "artifacts":   getRequiredArtifacts(phase.Name),
            },
        }, nil
        
    case "advance":
        transition, err := engine.AdvancePhase(ctx, phase.Name)
        if err != nil {
            return formatError(err)
        }
        return &ToolResult{
            Success: true,
            Data: map[string]interface{}{
                "from":     transition.From,
                "to":       transition.To,
                "timestamp": transition.Timestamp,
            },
            Message: fmt.Sprintf("Phase advanced from %s to %s", transition.From, transition.To),
        }, nil
    }
    
    return nil, ErrInvalidAction
}
```

### Artifacts Handler

```go
func handleArtifacts(ctx context.Context, input map[string]interface{}) (*ToolResult, error) {
    engine := workspace.NewEngine()
    phase := input["phase"].(string)
    verify := input["verify"].(bool)
    
    if verify {
        result, err := engine.VerifyArtifacts(ctx, parsePhase(phase))
        if err != nil {
            return formatError(err)
        }
        
        return &ToolResult{
            Success:   result.AllVerified,
            Data: map[string]interface{}{
                "verified":     result.Verified,
                "missing":      result.Missing,
                "invalid":      result.Invalid,
                "all_verified": result.AllVerified,
            },
        }, nil
    }
    
    artifacts, err := engine.GetArtifacts(ctx, parsePhase(phase))
    if err != nil {
        return formatError(err)
    }
    
    return &ToolResult{
        Success: true,
        Data: map[string]interface{}{
            "phase":     phase,
            "artifacts": artifacts,
        },
    }, nil
}
```

### Knowledge Handler

```go
func handleKnowledge(ctx context.Context, input map[string]interface{}) (*ToolResult, error) {
    engine := workspace.NewEngine()
    action := input["action"].(string)
    
    switch action {
    case "collect":
        collection, err := engine.CollectKnowledge(ctx)
        if err != nil {
            return formatError(err)
        }
        return &ToolResult{
            Success: true,
            Data: map[string]interface{}{
                "artifacts": collection.Artifacts,
                "count":     len(collection.Artifacts),
            },
        }, nil
        
    case "promote":
        path := input["path"].(string)
        err := engine.PromoteKnowledge(ctx, workspace.Artifact{Path: path})
        if err != nil {
            return formatError(err)
        }
        return &ToolResult{
            Success: true,
            Message: fmt.Sprintf("Knowledge artifact promoted: %s", path),
        }, nil
        
    case "list":
        collection, err := engine.CollectKnowledge(ctx)
        if err != nil {
            return formatError(err)
        }
        return &ToolResult{
            Success: true,
            Data: map[string]interface{}{
                "artifacts": collection.Artifacts,
            },
        }, nil
    }
    
    return nil, ErrInvalidAction
}
```

---

## Response Formatting

### Tool Result Structure

```go
// ToolResult is the result returned by a tool handler
type ToolResult struct {
    Success bool                   // Whether the tool succeeded
    Data    map[string]interface{} // Structured data
    Message string                 // Human-readable message
    Error   *ToolError            // Error details (if failed)
}

// ToolError contains error details
type ToolError struct {
    Code        string                 // Error code
    Message     string                 // Error message
    Details     map[string]interface{} // Additional details
    Remediation string                 // How to fix the error
}
```

### Error Response

```go
func formatError(err error) (*ToolResult, error) {
    var engineErr *workspace.EngineError
    if errors.As(err, &engineErr) {
        return &ToolResult{
            Success: false,
            Error: &ToolError{
                Code:        engineErr.Code,
                Message:     engineErr.Message,
                Details:     engineErr.Details,
                Remediation: engineErr.Remediation,
            },
        }, nil
    }
    
    return &ToolResult{
        Success: false,
        Error: &ToolError{
            Code:        "INTERNAL_ERROR",
            Message:     err.Error(),
            Remediation: "Check logs for details",
        },
    }, nil
}
```

### JSON-RPC Response Format

```json
{
  "jsonrpc": "2.0",
  "result": {
    "success": true,
    "data": {
      "phase": "knowledge",
      "runtime_version": "1.0.0"
    },
    "message": "Workspace verified successfully"
  },
  "id": 1
}
```

---

## MCP Server Implementation

### Server Structure

```go
// cmd/kdse-mcp/main.go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    
    "kdse/internal/workspace"
    "kdse/internal/mcp"
)

type Server struct {
    engine *workspace.Engine
    tools  []mcp.Tool
}

func NewServer() *Server {
    s := &Server{
        engine: workspace.NewEngine(),
        tools:  mcp.RegisterTools(),
    }
    return s
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
    var req mcp.JSONRPCRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, -32700, "Parse error")
        return
    }
    
    // Find tool handler
    handler := s.findTool(req.Method)
    if handler == nil {
        sendError(w, -32601, "Method not found")
        return
    }
    
    // Execute tool
    result, err := handler(r.Context(), req.Params)
    if err != nil {
        result, _ = formatError(err)
    }
    
    sendResponse(w, req.ID, result)
}
```

### Main Entry Point

```go
func main() {
    server := NewServer()
    
    http.HandleFunc("/mcp", server.HandleRequest)
    
    log.Println("KDSE MCP server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
```

---

## Package Structure

```
cmd/
└── kdse-mcp/
    ├── main.go              # Entry point
    ├── server/
    │   ├── server.go        # MCP server
    │   └── handler.go       # Request handler
    ├── tools/
    │   ├── tools.go         # Tool registry
    │   ├── init.go          # kdse_init tool
    │   ├── verify.go        # kdse_verify tool
    │   ├── phase.go         # kdse_phase tool
    │   ├── artifacts.go     # kdse_artifacts tool
    │   ├── report.go        # kdse_report tool
    │   └── knowledge.go     # kdse_knowledge tool
    ├── formatter/
    │   ├── formatter.go     # Response formatting
    │   └── error.go         # Error formatting
    └── types/
        └── types.go         # MCP types
```

---

## Dependency Rules

### Allowed Dependencies

```go
// cmd/kdse-mcp/main.go
package main

import (
    // OK: Standard library
    "context"
    "encoding/json"
    "log"
    "net/http"
    
    // OK: Internal workspace engine
    "kdse/internal/workspace"
    
    // OK: Internal types
    "kdse/internal/types"
    
    // OK: Internal methodology interfaces
    "kdse/internal/methodology/lifecycle"
)
```

### Forbidden Dependencies

```go
// cmd/kdse-mcp/main.go - DO NOT DO THIS
import (
    // FORBIDDEN: Direct CLI invocation
    "os/exec"  // NO - never exec kdse CLI
    
    // FORBIDDEN: Shell commands
    "os"  // Be careful with os package
    
    // FORBIDDEN: Direct state management
    "kdse/internal/state"    // NO
    
    // FORBIDDEN: Direct runtime logic
    "kdse/internal/runtime"  // NO
    
    // FORBIDDEN: Direct bootstrap logic
    "kdse/internal/bootstrap" // NO (use workspace engine)
)
```

---

## Verification Gate

### Pre-Tool Verification

```go
func withVerification(tool Tool, handler ToolHandler) ToolHandler {
    return func(ctx context.Context, input map[string]interface{}) (*ToolResult, error) {
        // Skip verification for init tool
        if tool.Name == "kdse_init" {
            return handler(ctx, input)
        }
        
        // Verify workspace before tool
        engine := workspace.NewEngine()
        result, err := engine.VerifyWorkspace(ctx)
        if err != nil {
            return formatError(err)
        }
        
        if !result.Valid {
            return &ToolResult{
                Success: false,
                Error: &ToolError{
                    Code:        "WORKSPACE_NOT_VERIFIED",
                    Message:     "Workspace verification failed",
                    Details:     formatErrors(result.Errors),
                    Remediation: "Run kdse_init to initialize workspace",
                },
            }, nil
        }
        
        // Add workspace info to context
        ctx = context.WithValue(ctx, WorkspaceKey, result)
        
        return handler(ctx, input)
    }
}
```

---

## Testing

### MCP Test Strategy

```go
// cmd/kdse-mcp/tools/tools_test.go
package tools

func TestVerifyTool(t *testing.T) {
    // Create temp workspace
    dir := t.TempDir()
    initWorkspace(t, dir)
    
    // Call verify tool
    handler := withVerification(tools["kdse_verify"], handleVerify)
    result, err := handler(context.Background(), map[string]interface{}{})
    
    // Verify results
    assert.NoError(t, err)
    assert.True(t, result.Success)
    assert.Equal(t, "knowledge", result.Data["phase"])
}
```

---

## Error Codes

### MCP Error Codes

| Code | Description | Remediation |
|------|-------------|-------------|
| E001 | Workspace not found | Use kdse_init tool |
| E002 | Verification failed | Check .kdse directory |
| E003 | Invalid phase | Use kdse_phase with 'show' |
| E004 | Invalid transition | Complete current phase first |
| E005 | Report generation failed | Check workspace state |
| E006 | Unknown tool | Check tool registry |
| E007 | Invalid parameters | Check tool schema |

---

## Document Relationships

```
MCP_RUNTIME.md
    │
    ├── Defines: MCP adapter responsibilities, tools, formatting
    │
    ├── Referenced By:
    │   ├── All MCP implementations
    │   └── Documentation
    │
    ├── References:
    │   ├── RUNTIME_ARCHITECTURE.md
    │   ├── WORKSPACE_ENGINE.md
    │   └── PRINCIPLES.md (I-001: Thin Adapters)
    │
    └── Related Documents:
        └── CLI_RUNTIME.md
```

---

*This document is normative. The MCP Runtime MUST be a thin adapter with no business logic.*
