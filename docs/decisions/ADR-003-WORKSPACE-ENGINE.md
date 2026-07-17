# ADR-003: Workspace Engine

**Decision ID:** ADR-003  
**Type:** Architecture Decision  
**Status:** Accepted  
**Date:** 2026-07-17  
**Author:** KDSE Architecture Team

---

## Context

We need a single component that owns all engineering state and enforces the KDSE methodology.

### Problem

Current architecture has:
- State distributed across components
- No single owner of project state
- Inconsistent state management
- No enforcement of methodology

### Evidence

| Current State | Problem |
|--------------|---------|
| `internal/state/` | Fragmented state management |
| `internal/runtime/` | Mixed with business logic |
| `internal/bootstrap/` | Separate from state |
| `internal/workspace/` | Unclear ownership |

---

## Decision

**The Workspace Engine is the ONLY owner of project state. All other components MUST go through the Workspace Engine.**

### Implication

| Component | Role |
|-----------|------|
| Workspace Engine | Owns all state, enforces methodology |
| CLI Runtime | Thin adapter, calls engine |
| MCP Runtime | Thin adapter, calls engine |
| Methodology | Defines rules, called by engine |

### Ownership

```
Workspace Engine owns:
├── .kdse/ directory structure
├── runtime.yaml
├── workspace.yaml
├── phase.yaml
├── session.yaml
├── metadata.yaml
├── All phase artifacts
└── All reports
```

### Not owned by Workspace Engine:

```
├── Source code (project files)
├── Dependencies (go.mod, etc.)
├── Git history
└── External files
```

---

## Responsibilities

### Workspace Engine

| Responsibility | Description |
|----------------|-------------|
| Create Runtime | Initialize .kdse directory |
| Load Runtime | Load runtime configuration |
| Verify Runtime | Verify runtime integrity |
| Load Phase | Load current phase state |
| Persist Phase | Save phase state |
| Validate Workspace | Validate workspace structure |
| Manage Metadata | Store/retrieve metadata |
| Manage Sessions | Track sessions |
| Generate Reports | Produce reports |

### Not Workspace Engine Responsibilities

| Responsibility | Owner |
|----------------|-------|
| Parse commands | CLI Runtime |
| Handle MCP protocol | MCP Runtime |
| Format output | CLI/MCP Runtime |
| Define methodology | Methodology package |
| Validate artifacts | Methodology/Authority |

---

## Interface

### Main Interface

```go
// Engine is the main interface for the Workspace Engine
type Engine interface {
    // Workspace lifecycle
    InitializeWorkspace(ctx context.Context, opts InitOptions) (*Workspace, error)
    VerifyWorkspace(ctx context.Context) (*VerificationResult, error)
    LoadWorkspace(ctx context.Context) (*Workspace, error)
    
    // Phase management
    GetPhase(ctx context.Context) (*Phase, error)
    AdvancePhase(ctx context.Context, target Phase) (*Transition, error)
    
    // Reporting
    GenerateReport(ctx context.Context, opts ReportOptions) (*Report, error)
    
    // Knowledge
    CollectKnowledge(ctx context.Context) (*KnowledgeCollection, error)
}
```

---

## Consequences

### Positive

- Single source of truth for state
- Clear ownership
- Consistent state management
- Easy to enforce methodology
- Simple verification

### Negative

- Single point of failure
- More complex component
- Interface maintenance

### Neutral

- All operations go through engine
- Thin adapters for runtimes

---

## Package Structure

```
internal/
├── workspace/
│   ├── engine/
│   │   ├── engine.go       # Main interface
│   │   ├── verify.go       # Verification logic
│   │   ├── init.go         # Initialization
│   │   ├── phase.go        # Phase management
│   │   ├── artifact.go     # Artifact management
│   │   └── report.go       # Report generation
│   ├── loader/
│   │   ├── loader.go       # Workspace loading
│   │   └── config.go       # Config loading
│   ├── validator/
│   │   ├── validator.go    # General validation
│   │   └── phase.go       # Phase validation
│   └── state/
│       ├── state.go       # State management
│       └── persistence.go # State persistence
```

---

## Verification Gate

### Every Operation

```go
func (e *Engine) AnyOperation(ctx context.Context, ...) error {
    // MUST verify workspace first
    result, err := e.VerifyWorkspace(ctx)
    if err != nil {
        return err
    }
    
    if !result.Valid {
        return ErrWorkspaceInvalid
    }
    
    // Proceed with operation
    return e.doOperation(ctx, ...)
}
```

---

## Alternatives Considered

### Alternative 1: Distributed State

**Rejected because:**
- No single source of truth
- Consistency problems
- Hard to verify
- State conflicts

### Alternative 2: Database State

**Rejected because:**
- Adds dependency
- Not portable
- More complex
- Violates simplicity

### Alternative 3: Runtime Owns State

**Rejected because:**
- Ties state to runtime
- Hard to test
- Doesn't fit thin adapter model
- Violates separation of concerns

---

## Related Decisions

| Decision | Relationship |
|----------|--------------|
| ADR-001 | Runtime is the Authority |
| ADR-002 | One Methodology, Multiple Runtimes |

---

## Related Principles

| Principle | Reference |
|-----------|-----------|
| P-002 | Runtime is the Authority |
| R-002 | State Ownership |
| D-003 | Single Responsibility |

---

## Change Log

| Date | Change |
|------|--------|
| 2026-07-17 | Initial decision |

---

*This decision is normative. The Workspace Engine is the sole owner of project state.*
