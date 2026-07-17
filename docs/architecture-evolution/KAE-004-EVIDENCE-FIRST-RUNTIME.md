# KAE-004: Evidence-First Runtime

**Evolution ID:** KAE-004  
**Type:** Architecture Evolution  
**Status:** Planned  
**Author:** KDSE Architecture Team  
**Date:** 2026-07-17

---

## Summary

Implements evidence-first verification that ensures no engineering action proceeds without verified runtime state.

---

## Context

P-001 (Evidence First) requires that no engineering work proceeds without evidence. The runtime must enforce this.

---

## Problem Statement

### Current State

- AI can claim KDSE initialization without .kdse
- No verification before operations
- Claims are accepted without evidence
- No enforcement of runtime existence

### Desired State

```
AI: "I will initialize KDSE."
Runtime: [Creates .kdse]
Runtime: [Verifies .kdse exists]
Runtime: [Verifies files valid]
Runtime: [Returns success]
AI: "KDSE initialized." ← Now backed by evidence
```

---

## Implementation

### Verification Gate

Every operation MUST verify runtime before proceeding:

```go
func (e *Engine) verifyBeforeAction(ctx context.Context, action string) error {
    // Step 1: Check .kdse exists
    if !e.exists(".kdse") {
        return &EngineError{
            Code:        "RUNTIME_MISSING",
            Message:     "KDSE runtime not found",
            Remediation: "Run 'kdse init' to initialize",
        }
    }
    
    // Step 2: Load runtime config
    runtime, err := e.loadRuntimeConfig()
    if err != nil {
        return &EngineError{
            Code:        "RUNTIME_INVALID",
            Message:     err.Error(),
            Remediation: "Check .kdse/runtim.yaml",
        }
    }
    
    // Step 3: Load phase
    phase, err := e.loadPhase()
    if err != nil {
        return &EngineError{
            Code:        "PHASE_INVALID",
            Message:     err.Error(),
            Remediation: "Check .kdse/phase.yaml",
        }
    }
    
    // Step 4: Verify required artifacts
    if err := e.verifyRequiredArtifacts(phase); err != nil {
        return &EngineError{
            Code:        "ARTIFACTS_MISSING",
            Message:     err.Error(),
            Remediation: "Complete current phase artifacts",
        }
    }
    
    return nil
}
```

### Evidence Collection

```go
type Evidence struct {
    Type        string
    Path        string
    Checksum    string
    Timestamp   time.Time
    Verified    bool
    Verifier    string
}

// CollectEvidence gathers evidence for an artifact
func (e *Engine) CollectEvidence(artifact string) (*Evidence, error) {
    // Verify file exists
    if !e.exists(artifact) {
        return nil, ErrArtifactMissing
    }
    
    // Calculate checksum
    checksum, err := e.checksum(artifact)
    if err != nil {
        return nil, err
    }
    
    return &Evidence{
        Type:      "file",
        Path:      artifact,
        Checksum:  checksum,
        Timestamp: time.Now(),
        Verified:  true,
        Verifier:  "workspace-engine",
    }, nil
}
```

---

## Evidence Hierarchy

### Level 1: Primary Evidence (Required)

| Evidence | Description | Required |
|----------|-------------|----------|
| .kdse directory | Runtime directory | YES |
| .kdse/runtime.yaml | Runtime configuration | YES |
| .kdse/phase.yaml | Current phase | YES |

### Level 2: Secondary Evidence

| Evidence | Description | Required |
|----------|-------------|----------|
| Artifact files | Phase artifacts | YES |
| Phase history | Transition history | NO |

### Level 3: Tertiary Evidence

| Evidence | Description | Required |
|----------|-------------|----------|
| Logs | Operation logs | NO |
| Reports | Generated reports | NO |

---

## Verification Sequence

### Before Any Engineering Action

```
┌─────────────────────────┐
│ Verify .kdse exists     │
└───────────┬─────────────┘
            │
            ▼ (fail if missing)
┌─────────────────────────┐
│ Verify runtime.yaml     │
└───────────┬─────────────┘
            │
            ▼ (fail if invalid)
┌─────────────────────────┐
│ Verify phase.yaml       │
└───────────┬─────────────┘
            │
            ▼ (fail if invalid)
┌─────────────────────────┐
│ Verify required         │
│ artifacts for phase     │
└───────────┬─────────────┘
            │
            ▼ (fail if missing)
┌─────────────────────────┐
│ CONTINUE                │
└─────────────────────────┘
```

---

## Anti-Pattern Prevention

### Pattern: Claim Without Evidence

```go
// ❌ BEFORE: AI claims without evidence
func (ai *AI) InitializeKDSE() {
    // AI just claims
    fmt.Println("KDSE initialized successfully")
    // Reality: .kdse might not exist
}

// ✅ AFTER: Runtime ensures evidence
func (runtime *Runtime) Initialize() error {
    // Create .kdse
    if err := createRuntime(); err != nil {
        return err
    }
    
    // Verify evidence exists
    if !verifyEvidence() {
        return ErrEvidenceMissing
    }
    
    // Only now can we claim success
    return nil
}
```

### Pattern: Implicit Phase

```go
// ❌ BEFORE: Phase stored in memory
func (ai *AI) DoWork() {
    currentPhase = "knowledge" // In memory only
}

// ✅ AFTER: Phase in filesystem
func (engine *Engine) GetPhase() (Phase, error) {
    data, err := os.ReadFile(".kdse/phase.yaml")
    if err != nil {
        return "", err
    }
    return parsePhase(data)
}
```

---

## Evidence Recording

### Evidence File

```yaml
# .kdse/evidence.yaml
evidence:
  initialized: true
  initialized_at: 2026-07-17T10:00:00Z
  
  verified:
    - type: directory
      path: .kdse
      exists: true
      
    - type: file
      path: .kdse/runtime.yaml
      exists: true
      checksum: abc123
      
    - type: file
      path: .kdse/phase.yaml
      exists: true
      checksum: def456
      
  current_phase:
    name: initialization
    verified_at: 2026-07-17T10:00:01Z
```

---

## Related Documents

| Document | Relationship |
|----------|--------------|
| [PRINCIPLES.md](../architecture/PRINCIPLES.md) | P-001 Evidence First |
| [KAE-001](KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md) | Parent evolution |

---

## Status

**Status:** Planned  
**Depends on:** KAE-001, KAE-002, KAE-003  
**Estimated Effort:** Low

---

*This document describes the evidence-first runtime evolution.*
