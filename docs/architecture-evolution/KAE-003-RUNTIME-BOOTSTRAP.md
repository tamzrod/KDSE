# KAE-003: Runtime Bootstrap

**Evolution ID:** KAE-003  
**Type:** Architecture Evolution  
**Status:** Planned  
**Author:** KDSE Architecture Team  
**Date:** 2026-07-17

---

## Summary

Implements verification-first runtime bootstrap that uses templates, ensures atomic initialization, and rolls back on failure.

---

## Context

Runtime initialization must be evidence-based, template-driven, and atomic. No partial initialization is acceptable.

---

## Problem Statement

### Current State

- Initialization is hardcoded
- No template support
- Partial initialization possible
- No rollback on failure

### Desired State

- Templates define structure
- Atomic operations
- Rollback on failure
- Verification before completion

---

## Implementation

### Bootstrap Flow

```
1. Download Runtime Template (if remote)
       ↓
2. Extract Runtime Files
       ↓
3. Create .kdse Directory
       ↓
4. Generate Metadata
       ↓
5. Initialize Phase
       ↓
6. Verify Runtime ← VERIFICATION GATE
       ↓ (only if verification passes)
7. Return Success
```

### Atomic Initialization

```go
func (e *Engine) InitializeWorkspace(ctx context.Context, opts InitOptions) (*Workspace, error) {
    // Create staging directory
    tmpDir, err := os.MkdirTemp("", "kdse-init-*")
    if err != nil {
        return nil, err
    }
    defer os.RemoveAll(tmpDir) // Cleanup on any failure
    
    // Stage all operations in tmpDir
    if err := e.stageRuntime(tmpDir, opts); err != nil {
        return nil, err
    }
    
    // Verify staged runtime
    if err := e.verifyStagedRuntime(tmpDir); err != nil {
        return nil, err
    }
    
    // Atomic move to final location
    if err := atomicMove(tmpDir, opts.Path); err != nil {
        return nil, err
    }
    
    // Now .kdse exists with complete state
    return e.LoadWorkspace(ctx)
}
```

### Rollback Strategy

```go
func (e *Engine) rollback(dir string) {
    // Remove all created files
    os.RemoveAll(dir)
    
    // Log rollback for debugging
    log.Printf("Rollback completed for %s", dir)
}
```

---

## Template Structure

### Runtime Template

```
templates/
└── runtime/
    ├── default/
    │   ├── runtime.yaml        # Template
    │   ├── workspace.yaml     # Template
    │   ├── methodology.yaml   # Template
    │   ├── phase.yaml         # Template
    │   ├── session.yaml       # Template
    │   ├── metadata.yaml      # Template
    │   └── directories/
    │       ├── knowledge/
    │       ├── architecture/
    │       ├── implementation/
    │       ├── verification/
    │       └── reports/
    └── v2/
        └── ...
```

### Template Variables

```yaml
# runtime.yaml template
runtime:
  type: {{.RuntimeType}}
  version: {{.RuntimeVersion}}
  commit: {{.TemplateCommit}}

template:
  version: {{.TemplateVersion}}
  commit: {{.TemplateCommit}}

workspace:
  version: {{.WorkspaceVersion}}
  root: {{.WorkspaceRoot}}

session:
  id: {{.SessionID}}
  created: {{.CreatedAt}}

phase:
  current: initialization
  previous: none
```

---

## Verification Requirements

### Verification Gate

Initialization MUST pass verification before returning success:

```go
func (e *Engine) InitializeWorkspace(...) (*Workspace, error) {
    // ... initialization steps ...
    
    // Critical: Verify before returning
    result, err := e.VerifyWorkspace(ctx)
    if err != nil || !result.Valid {
        e.rollback(tmpDir)
        return nil, ErrVerificationFailed
    }
    
    return ws, nil
}
```

### Verification Checklist

| Check | Description | On Failure |
|-------|-------------|------------|
| Directory exists | .kdse/ created | Rollback |
| Files exist | All required files | Rollback |
| Files valid | YAML syntax valid | Rollback |
| Phase valid | Phase is "initialization" | Rollback |
| Metadata valid | All metadata fields | Rollback |

---

## Error Handling

### Error Types

```go
var (
    ErrTemplateNotFound     = errors.New("runtime template not found")
    ErrTemplateInvalid      = errors.New("runtime template invalid")
    ErrInitializationFailed = errors.New("runtime initialization failed")
    ErrVerificationFailed   = errors.New("runtime verification failed")
    ErrRollbackFailed        = errors.New("rollback failed")
)
```

### Error Response

```json
{
  "error": {
    "code": "INITIALIZATION_FAILED",
    "message": "Runtime initialization failed",
    "details": {
      "stage": "verification",
      "reason": "runtime.yaml not found"
    },
    "remediation": "Run kdse init again or check template"
  }
}
```

---

## Related Documents

| Document | Relationship |
|----------|--------------|
| [KAE-001](KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md) | Parent evolution |
| [KAE-004](KAE-004-EVIDENCE-FIRST-RUNTIME.md) | Related evolution |

---

## Status

**Status:** Planned  
**Depends on:** KAE-001, KAE-002  
**Estimated Effort:** Medium

---

*This document describes the runtime bootstrap evolution.*
