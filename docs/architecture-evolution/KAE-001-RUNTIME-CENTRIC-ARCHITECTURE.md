# KAE-001: Runtime-Centric Architecture

**Evolution ID:** KAE-001  
**Type:** Architecture Evolution  
**Status:** In Progress  
**Author:** KDSE Architecture Team  
**Date:** 2026-07-17

---

## Summary

Transform KDSE into a Runtime-Centric Engineering Platform where the KDSE Runtime (`.kdse/`) becomes the authoritative foundation of every engineering project.

---

## Problem Statement

### Current Problem

The current architecture allows an AI to claim KDSE initialization without a verified runtime:

```
AI: "KDSE initialized successfully."
Reality: .kdse does not exist.
```

This violates P-001 (Evidence First) and P-006 (Filesystem is Evidence).

### Root Causes

1. **No verification gate**: Engineering can proceed without verifying runtime existence
2. **Implicit state**: Phase and workspace state are not enforced
3. **Runtime coupling**: Methodology and runtime are not separated
4. **Thick adapters**: CLI and MCP contain business logic
5. **No single source of truth**: State is fragmented across components

---

## Solution

### Architecture Transformation

#### Before (Current)

```
Participant → CLI/MCP (with business logic) → Fragmented State
                     ↓
              Multiple sources of truth
```

#### After (Target)

```
Participant → CLI/MCP (thin adapter) → Workspace Engine → Methodology
                                                ↓
                                          .kdse/
                                    (Single Source of Truth)
```

### Key Changes

| Aspect | Before | After |
|--------|--------|-------|
| Source of Truth | Fragmented | .kdse/ directory |
| Runtime Role | Mixed with logic | Thin adapter only |
| State Ownership | Distributed | Workspace Engine only |
| Phase Enforcement | None | Enforced by engine |
| Initialization | Hardcoded | Template-based |

---

## Implementation Plan

### Phase 1: Architecture Documentation (COMPLETE)
- [x] Create docs/architecture/ documents
- [x] Create docs/architecture-evolution/ documents
- [x] Create docs/decisions/ ADRs

### Phase 2: Repository Refactor (PENDING)
- [ ] Create cmd/kdse/ (CLI)
- [ ] Create cmd/kdse-mcp/ (MCP)
- [ ] Create internal/methodology/
- [ ] Create internal/workspace/
- [ ] Create internal/runtime/
- [ ] Create internal/bootstrap/

**Target Structure:**
```
cmd/
├── kdse/          # CLI runtime
│   ├── main.go
│   └── commands/
└── kdse-mcp/     # MCP runtime
    ├── main.go
    └── tools/

internal/
├── methodology/
│   ├── lifecycle/
│   ├── phases/
│   ├── authority/
│   ├── verification/
│   ├── knowledge/
│   └── architecture/
├── workspace/
│   ├── engine/
│   ├── loader/
│   ├── validator/
│   └── state/
├── runtime/
│   ├── runtime.go
│   ├── cli/
│   └── mcp/
├── bootstrap/
└── templates/
```

### Phase 3: Workspace Engine (PENDING)
- [ ] Create engine interface
- [ ] Implement state management
- [ ] Implement verification
- [ ] Implement phase management
- [ ] Implement session management
- [ ] Implement report generation

### Phase 4: Runtime Abstraction (PENDING)
- [ ] Define Runtime interface
- [ ] Implement CLI runtime
- [ ] Implement MCP runtime
- [ ] Ensure identical behavior

### Phase 5: Runtime Bootstrap (PENDING)
- [ ] Remove hardcoded initialization
- [ ] Implement template download
- [ ] Implement atomic initialization
- [ ] Implement rollback on failure
- [ ] Implement verification-first init

### Phase 6: KDSE Runtime (PENDING)
- [ ] Define mandatory .kdse structure
- [ ] Implement runtime.yaml schema
- [ ] Implement workspace.yaml schema
- [ ] Implement phase.yaml schema
- [ ] Enforce runtime existence

### Phase 7: Workspace Verification (PENDING)
- [ ] Implement verification gate
- [ ] Verify before every operation
- [ ] Fail fast on verification failure
- [ ] Return clear error messages

### Phase 8: CLI Runtime (PENDING)
- [ ] Strip business logic from CLI
- [ ] Use only Workspace Engine
- [ ] Implement thin adapter pattern

### Phase 9: MCP Runtime (PENDING)
- [ ] Strip business logic from MCP
- [ ] Use only Workspace Engine
- [ ] Implement thin adapter pattern
- [ ] Never call CLI internally

### Phase 10: Templates (PENDING)
- [ ] Move templates outside binary
- [ ] Create template registry
- [ ] Implement template download
- [ ] Support custom templates

### Phase 11: Runtime Metadata (PENDING)
- [ ] Implement metadata.yaml
- [ ] Track runtime type
- [ ] Track runtime version
- [ ] Track template version
- [ ] Track session info

### Phase 12: Engineering Rules (PENDING)
- [ ] Enforce phase sequence
- [ ] Prohibit phase skipping
- [ ] Implement completion criteria
- [ ] Implement evidence requirements

### Phase 13: Testing (PENDING)
- [ ] Test CLI produces identical workspace
- [ ] Test MCP produces identical workspace
- [ ] Test workspace verification
- [ ] Test phase transitions
- [ ] Test artifact validation

### Phase 14: Migration (PENDING)
- [ ] Create migration guide
- [ ] Document breaking changes
- [ ] Provide migration tooling
- [ ] Test migration paths

### Phase 15: Repository Documentation (PENDING)
- [ ] Update README
- [ ] Update Developer Guide
- [ ] Update Architecture Guide
- [ ] Update Runtime Guide
- [ ] Update Bootstrap Guide

---

## Verification

### Verification Checklist

| Criterion | Description | Method |
|-----------|-------------|--------|
| V-001 | Exactly one methodology implementation | Package inspection |
| V-002 | CLI is thin adapter | No business logic imports |
| V-003 | MCP is thin adapter | No business logic imports |
| V-004 | Workspace Engine owns state | Package inspection |
| V-005 | .kdse is mandatory | Cannot proceed without |
| V-006 | Verification before action | Gate implementation |
| V-007 | Filesystem is source of truth | Verification test |
| V-008 | No claims without runtime | Test AI claims |
| V-009 | Both runtimes equivalent | Integration test |
| V-010 | Evolution documented | ADR completeness |

### Verification Tests

```go
func TestRuntimeIsMandatory(t *testing.T) {
    // Create temp directory without .kdse
    dir := t.TempDir()
    
    // Try to verify (should fail)
    engine := workspace.NewEngine()
    _, err := engine.VerifyWorkspace(context.Background())
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "runtime not found")
}

func TestCLIProducesIdenticalWorkspace(t *testing.T) {
    // Initialize with CLI
    cliDir := t.TempDir()
    runCLI("init", cliDir)
    
    // Initialize with MCP
    mcpDir := t.TempDir()
    runMCP("kdse_init", mcpDir)
    
    // Compare (excluding runtime-specific metadata)
    cliState := loadWorkspaceState(cliDir)
    mcpState := loadWorkspaceState(mcpDir)
    
    assert.Equal(t, cliState.Phase, mcpState.Phase)
    assert.Equal(t, cliState.Artifacts, mcpState.Artifacts)
}
```

---

## Migration Guide

### Breaking Changes

| Change | Impact | Mitigation |
|--------|--------|------------|
| CLI behavior change | Runtimes are thin adapters | Update documentation |
| Initialization change | Template-based | Migration guide |
| State change | Centralized in engine | Migration tooling |
| API changes | Engine interface | Deprecation period |

### Migration Steps

1. **Backup existing project**
   ```bash
   cp -r project project.backup
   ```

2. **Update KDSE binary**
   ```bash
   # Download new version with architecture changes
   ```

3. **Run migration tool**
   ```bash
   kdse migrate --from=v1 --to=v2
   ```

4. **Verify migration**
   ```bash
   kdse verify
   ```

5. **Test functionality**
   ```bash
   kdse phase show
   kdse artifacts knowledge
   ```

---

## Success Criteria

| # | Criterion | Verification |
|---|-----------|--------------|
| 1 | Exactly one implementation of methodology | Package inspection |
| 2 | CLI and MCP are thin adapters | Code review |
| 3 | Workspace Engine owns all engineering state | Architecture review |
| 4 | KDSE Runtime (.kdse) is mandatory | Integration test |
| 5 | Engineering cannot proceed without verified runtime | Integration test |
| 6 | Runtime is evidence-based | Test verification |
| 7 | Filesystem artifacts are source of truth | Architecture review |
| 8 | AI cannot claim KDSE init without verified runtime | Test AI claims |
| 9 | Both runtimes produce equivalent workspaces | Integration test |
| 10 | Architecture evolution is fully documented | Document review |

---

## Related Documents

| Document | Relationship |
|----------|--------------|
| [PRINCIPLES.md](../architecture/PRINCIPLES.md) | Defines core principles |
| [CURRENT_ARCHITECTURE.md](../architecture/CURRENT_ARCHITECTURE.md) | Current state analysis |
| [RUNTIME_ARCHITECTURE.md](../architecture/RUNTIME_ARCHITECTURE.md) | Target architecture |
| [ADR-001](../decisions/ADR-001-RUNTIME-IS-THE-AUTHORITY.md) | Runtime authority decision |

---

## Change Log

| Date | Change | Status |
|------|--------|--------|
| 2026-07-17 | Initial document | In Progress |

---

*This document describes the Runtime-Centric Architecture evolution for KDSE.*
