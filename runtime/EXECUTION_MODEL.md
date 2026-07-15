# KDSE Runtime Execution Model

**Document Version:** 2.0  
**Type:** Informative Reference Implementation  
**Effective Date:** 2026-07-15

---

## Purpose

This document describes the state-based orchestration engine of the KDSE Runtime. The Runtime is the operational component that orchestrates KDSE sessions, consuming the Standard and producing actionable engineering guidance.

**Key Difference from v1.0**: This execution model is **state-based**, not linear. Each execute cycle evaluates the current state and decides what action to take next, rather than following a fixed sequence.

---

## State-Based Orchestration Principles

### The Seven-Step Cycle

Each execute cycle performs these steps in order:

```
┌─────────────────────────────────────────────────────────────────┐
│                    EXECUTE CYCLE                                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. RESOLVE WORKSPACE                                            │
│     ↓ Determine repository, project, or temporary workspace       │
│                                                                  │
│  2. EVALUATE CURRENT STATE                                        │
│     ↓ Assess what phase we're in, what's blocked                 │
│                                                                  │
│  3. EVALUATE CONFIDENCE                                          │
│     ↓ Calculate foundation, repository, evidence confidence       │
│                                                                  │
│  4. EVALUATE MISSING EVIDENCE                                    │
│     ↓ Determine what evidence is required for current phase      │
│                                                                  │
│  5. DECIDE NEXT PHASE                                            │
│     ↓ Based on state, confidence, evidence - what's next?        │
│                                                                  │
│  6. EXECUTE ONLY THAT PHASE                                      │
│     ↓ Take action only on the decided phase                       │
│                                                                  │
│  7. RE-EVALUATE                                                   │
│     ↓ After execution, re-assess everything                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### State-Based vs Linear

| Aspect | Linear Model | State-Based Model |
|--------|--------------|-------------------|
| Path | Fixed sequence | Dynamic decision |
| Phase Selection | Determined by previous | Determined by state |
| Confidence | Calculated once | Re-evaluated each cycle |
| Evidence | Tracked manually | Auto-evaluated |
| Implementation Gate | Manual check | Automatic threshold |

---

## Workspace Resolution Hierarchy

The orchestrator supports a three-level workspace hierarchy:

```
Repository
    ↓
Project Folder
    ↓
Temporary Workspace
```

### Workspace Types

1. **Repository**: A git repository or directory containing source code
2. **Project**: A subdirectory within repository representing a specific project
3. **Temporary**: A workspace created under `./temp/.kdse/` for isolated work

### Temporary Workspace Rules

- Temporary workspaces are created under `./temp/.kdse/<project>/`
- **Never hardcode** `/app` or `/workspace` paths
- All filesystem paths come from the **Workspace Resolver**
- When a project is later created, `.kdse` is migrated automatically

---

## Orchestration Phases

```
┌─────────────────────────────────────────────────────────────────┐
│                     PHASE HIERARCHY                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────┐                                                    │
│  │  Idle   │ ← No active session                                │
│  └────┬────┘                                                    │
│       │ Initialize                                              │
│       ▼                                                         │
│  ┌───────────┐                                                  │
│  │  Resolve  │ ← Resolve workspace hierarchy                      │
│  └─────┬─────┘                                                  │
│        │ Workspace resolved                                     │
│        ▼                                                        │
│  ┌───────────┐                                                  │
│  │  Assess   │ ← Evaluate current repository state               │
│  └─────┬─────┘                                                  │
│        │ Assessment complete                                     │
│        ▼                                                        │
│  ┌─────────────┐                                                │
│  │ Foundation  │ ← Verify/establish Foundation documents         │
│  └──────┬──────┘                                                │
│         │ Foundation threshold met (REQUIRED GATE)               │
│         ▼                                                       │
│  ┌───────────┐                                                  │
│  │  Collect  │ ← Gather evidence for next phases                 │
│  └─────┬─────┘                                                  │
│        │ Evidence sufficient                                    │
│        ▼                                                        │
│  ┌───────────┐                                                  │
│  │  Analyze  │ ← Analyze collected evidence                      │
│  └─────┬─────┘                                                  │
│        │ Analysis complete                                       │
│        ▼                                                        │
│  ┌───────────┐                                                  │
│  │  Design   │ ← Design solution based on analysis              │
│  └─────┬─────┘                                                  │
│        │ Design complete                                         │
│        ▼                                                        │
│  ┌─────────────┐ ← Only if Foundation threshold met             │
│  │ Implement   │                                                  │
│  └─────┬───────┘                                                  │
│        │ Implementation complete                                 │
│        ▼                                                        │
│  ┌───────────┐                                                  │
│  │  Verify   │ ← Verify implementation results                  │
│  └─────┬─────┘                                                  │
│        │ Verification passes                                     │
│        ▼                                                        │
│  ┌───────────┐                                                  │
│  │ Complete  │ ← Session complete                                │
│  └───────────┘                                                  │
│                                                                  │
│         OR                                                       │
│                                                                  │
│  ┌───────────┐                                                  │
│  │ Blocked   │ ← Missing evidence/confidence threshold           │
│  └───────────┘                                                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Confidence Evaluation

### Confidence Dimensions

| Dimension | Weight | Description |
|-----------|--------|-------------|
| Foundation | 50% | Foundation document completeness |
| Repository | 30% | Repository structure and artifacts |
| Evidence | 20% | Evidence collected for current phase |

### Foundation Threshold Gate

**CRITICAL**: Implementation is **forbidden** until the Foundation threshold is met.

```
Implementation Gate:
┌─────────────────────────────────────────────────────────────┐
│                                                              │
│   Foundation Confidence ≥ Threshold?                          │
│         │                                                    │
│         ├── YES → Implementation ALLOWED                     │
│         │                                                    │
│         └── NO  → Implementation BLOCKED                     │
│                      Return to Foundation phase              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

The default Foundation threshold is **0.7 (70%)**, but this can be configured.

---

## Evidence Evaluation

### Per-Phase Evidence Requirements

Each phase has specific evidence requirements:

| Phase | Required Evidence | Critical |
|-------|-------------------|----------|
| Resolve | Repository manifest | Yes |
| Assess | README, SPEC | Yes |
| Foundation | Foundation docs (6 files) | Yes |
| Collect | Evidence directory | Yes |
| Analyze | Collected evidence | Yes |
| Design | Analysis results, Architecture | Yes |
| Implement | Design spec, Context handoff | Yes |
| Verify | Implementation, Tests | Yes |

### Evidence Completeness

Evidence completeness is calculated as:
```
Completeness = Evidence Present / Evidence Required
```

---

## Decision Logic

The orchestrator uses the following decision logic:

```
decideNextPhase():
    if blocked_by_missing_evidence:
        return STAY (cannot progress)
    
    if current_phase == Implement and not foundation_threshold_met:
        return Foundation (blocked gate)
    
    if current_phase_complete:
        return next_phase_in_sequence
    
    return current_phase (continue working)
```

---

## Session Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| foundation_threshold | 0.7 | Min Foundation confidence for Implementation |
| evidence_threshold | 0.6 | Min evidence completeness to proceed |
| max_cycles | 100 | Maximum orchestration cycles |
| temp_workspace_base | "temp" | Base directory for temp workspaces |

---

## Workspace Resolver API

All filesystem paths are obtained through the Workspace Resolver:

```go
// Resolve workspace hierarchy
workspace := resolver.ResolveWorkspace("/path/to/start")

// Create temporary workspace
tempWorkspace := resolver.ResolveTemporaryWorkspace("project-name")

// Migrate .kdse to project
resolver.MigrateToProject(tempKDSEPath, projectPath)

// Get subdirectory path
subPath := resolver.ResolveSubPath(workspacePath, "foundation")
```

**Rule**: Never hardcode paths like `/app` or `/workspace`. Always use the resolver.

---

## Progress Measurement

Progress is measured through:

| Metric | Description |
|--------|-------------|
| Confidence Score | Weighted composite of dimension scores |
| Foundation Score | Specific Foundation document coverage |
| Evidence Completeness | Evidence gathered vs required |
| Phase Progress | Cycles spent in each phase |

---

## Session Completion

A session completes when:

1. **Complete Phase Reached**: All phases verified
2. **Max Cycles**: Reached configured maximum
3. **Operator Closes**: Human ends session
4. **No Progress**: Repeated cycles without state change

---

## Key Differences from Linear Model

| Feature | Linear v1 | State-Based v2 |
|---------|-----------|-----------------|
| Phase Selection | Sequential | Dynamic based on state |
| Confidence Check | Manual | Automatic each cycle |
| Implementation Gate | Per-decision | Automatic threshold |
| Evidence Tracking | Separate process | Integrated in cycle |
| Path Resolution | Potentially hardcoded | Always through resolver |

---

## Implementation Notes

### For Developers

1. **Never hardcode paths** - Use `WorkspaceResolver` for all filesystem operations
2. **Check confidence before implement** - The engine blocks Implementation until threshold met
3. **Track evidence per phase** - Each phase has specific requirements
4. **Re-evaluate in each cycle** - Don't assume state persists between cycles

### For Operators

1. **Monitor confidence** - Watch Foundation score during session
2. **Provide evidence** - Help collect missing evidence when blocked
3. **Understand the gate** - Implementation requires Foundation at threshold

---

## References

- [SESSION_PROTOCOL.md](SESSION_PROTOCOL.md) - Session lifecycle details
- [REPORT_SPEC.md](REPORT_SPEC.md) - Runtime Report structure
- [ARTIFACT_VERIFICATION.md](ARTIFACT_VERIFICATION.md) - Artifact verification
- [docs/foundation/](../docs/foundation/) - Foundation documents
- [docs/audit/](../docs/audit/) - Audit standards

---

*This document describes the state-based execution model (v2.0). It replaces the linear model from v1.0.*
