# ADR-002: One Methodology, Multiple Runtimes

**Decision ID:** ADR-002  
**Type:** Architecture Decision  
**Status:** Accepted  
**Date:** 2026-07-17  
**Author:** KDSE Architecture Team

---

## Context

We need to support multiple execution environments (CLI, MCP, IDE integrations) while maintaining a single, consistent engineering methodology.

### Problem

- Multiple execution environments needed (CLI, MCP, etc.)
- Each environment should behave consistently
- Methodology must not be tied to implementation
- No duplication of engineering rules

### Evidence

Current architecture mixes:
- Runtime implementation with engineering rules
- CLI behavior with methodology
- MCP tools with knowledge processing

---

## Decision

**There is exactly ONE KDSE methodology, implemented by multiple thin-adapter runtimes.**

### Implication

| Aspect | Implication |
|--------|-------------|
| Methodology | Single implementation in methodology/ package |
| Runtimes | CLI, MCP are thin adapters that call methodology |
| State | Workspace Engine owns state, not runtimes |
| Behavior | All runtimes produce identical results |

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       KDSE Methodology                        │
│                      (Single Implementation)                  │
│                                                                     │
│  internal/methodology/                                          │
│  ├── lifecycle/     - Phase definitions                        │
│  ├── phases/       - Phase logic                              │
│  ├── authority/    - Authority rules                          │
│  ├── verification/ - Verification logic                       │
│  └── knowledge/    - Knowledge management                     │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Implements
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Workspace Engine                         │
│                      (State Owner)                           │
│                                                                     │
│  internal/workspace/engine/                                   │
│  ├── InitializeWorkspace()                                    │
│  ├── VerifyWorkspace()                                        │
│  ├── AdvancePhase()                                          │
│  └── ...                                                      │
└─────────────────────────────────────────────────────────────┘
                    ┌───────────────────┬───────────────────┐
                    │                   │                   │
                    ▼                   ▼                   ▼
            ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
            │ CLI Runtime │     │ MCP Runtime │     │ IDE Runtime │
            │(Thin Adapter│     │(Thin Adapter│     │(Future)     │
            └─────────────┘     └─────────────┘     └─────────────┘
```

---

## Consequences

### Positive

- Single source of methodology truth
- Consistent behavior across runtimes
- Easy to add new runtimes
- Clear separation of concerns

### Negative

- Abstraction overhead
- Interface maintenance
- More packages to maintain

### Neutral

- Thin adapters for each runtime
- Workspace Engine as intermediary

---

## Dependency Rules

### Forbidden Dependencies

```go
// ❌ FORBIDDEN: Methodology imports runtime
// internal/methodology/lifecycle/lifecycle.go
import "kdse/internal/runtime"  // NO!

// ❌ FORBIDDEN: Methodology imports CLI
import "kdse/cmd/kdse"  // NO!

// ❌ FORBIDDEN: Methodology imports MCP
import "kdse/cmd/kdse-mcp"  // NO!
```

### Allowed Dependencies

```go
// ✅ ALLOWED: Runtime imports methodology
// cmd/kdse/main.go
import (
    "kdse/internal/workspace/engine"
    "kdse/internal/methodology/lifecycle"
)

// ✅ ALLOWED: Workspace Engine imports methodology
// internal/workspace/engine/engine.go
import (
    "kdse/internal/methodology/lifecycle"
    "kdse/internal/methodology/authority"
)
```

### Dependency Graph

```
cmd/kdse ─────┐
cmd/kdse-mcp ─┼──→ internal/workspace/engine ──→ internal/methodology
                 ↑
                 │
internal/workspace/loader, validator, state
```

---

## Runtime Interface

### Interface Definition

```go
// Runtime defines the interface for all KDSE runtimes
type Runtime interface {
    Initialize(ctx context.Context, opts InitOptions) (*Workspace, error)
    Verify(ctx context.Context) (*VerificationResult, error)
    Load(ctx context.Context, path string) (*Workspace, error)
    GetPhase(ctx context.Context) (*Phase, error)
    AdvancePhase(ctx context.Context, target Phase) (*Transition, error)
    GenerateReport(ctx context.Context, reportType ReportType) (*Report, error)
}
```

### CLI Implementation

```go
// cmd/kdse/commands/runtime.go
type CLIRuntime struct {
    engine *engine.Engine
}

func (r *CLIRuntime) Initialize(ctx context.Context, opts InitOptions) (*Workspace, error) {
    // CLI just calls engine
    return r.engine.InitializeWorkspace(ctx, opts)
}

func (r *CLIRuntime) Verify(ctx context.Context) (*VerificationResult, error) {
    // CLI just calls engine
    return r.engine.VerifyWorkspace(ctx)
}
```

### MCP Implementation

```go
// cmd/kdse-mcp/tools/runtime.go
type MCPRuntime struct {
    engine *engine.Engine
}

func (r *MCPRuntime) Initialize(ctx context.Context, opts InitOptions) (*Workspace, error) {
    // MCP just calls engine
    return r.engine.InitializeWorkspace(ctx, opts)
}

func (r *MCPRuntime) Verify(ctx context.Context) (*VerificationResult, error) {
    // MCP just calls engine
    return r.engine.VerifyWorkspace(ctx)
}
```

---

## Equivalence Testing

### Test Requirement

Both CLI and MCP MUST produce identical workspaces:

```go
func TestRuntimesEquivalent(t *testing.T) {
    // Initialize with CLI
    cliDir := t.TempDir()
    cli := &CLIRuntime{engine: engine.New()}
    cliWs, _ := cli.Initialize(context.Background(), InitOptions{Path: cliDir})
    
    // Initialize with MCP
    mcpDir := t.TempDir()
    mcp := &MCPRuntime{engine: engine.New()}
    mcpWs, _ := mcp.Initialize(context.Background(), InitOptions{Path: mcpDir})
    
    // Compare
    assert.Equal(t, cliWs.State.Phase, mcpWs.State.Phase)
    assert.Equal(t, cliWs.State.Artifacts, mcpWs.State.Artifacts)
    assert.Equal(t, len(cliWs.State.History), len(mcpWs.State.History))
}
```

---

## Alternatives Considered

### Alternative 1: Runtime-Specific Methodologies

**Rejected because:**
- Duplication of logic
- Inconsistent behavior
- Hard to maintain
- Violates P-003 (One Methodology)

### Alternative 2: Single Monolithic Runtime

**Rejected because:**
- CLI overhead for MCP
- MCP overhead for CLI
- No flexibility
- Doesn't support different environments

### Alternative 3: Multiple Independent Implementations

**Rejected because:**
- Duplication
- Inconsistency
- Maintenance burden
- Violates P-003 (One Methodology)

---

## Related Decisions

| Decision | Relationship |
|----------|--------------|
| ADR-001 | Runtime is the Authority |
| ADR-003 | Workspace Engine |

---

## Related Principles

| Principle | Reference |
|-----------|-----------|
| P-003 | One Methodology |
| P-004 | Multiple Runtimes |
| P-005 | Runtime Independence |
| I-001 | Thin Adapters |

---

## Change Log

| Date | Change |
|------|--------|
| 2026-07-17 | Initial decision |

---

*This decision is normative. There must be exactly one methodology implementation, with multiple thin-adapter runtimes.*
