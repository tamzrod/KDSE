# KAE-002: CLI/MCP Runtime Separation

**Evolution ID:** KAE-002  
**Type:** Architecture Evolution  
**Status:** Planned  
**Author:** KDSE Architecture Team  
**Date:** 2026-07-17

---

## Summary

Separates CLI and MCP into thin adapters that only translate between user interfaces and the Workspace Engine, with no business logic.

---

## Context

After implementing the Runtime-Centric Architecture (KAE-001), CLI and MCP runtimes must be refactored into thin adapters.

---

## Problem Statement

### Current State

CLI and MCP currently contain business logic:
- Direct artifact validation
- State management
- Phase enforcement
- Knowledge processing

### Desired State

CLI and MCP are thin adapters:
- Parse commands/requests
- Call Workspace Engine
- Format output/responses
- NO business logic

---

## Implementation

### Dependency Structure

```
cmd/kdse/          cmd/kdse-mcp/
    ↓                   ↓
workspace/engine ←←←←←
    ↓
methodology/lifecycle
```

### Forbidden Imports

```go
// cmd/kdse/commands/init.go - FORBIDDEN
package commands

import (
    // FORBIDDEN
    "kdse/internal/runtime"      // NO
    "kdse/internal/bootstrap"    // NO
    "kdse/internal/state"         // NO
    "kdse/internal/knowledge"    // NO
)

// ALLOWED
import (
    "kdse/internal/workspace/engine"  // YES
    "kdse/internal/types"             // YES
)
```

### Command Handler Pattern

```go
// Before (with business logic)
func handleInit(args []string) error {
    // Create directory
    os.MkdirAll(".kdse", 0755)
    
    // Create files
    createRuntimeConfig()
    createPhaseConfig()
    
    // Validate
    validateWorkspace()
    
    return nil
}

// After (thin adapter)
func handleInit(ctx context.Context, args []string) error {
    engine := workspace.NewEngine()
    
    opts := workspace.InitOptions{
        Path: getWorkingDir(),
    }
    
    ws, err := engine.InitializeWorkspace(ctx, opts)
    if err != nil {
        return fmt.Errorf("init failed: %w", err)
    }
    
    fmt.Printf("Initialized at %s\n", ws.Path)
    return nil
}
```

---

## Verification

### Test: No Business Logic in Runtimes

```bash
# Check for forbidden imports in CLI
grep -r "internal/runtime" cmd/kdse/
grep -r "internal/bootstrap" cmd/kdse/

# Should return no results
```

### Test: All Operations Go Through Engine

```bash
# Trace function calls
go test -trace cmd/kdse/

# All operations should call workspace.NewEngine()
```

---

## Migration Path

1. Create new thin adapter structure
2. Implement Workspace Engine
3. Update CLI commands to use engine
4. Update MCP tools to use engine
5. Remove old business logic
6. Verify identical behavior

---

## Related Documents

| Document | Relationship |
|----------|--------------|
| [CLI_RUNTIME.md](../architecture/CLI_RUNTIME.md) | CLI specification |
| [MCP_RUNTIME.md](../architecture/MCP_RUNTIME.md) | MCP specification |
| [KAE-001](KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md) | Parent evolution |

---

## Status

**Status:** Planned  
**Depends on:** KAE-001 (Runtime-Centric Architecture)  
**Estimated Effort:** Medium

---

*This document describes the CLI/MCP runtime separation evolution.*
