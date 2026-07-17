# KDSE Architectural Principles

**Document Version:** 1.0  
**Type:** Normative  
**Effective Date:** 2026-07-17

---

## Purpose

This document defines the core architectural principles that govern KDSE. All decisions, implementations, and evolutions must adhere to these principles.

---

## Core Principles

### P-001: Evidence First

**Statement:** No engineering work may proceed without evidence.

**Rationale:** Knowledge without evidence is merely opinion. Engineering decisions must be grounded in verifiable facts, measurements, and documentation.

**Implications:**
- Every claim must be backed by evidence
- Every decision must reference source documents
- Every artifact must be verifiable
- The filesystem is the authoritative source of truth

**Evidence Hierarchy:**
1. Primary artifacts (source code, configs, tests)
2. Secondary artifacts (documentation, reports)
3. Tertiary artifacts (comments, commit messages)
4. Assertions (claims without supporting artifacts)

### P-002: Runtime is the Authority

**Statement:** The KDSE Runtime (.kdse) is the authoritative engineering state.

**Rationale:** Without an authoritative runtime, there is no single source of truth. Multiple interpretations lead to inconsistency and non-compliance.

**Implications:**
- The .kdse directory defines what is and is not a KDSE project
- No engineering claim is valid without a verified runtime
- The runtime persists all engineering state
- The runtime is mandatory, not optional

### P-003: One Methodology

**Statement:** There is only one engineering methodology.

**Rationale:** Multiple methodologies create conflict, confusion, and inconsistency. KDSE provides one comprehensive approach.

**Implications:**
- All KDSE projects use the same methodology
- The methodology is enforced, not suggested
- Deviations require formal change control
- The methodology spans all engineering phases

### P-004: Multiple Runtimes

**Statement:** CLI, MCP, and future runtimes are execution adapters.

**Rationale:** The methodology must remain independent of execution contexts. Different environments require different interfaces.

**Implications:**
- CLI, MCP, and IDE integrations are adapters
- All adapters produce equivalent results
- The methodology never changes based on runtime
- Runtime implementations may vary; methodology cannot

### P-005: Runtime Independence

**Statement:** Methodology must never depend on runtime implementation.

**Rationale:** Coupling methodology to runtime creates fragility. The methodology must outlast any specific runtime implementation.

**Implications:**
- Methodology packages do not import runtime packages
- Runtime packages may import methodology packages
- The interface is runtime → methodology, never methodology → runtime
- Abstraction boundaries are strictly enforced

### P-006: Filesystem is Evidence

**Statement:** The filesystem is authoritative. If .kdse does not exist, the project is NOT a KDSE project.

**Rationale:** In distributed systems, the filesystem provides the only consistent ground truth. Memory, network, and cache are ephemeral.

**Implications:**
- .kdse must exist to declare a project KDSE-compliant
- Filesystem artifacts take precedence over claims
- No verbal or documented claim replaces filesystem evidence
- Verification must always check filesystem state

---

## Design Principles

### D-001: Explicit Over Implicit

**Statement:** Be explicit about all assumptions, dependencies, and requirements.

**Application:**
- Document all decisions with rationale
- Declare all dependencies explicitly
- State all requirements in requirements files
- Never assume implicit behavior

### D-002: Fail Fast and Loud

**Statement:** Fail immediately and clearly when invariants are violated.

**Application:**
- Verification happens before operations
- Errors include context and remediation
- No silent failures or graceful degradation
- Partial states are never acceptable

### D-003: Single Responsibility

**Statement:** Each component has one clear responsibility.

**Application:**
- Workspace Engine owns project state
- Runtimes own execution only
- Methodology owns engineering rules
- Each package has one owner

### D-004: Dependency Inversion

**Statement:** High-level modules must not depend on low-level modules.

**Application:**
- Methodology defines interfaces
- Runtimes implement interfaces
- Dependencies point toward abstraction
- No circular dependencies

### D-005: Least Surprise

**Statement:** Behavior must match reasonable expectations.

**Application:**
- Consistent behavior across runtimes
- Predictable phase transitions
- Clear error messages
- Intuitive workflows

---

## Enforcement Principles

### E-001: Verification Before Action

**Statement:** Always verify runtime state before any engineering action.

**Sequence:**
1. Verify Runtime Exists
2. Load Runtime
3. Load Current Phase
4. Validate Required Artifacts
5. Continue

If verification fails, STOP.

### E-002: Atomic Operations

**Statement:** Operations succeed completely or fail completely with rollback.

**Application:**
- Initialization either completes or rolls back
- Phase transitions are atomic
- No partial states persist
- Transactions are used where applicable

### E-003: Immutable Evidence

**Statement:** Once evidence is recorded, it cannot be altered.

**Application:**
- Audit trails are append-only
- Historical data is preserved
- Corrections create new evidence, not modifications
- Reports reference immutable sources

---

## Implementation Principles

### I-001: Thin Adapters

**Statement:** Runtimes are thin adapters, not thick implementations.

**CLI Responsibilities:**
- Parse commands
- Call Workspace Engine
- Display output
- Nothing else

**MCP Responsibilities:**
- Receive tool requests
- Call Workspace Engine
- Return structured responses
- Nothing else

### I-002: No Magic

**Statement:** No hidden behavior, automatic actions, or implicit transformations.

**Application:**
- Explicit configuration required
- No auto-detection without evidence
- No assumptions about project structure
- No default behaviors that bypass methodology

### I-003: Template-Based Bootstrap

**Statement:** Initialization uses templates, never hardcoded values.

**Application:**
- Runtime templates define structure
- Configuration comes from templates
- No embedded methodology in binaries
- Templates enable customization

---

## Relationship Principles

### R-001: Runtime → Methodology

**Statement:** Runtimes consume methodology; methodology never consumes runtime.

**Valid Dependency Graph:**
```
Runtime → Workspace Engine → Methodology
     ↓           ↓              ↓
   CLI/MCP    State Mgmt    Engineering
```

**Invalid Dependency Graph:**
```
Methodology → Runtime
```

### R-002: State Ownership

**Statement:** Only one component owns state at any time.

**Ownership Rules:**
- Workspace Engine owns .kdse state
- Runtime owns session state
- Methodology owns engineering rules
- No shared mutable state

### R-003: Interface Segregation

**Statement:** Clients depend only on methods they use.

**Application:**
- Runtimes use only Runtime interface
- Workspace Engine uses only Methodology interface
- Small, focused interfaces
- No fat interfaces

---

## Evidence Principles

### V-001: Artifacts are Primary Evidence

**Statement:** Filesystem artifacts take precedence over all other evidence.

**Evidence Priority:**
1. `.kdse/` directory contents
2. Source code files
3. Configuration files
4. Test files
5. Documentation files
6. Runtime reports
7. Verbal/written claims

### V-002: Traceability Required

**Statement:** Every artifact must be traceable to a requirement or decision.

**Traceability Chain:**
```
Requirement → Decision → Implementation → Verification → Artifact
```

### V-003: Verification Required

**Statement:** Every artifact must be verified before promotion.

**Verification Sequence:**
1. Create artifact
2. Verify completeness
3. Verify correctness
4. Record verification evidence
5. Promote to next level

---

## Phase Principles

### PH-001: Sequential Phases

**Statement:** Engineering proceeds through phases in strict sequence.

**Prohibited Actions:**
- Skipping phases
- Regressing phases
- Parallel phases
- Concurrent phase work

### PH-002: Phase Completion Criteria

**Statement:** Each phase has explicit completion criteria.

**Completion Requirements:**
- All required artifacts exist
- All artifacts are verified
- Phase metadata is recorded
- Transition is authorized

### PH-003: Phase Isolation

**Statement:** Each phase produces isolated, self-contained artifacts.

**Phase Characteristics:**
- Clear inputs and outputs
- No cross-phase dependencies
- Complete within phase boundaries
- Evidence preserved for audit

---

## Security Principles

### S-001: No Secret Propagation

**Statement:** Secrets must never leave their context.

**Rules:**
- Credentials not stored in .kdse
- API keys not in source code
- Tokens not in logs
- Secrets in environment only

### S-002: Minimal Trust

**Statement:** Trust nothing by default.

**Trust Levels:**
- Filesystem artifacts: trusted
- Runtime claims: verified
- AI claims: verified
- External systems: untrusted until verified

### S-003: Audit Everything

**Statement:** All significant actions must be auditable.

**Audit Requirements:**
- Log all state changes
- Record all decisions
- Track all phase transitions
- Preserve all verification evidence

---

## Anti-Patterns

### AP-001: Claim Without Evidence

**Statement:** Never claim KDSE initialization without verified .kdse existence.

**Anti-Pattern:**
```
AI: "KDSE initialized successfully."
Reality: .kdse does not exist.
```

**Correct Pattern:**
```
1. Create .kdse directory
2. Verify .kdse exists
3. Initialize runtime.yaml
4. Verify runtime.yaml valid
5. Report success
```

### AP-002: Runtime Coupling

**Statement:** Never import runtime packages in methodology code.

**Anti-Pattern:**
```go
// methodology/lifecycle/lifecycle.go
import "kdse/internal/runtime"  // VIOLATION
```

**Correct Pattern:**
```go
// methodology/lifecycle/lifecycle.go
// No runtime imports

// runtime/cli/cli.go
import "kdse/internal/methodology"  // VALID
```

### AP-003: Implicit State

**Statement:** Never rely on implicit state that is not persisted.

**Anti-Pattern:**
```
Memory: "Current phase is Architecture"
Reality: Phase not recorded in .kdse
```

**Correct Pattern:**
```
1. Set current phase in .kdse/phase.yaml
2. Persist immediately
3. Verify persisted correctly
4. Continue with verified state
```

### AP-004: Hardcoded Values

**Statement:** Never hardcode configuration that should come from templates.

**Anti-Pattern:**
```go
const DefaultPhase = "knowledge"
```

**Correct Pattern:**
```go
// Templates define defaults
// Runtime loads template values
// No hardcoded phase names
```

---

## Principle Hierarchy

```
┌─────────────────────────────────────────────────────────────┐
│                    Core Principles (P)                       │
│                                                              │
│  Evidence First → Runtime Authority → One Methodology       │
│                                                              │
│  Multiple Runtimes → Runtime Independence → Filesystem      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Design Principles (D)                       │
│                                                              │
│  Explicit Over Implicit                                      │
│  Fail Fast and Loud                                          │
│  Single Responsibility                                       │
│  Dependency Inversion                                        │
│  Least Surprise                                              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Enforcement Principles (E)                   │
│                                                              │
│  Verification Before Action                                  │
│  Atomic Operations                                           │
│  Immutable Evidence                                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│               Implementation Principles (I)                 │
│                                                              │
│  Thin Adapters                                               │
│  No Magic                                                    │
│  Template-Based Bootstrap                                    │
└─────────────────────────────────────────────────────────────┘
```

---

## Document Relationships

```
PRINCIPLES.md
    │
    ├── Defines: Core, Design, Enforcement, Implementation, Relationship, Evidence, Phase, Security
    │
    ├── Referenced By:
    │   ├── RUNTIME_ARCHITECTURE.md
    │   ├── METHODOLOGY.md
    │   ├── WORKSPACE_ENGINE.md
    │   ├── CLI_RUNTIME.md
    │   ├── MCP_RUNTIME.md
    │   ├── ADR-001-RUNTIME-IS-THE-AUTHORITY.md
    │   ├── ADR-002-ONE-METHODOLOGY-MULTIPLE-RUNTIMES.md
    │   └── ADR-003-WORKSPACE-ENGINE.md
    │
    └── Evolution Documents:
        ├── KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md
        └── KAE-002-CLI-MCP-RUNTIME-SEPARATION.md
```

---

*This document is normative. All KDSE implementations must adhere to these principles.*
