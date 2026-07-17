# Current Architecture Analysis

**Document Version:** 1.0  
**Type:** Informative  
**Analysis Date:** 2026-07-17

---

## Purpose

This document describes the current KDSE architecture as of 2026-07-17. It identifies the problems that necessitate the Runtime-Centric Architecture evolution.

---

## Current Architecture Overview

### Directory Structure

```
kdse/
├── cmd/
│   ├── kdse/           # CLI implementation
│   └── mcp/            # MCP server implementation
├── internal/
│   ├── agreement/      # Agreement handling
│   ├── bootstrap/      # Bootstrap logic
│   ├── collect/        # Collection logic
│   ├── config/         # Configuration management
│   ├── constraints/    # Constraint checking
│   ├── context/        # Context management
│   ├── detection/      # Detection logic
│   ├── enforcer/       # Enforcement logic
│   ├── guard/          # Session guard
│   ├── knowledge/      # Knowledge management
│   ├── mcpclient/      # MCP client
│   ├── metrics/        # Metrics collection
│   ├── normalize/      # Normalization logic
│   ├── orchestration/  # Orchestration engine
│   ├── report/         # Report generation
│   ├── runtime/        # Runtime logic
│   ├── shttp/          # HTTP server
│   ├── someday/        # Future planning
│   ├── state/          # State management
│   ├── types/          # Type definitions
│   └── workspace/      # Workspace management
├── docs/               # Documentation
├── runtime/            # Runtime documentation and templates
└── templates/          # Project templates
```

### Current Component Responsibilities

#### Command Components

| Component | Responsibility |
|-----------|----------------|
| cmd/kdse | CLI entry point |
| cmd/mcp | MCP server entry point |

#### Internal Components

| Component | Responsibility |
|-----------|----------------|
| orchestration | Orchestration engine, confidence scoring, evidence collection |
| runtime | Runtime compliance, invariant checking |
| workspace | Workspace management and initialization |
| bootstrap | Bootstrap logic |
| state | State management |
| guard | Session guard |
| knowledge | Knowledge management |
| collect | Collection logic |
| report | Report generation |
| mcpclient | MCP client orchestration |
| config | Configuration management |

---

## Identified Problems

### P-001: Evidence Without Foundation

**Problem:** The architecture allows AI to claim KDSE initialization without verified runtime existence.

**Example:**
```
AI: "KDSE initialized successfully."
Reality: .kdse does not exist.
```

**Impact:** Violates P-001 (Evidence First) and P-006 (Filesystem is Evidence).

**Root Cause:** No enforcement of runtime verification before claims.

### P-002: Runtime Coupling

**Problem:** Methodology is not separated from runtime implementation.

**Current State:**
- `internal/runtime/` contains both runtime logic and engineering rules
- Engineering rules are entangled with execution logic
- No clear boundary between methodology and runtime

**Impact:** 
- Cannot add new runtimes without modifying methodology
- Testing requires runtime dependencies
- Evolution is constrained by implementation

**Root Cause:** No architectural separation between layers.

### P-003: Implicit Initialization

**Problem:** Initialization is hardcoded, not template-driven.

**Current State:**
```go
const DefaultPhase = "knowledge"
```

**Impact:**
- Cannot customize without code changes
- No version control for initialization
- Templates are not used for bootstrap

**Root Cause:** Initialization logic embedded in code.

### P-004: State Fragmentation

**Problem:** Engineering state is not centralized.

**Current State:**
- State distributed across multiple components
- No single owner of project state
- Inconsistent state management

**Impact:**
- Conflicting state updates
- No authoritative source of truth
- Difficult to verify compliance

**Root Cause:** No Workspace Engine component.

### P-005: Thick Adapters

**Problem:** CLI and MCP are not thin adapters.

**Current State:**
- `cmd/kdse` likely contains business logic
- `cmd/mcp` likely contains business logic
- Runtimes are not pure adapters

**Impact:**
- Duplicated logic between runtimes
- Inconsistent behavior
- Difficult to maintain

**Root Cause:** No Runtime interface abstraction.

### P-006: No Phase Enforcement

**Problem:** Phase transitions are not enforced.

**Current State:**
- Phases exist in documentation
- No enforcement of phase sequence
- AI can skip phases

**Impact:**
- Methodology violations go undetected
- No evidence of phase compliance
- Architecture drift

**Root Cause:** No phase enforcement in runtime.

### P-007: Non-Atomic Operations

**Problem:** Operations may leave partial state on failure.

**Current State:**
- Initialization may partially complete
- No rollback on failure
- Inconsistent project state possible

**Impact:**
- Corrupted project state
- Difficult recovery
- Non-deterministic behavior

**Root Cause:** No atomic operation guarantee.

### P-008: No Verification Gate

**Problem:** No verification before engineering actions.

**Current State:**
- AI proceeds without runtime verification
- No check for .kdse existence
- Claims without evidence

**Impact:**
- Actions against invalid state
- Corrupted project structure
- Methodology violations

**Root Cause:** No verification gate implementation.

---

## Current Data Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Participant (AI/Human)                       │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Command/Request
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                          CLI / MCP Runtime                          │
│  • Parses commands                                                   │
│  • Calls internal components                                         │
│  • Contains business logic                                           │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Calls
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       Internal Components                            │
│  • orchestration                                                     │
│  • runtime                                                           │
│  • workspace                                                         │
│  • state                                                             │
│  • guard                                                             │
│  • ... (many more)                                                   │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Reads/Writes
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                           Filesystem                                  │
│  • Project files                                                     │
│  • .kdse (maybe exists)                                              │
│  • Templates                                                          │
└─────────────────────────────────────────────────────────────────────┘
```

**Problems with Current Flow:**
1. CLI/MCP contain business logic (violates thin adapter principle)
2. No verification gate before operations
3. State is fragmented across components
4. No single source of truth

---

## Target Architecture

### Target Data Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Participant (AI/Human)                       │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Command/Request
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        CLI / MCP Runtime (Thin)                      │
│  • Parses commands                                                   │
│  • Calls Workspace Engine                                            │
│  • Displays output                                                   │
│  NO BUSINESS LOGIC                                                   │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Workspace Engine Interface
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Workspace Engine                              │
│  • Owns project state                                                │
│  • Verifies runtime                                                  │
│  • Enforces phase transitions                                        │
│  • Manages lifecycle                                                 │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Methodology Interface
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         KDSE Methodology                              │
│  • Engineering rules                                                 │
│  • Phase definitions                                                 │
│  • Artifact validation                                               │
│  NO RUNTIME DEPENDENCIES                                             │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ State Persistence
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      .kdse/ (Runtime Directory)                      │
│  • runtime.yaml                                                      │
│  • workspace.yaml                                                    │
│  • methodology.yaml                                                  │
│  • phase.yaml                                                        │
│  • session.yaml                                                      │
│  • knowledge/                                                        │
│  • architecture/                                                     │
│  • implementation/                                                   │
│  • verification/                                                     │
│  • reports/                                                          │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Migration Requirements

### From Current to Target

| Current State | Target State | Migration Step |
|--------------|--------------|----------------|
| Mixed business logic | Thin adapters | Extract logic to Workspace Engine |
| Fragmented state | Centralized state | Create Workspace Engine |
| No verification gate | Verification required | Add verification before operations |
| Hardcoded initialization | Template-based | Extract to templates |
| Implicit phases | Enforced phases | Add phase enforcement |
| Non-atomic operations | Atomic operations | Add transaction support |
| Runtime coupling | Runtime independence | Separate methodology packages |

### Breaking Changes

1. **CLI/MCP Behavior Change:** Runtimes will only be adapters
2. **Initialization Change:** Template-based, not hardcoded
3. **State Change:** Centralized in Workspace Engine
4. **Verification Change:** Required before all operations

### Backward Compatibility

- Existing .kdse directories will require migration
- Existing workflows may need adjustment
- API changes will require updates

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Migration complexity | High | High | Phased rollout, documentation |
| Breaking changes | High | Medium | Clear migration guide |
| Testing gaps | Medium | High | Comprehensive test suite |
| Performance regression | Low | Medium | Benchmarking |

---

## Document Relationships

```
CURRENT_ARCHITECTURE.md
    │
    ├── Describes: Current state, problems, target state
    │
    ├── References:
    │   ├── RUNTIME_ARCHITECTURE.md (target)
    │   ├── WORKSPACE_ENGINE.md
    │   ├── METHODOLOGY.md
    │   ├── CLI_RUNTIME.md
    │   └── MCP_RUNTIME.md
    │
    └── Related:
        ├── KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md
        └── ADR-001-RUNTIME-IS-THE-AUTHORITY.md
```

---

*This document is informative. It describes the current state for reference during architecture evolution.*
