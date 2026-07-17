# ADR-001: Runtime is the Authority

**Decision ID:** ADR-001  
**Type:** Architecture Decision  
**Status:** Accepted  
**Date:** 2026-07-17  
**Author:** KDSE Architecture Team

---

## Context

We need to establish a single source of truth for KDSE project state. Without an authoritative runtime, multiple interpretations can lead to inconsistency and non-compliance.

### Problem

Current architecture allows:
- AI to claim KDSE initialization without verified runtime
- State fragmentation across components
- No enforcement of runtime existence
- Claims without evidence

### Evidence

```
AI: "KDSE initialized successfully."
Reality: .kdse does not exist.
```

This violates P-001 (Evidence First).

---

## Decision

**The KDSE Runtime (.kdse/) is the authoritative source of truth for all engineering state.**

### Implication

| Statement | Implication |
|-----------|-------------|
| .kdse/ exists | Project IS a KDSE project |
| .kdse/ missing | Project is NOT a KDSE project |
| Claims vs .kdse/ | .kdse/ takes precedence |

### Authority Hierarchy

```
1. .kdse/ directory (filesystem)
   └── runtime.yaml, phase.yaml, etc.
   
2. Source code, configs, tests
   
3. Documentation, comments

4. Verbal/written claims (LOWEST)
```

---

## Consequences

### Positive

- Single source of truth
- Enforceable verification
- Consistent state management
- Clear authority chain

### Negative

- Requires .kdse/ for all operations
- Migration for existing projects
- More rigid structure

### Neutral

- Template-based initialization
- Verification gate added

---

## Evidence Requirements

### For Runtime Authority

| Evidence | Description |
|----------|-------------|
| Directory exists | .kdse/ must exist |
| Files valid | All required files exist and valid |
| State consistent | No conflicting state |
| Phase valid | Phase is recognized |
| Metadata valid | All metadata fields present |

### For Claims

| Claim Type | Evidence Required |
|------------|-------------------|
| "KDSE project" | .kdse/ exists and verified |
| "Phase complete" | All phase artifacts verified |
| "Verification passed" | verification.yaml exists |
| "Initialization done" | All runtime files created and verified |

---

## Implementation

### Verification Sequence

```go
func (e *Engine) VerifyWorkspace(ctx context.Context) (*VerificationResult, error) {
    // 1. Check .kdse exists (PRIMARY EVIDENCE)
    if !e.exists(".kdse") {
        return &VerificationResult{
            Valid: false,
            Errors: []VerificationError{{
                Code:    "KDSE_MISSING",
                Message: ".kdse directory not found",
            }},
        }, ErrRuntimeMissing
    }
    
    // 2. Load runtime config
    runtime, err := e.loadRuntimeConfig()
    if err != nil {
        return nil, err
    }
    
    // 3. Verify phase
    phase, err := e.loadPhase()
    if err != nil {
        return nil, err
    }
    
    // 4. Verify required artifacts
    if err := e.verifyRequiredArtifacts(phase); err != nil {
        return nil, err
    }
    
    return &VerificationResult{
        Valid:      true,
        Phase:      phase,
        RuntimeInfo: runtime,
    }, nil
}
```

### Error Messages

```go
var (
    ErrRuntimeMissing = errors.New("runtime not found: .kdse directory missing")
    ErrRuntimeInvalid = errors.New("runtime configuration invalid")
)
```

---

## Alternatives Considered

### Alternative 1: Memory as Authority

**Rejected because:**
- Not persistent
- Not shareable between processes
- Can be lost on restart
- Violates P-006 (Filesystem is Evidence)

### Alternative 2: Database as Authority

**Rejected because:**
- Adds database dependency
- Not portable
- More complex setup
- Violates simplicity principle

### Alternative 3: Git as Authority

**Rejected because:**
- Not all projects use Git
- Requires Git initialization
- Adds external dependency
- Doesn't fit all use cases

---

## Related Decisions

| Decision | Relationship |
|----------|--------------|
| ADR-002 | One Methodology, Multiple Runtimes |
| ADR-003 | Workspace Engine |

---

## Related Principles

| Principle | Reference |
|-----------|-----------|
| P-001 | Evidence First |
| P-002 | Runtime is the Authority |
| P-006 | Filesystem is Evidence |

---

## Change Log

| Date | Change |
|------|--------|
| 2026-07-17 | Initial decision |

---

*This decision is normative. All implementations must treat the runtime as the authority.*
