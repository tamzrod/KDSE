# Runtime-Centric Architecture

**Document Version:** 1.0  
**Type:** Normative  
**Effective Date:** 2026-07-17

---

## Purpose

This document defines the Runtime-Centric Architecture for KDSE. It establishes the runtime as the authoritative foundation of every engineering project.

---

## Architecture Overview

### High-Level Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                       KDSE Methodology                                │
│                        (Normative)                                    │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  • Engineering Rules                                         │   │
│  │  • Phase Definitions                                         │   │
│  │  • Artifact Specifications                                   │   │
│  │  • Validation Logic                                          │   │
│  │  • Evidence Requirements                                     │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ⚠️  MUST NOT depend on runtime implementations                    │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Implements
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Workspace Engine                                 │
│                       (State Owner)                                   │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  • Create Runtime                                            │   │
│  │  • Load Runtime                                              │   │
│  │  • Verify Runtime                                            │   │
│  │  • Load Phase                                                │   │
│  │  • Persist Phase                                             │   │
│  │  • Validate Workspace                                        │   │
│  │  • Manage Metadata                                           │   │
│  │  • Manage Sessions                                           │   │
│  │  • Generate Reports                                          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ✓  Owns ALL engineering state                                      │
│  ✓  Single source of truth for project state                       │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                    ┌───────────────┴───────────────┐
                    │                               │
                    ▼                               ▼
┌───────────────────────────────────────────────────────────────────┐
│                      CLI Runtime                                   │
│                       (Adapter)                                     │
│                                                                   │
│  ┌───────────────────────────────────────────────────────────┐   │
│  │  • Parse commands                                          │   │
│  │  • Call Workspace Engine                                   │   │
│  │  • Display output                                           │   │
│  │  NO BUSINESS LOGIC                                         │   │
│  └───────────────────────────────────────────────────────────┘   │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘
┌───────────────────────────────────────────────────────────────────┐
│                      MCP Runtime                                   │
│                       (Adapter)                                     │
│                                                                   │
│  ┌───────────────────────────────────────────────────────────┐   │
│  │  • Receive tool requests                                  │   │
│  │  • Call Workspace Engine                                   │   │
│  │  • Return structured responses                             │   │
│  │  NO BUSINESS LOGIC                                         │   │
│  └───────────────────────────────────────────────────────────┘   │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌───────────────────────────────────────────────────────────────────┐
│                    Filesystem (.kdse/)                             │
│                     (Evidence Store)                               │
│                                                                   │
│  runtime.yaml    - Runtime configuration                         │
│  workspace.yaml  - Workspace state                                │
│  methodology.yaml - Methodology reference                         │
│  phase.yaml      - Current phase state                            │
│  session.yaml    - Session state                                 │
│  knowledge/      - Knowledge artifacts                            │
│  architecture/   - Architecture artifacts                          │
│  implementation/ - Implementation artifacts                        │
│  verification/   - Verification artifacts                          │
│  reports/        - Generated reports                              │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘
```

---

## Core Principle: Runtime is the Authority

### Evidence Hierarchy

```
1. KDSE Runtime (.kdse/) exists
   │
   ├── Verifiable → Project IS a KDSE project
   │
   └── Missing → Project is NOT a KDSE project
                  (regardless of any claims)
```

### Verification Sequence

Before ANY engineering action:

```
┌─────────────────┐
│ Verify Runtime  │
│   Exists        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Load Runtime  │
│   Configuration │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Load Current    │
│     Phase       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Validate        │
│ Required        │
│   Artifacts     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   CONTINUE      │
│ (if all pass)   │
└─────────────────┘

If any step FAILS → STOP, Return Error
```

---

## Architecture Boundaries

### Layer Definitions

| Layer | Responsibility | Dependencies |
|-------|----------------|--------------|
| Methodology | Engineering rules, phase definitions | None (standalone) |
| Workspace Engine | State management, verification | Methodology |
| CLI Runtime | Command parsing, output display | Workspace Engine |
| MCP Runtime | Request handling, response formatting | Workspace Engine |
| Filesystem | Evidence persistence | None |

### Dependency Rules

```
VALID:
  CLI/MCP → Workspace Engine → Methodology
  Workspace Engine → Methodology
  Runtime → Workspace Engine

INVALID:
  Methodology → Workspace Engine
  Methodology → Runtime
  Workspace Engine → CLI/MCP
```

### Forbidden Patterns

```go
// ❌ INVALID: Methodology imports runtime
// internal/methodology/phase/phase.go
import "kdse/internal/runtime"  // FORBIDDEN

// ✅ VALID: Runtime imports methodology
// cmd/kdse/main.go
import "kdse/internal/methodology"
```

---

## Workspace Engine Specification

### Responsibilities

The Workspace Engine is the ONLY owner of project state.

| Responsibility | Description |
|----------------|-------------|
| Create Runtime | Initialize .kdse directory structure |
| Load Runtime | Load runtime configuration from filesystem |
| Verify Runtime | Verify runtime integrity and compliance |
| Load Phase | Load current phase state |
| Persist Phase | Save phase state to filesystem |
| Validate Workspace | Validate workspace structure and artifacts |
| Manage Metadata | Store and retrieve runtime metadata |
| Manage Sessions | Track session state and history |
| Generate Reports | Produce verification and progress reports |
| Collect Knowledge | Aggregate knowledge artifacts |

### Interface Definition

```go
// Runtime interface - implemented by CLI and MCP
type Runtime interface {
    // Workspace lifecycle
    InitializeWorkspace(ctx context.Context, opts InitOptions) (*Workspace, error)
    VerifyWorkspace(ctx context.Context) (*VerificationResult, error)
    LoadWorkspace(ctx context.Context) (*Workspace, error)
    
    // Phase management
    CurrentPhase(ctx context.Context) (Phase, error)
    AdvancePhase(ctx context.Context) (*PhaseTransition, error)
    
    // Reporting
    GenerateReport(ctx context.Context, reportType ReportType) (*Report, error)
    
    // Knowledge
    CollectKnowledge(ctx context.Context) (*KnowledgeCollection, error)
    
    // Architecture
    GenerateArchitecture(ctx context.Context) (*Architecture, error)
    
    // Verification
    VerifyArtifacts(ctx context.Context) (*VerificationResult, error)
}
```

### State Ownership

```
Workspace Engine owns:
├── .kdse/runtime.yaml      - Runtime configuration
├── .kdse/workspace.yaml    - Workspace state
├── .kdse/methodology.yaml  - Methodology reference
├── .kdse/phase.yaml        - Current phase
├── .kdse/session.yaml      - Session state
├── .kdse/metadata.yaml     - Runtime metadata
├── .kdse/knowledge/        - Knowledge artifacts
├── .kdse/architecture/     - Architecture artifacts
├── .kdse/implementation/    - Implementation artifacts
├── .kdse/verification/      - Verification artifacts
└── .kdse/reports/           - Generated reports
```

---

## Runtime Abstraction

### Runtime Interface

All runtimes (CLI, MCP, future) MUST implement the same interface.

```go
// Runtime is the interface for all KDSE runtime implementations
type Runtime interface {
    // Initialization
    Initialize(ctx context.Context, config RuntimeConfig) error
    Verify(ctx context.Context) (*VerificationResult, error)
    
    // Workspace operations
    Create(ctx context.Context, path string) (*Workspace, error)
    Load(ctx context.Context, path string) (*Workspace, error)
    Validate(ctx context.Context) (*ValidationResult, error)
    
    // Phase operations
    GetPhase(ctx context.Context) (*Phase, error)
    Advance(ctx context.Context, target Phase) (*Transition, error)
    
    // Reporting
    Report(ctx context.Context, reportType ReportType) (*Report, error)
}
```

### CLI Runtime

**Responsibilities:**
- Parse command-line arguments
- Call Workspace Engine methods
- Display formatted output
- Handle user interaction

**Constraints:**
- NO business logic
- NO direct filesystem operations (go through Workspace Engine)
- NO direct methodology implementation

### MCP Runtime

**Responsibilities:**
- Receive tool requests via MCP protocol
- Call Workspace Engine methods
- Return structured JSON responses
- Handle async operations

**Constraints:**
- NO business logic
- NO CLI invocation
- NO shell commands (use native APIs)
- NO direct filesystem operations (go through Workspace Engine)

---

## KDSE Runtime Directory

### Required Structure

When a project is initialized for KDSE, the following structure MUST exist:

```
.kdse/
├── runtime.yaml           # Runtime configuration (REQUIRED)
├── workspace.yaml         # Workspace state (REQUIRED)
├── methodology.yaml        # Methodology reference (REQUIRED)
├── phase.yaml              # Current phase (REQUIRED)
├── session.yaml             # Session state (REQUIRED)
├── metadata.yaml            # Runtime metadata (REQUIRED)
├── knowledge/               # Knowledge artifacts
│   ├── README.md
│   └── .gitkeep
├── architecture/           # Architecture artifacts
│   ├── README.md
│   └── .gitkeep
├── implementation/          # Implementation artifacts
│   ├── README.md
│   └── .gitkeep
├── verification/            # Verification artifacts
│   ├── README.md
│   └── .gitkeep
└── reports/                 # Generated reports
    ├── README.md
    └── .gitkeep
```

### Metadata Specification

```yaml
# runtime.yaml
runtime:
  type: cli  # or "mcp"
  version: 1.0.0
  commit: abc123

template:
  version: 2.0
  commit: def456

workspace:
  version: 1.0.0
  root: /path/to/project

session:
  id: session-uuid
  created: 2026-07-17T10:00:00Z

phase:
  current: initialization
  previous: none
```

---

## Engineering Workflow

### Phase Sequence

```
Blank Project
      │
      ▼
┌─────────────────┐
│   Initialize    │
│    Runtime      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Verify        │
│    Runtime      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│     Load        │
│  Methodology    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Knowledge     │◄──── Cannot skip
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Architecture   │◄──── Cannot skip
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Implementation  │◄──── Cannot skip
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Verification    │◄──── Cannot skip
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Reports      │◄──── Cannot skip
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Completed     │
└─────────────────┘
```

### Anti-Skipping Enforcement

```go
func (e *WorkspaceEngine) AdvancePhase(ctx context.Context, target Phase) (*Transition, error) {
    current := e.CurrentPhase(ctx)
    
    if !isValidTransition(current, target) {
        return nil, &InvalidTransitionError{
            Current:  current,
            Target:   target,
            Required: validNextPhase(current),
        }
    }
    
    // Verify all required artifacts for current phase exist
    if err := e.verifyCurrentPhaseArtifacts(ctx, current); err != nil {
        return nil, &IncompletePhaseError{
            Phase:    current,
            Missing:  err.Missing,
        }
    }
    
    return e.transitionTo(ctx, target)
}
```

---

## Evidence Requirements

### Verification Artifacts

For each phase, the following MUST exist:

| Phase | Required Artifacts |
|-------|-------------------|
| Knowledge | knowledge/*.md (at least 1) |
| Architecture | architecture/*.md (at least 1) |
| Implementation | implementation/*.md (at least 1) |
| Verification | verification/*.md (at least 1) |
| Reports | reports/*.md (at least 1) |

### Verification Record

```yaml
# verification/verification.yaml
verification:
  timestamp: 2026-07-17T10:00:00Z
  phase: architecture
  status: passed
  
artifacts:
  - path: architecture/ARCHITECTURE.md
    verified: true
    checksum: abc123
  
  - path: architecture/DECISIONS.md
    verified: true
    checksum: def456

evidence:
  runtime_verified: true
  runtime_version: 1.0.0
  workspace_valid: true
```

---

## Error Handling

### Verification Failure Response

```go
type VerificationError struct {
    Phase    Phase
    Failures []Failure
    Message  string
}

func (e *VerificationError) Error() string {
    return fmt.Sprintf("verification failed for phase %s: %s", e.Phase, e.Message)
}

// Verification errors MUST include remediation guidance
func (e *VerificationError) Remediation() string {
    return "Run 'kdse init' to initialize runtime, then retry."
}
```

### Error Response Format

```json
{
  "error": {
    "code": "VERIFICATION_FAILED",
    "message": "Runtime verification failed",
    "details": {
      "missing": [".kdse/runtime.yaml"],
      "invalid": [],
      "phase": "initialization"
    },
    "remediation": "Run 'kdse init' to initialize runtime"
  }
}
```

---

## Document Relationships

```
RUNTIME_ARCHITECTURE.md
    │
    ├── Defines: Architecture overview, boundaries, components
    │
    ├── References:
    │   ├── PRINCIPLES.md (normative)
    │   ├── WORKSPACE_ENGINE.md
    │   ├── METHODOLOGY.md
    │   ├── CLI_RUNTIME.md
    │   └── MCP_RUNTIME.md
    │
    └── Related Documents:
        ├── ADR-001-RUNTIME-IS-THE-AUTHORITY.md
        ├── ADR-002-ONE-METHODOLOGY-MULTIPLE-RUNTIMES.md
        ├── ADR-003-WORKSPACE-ENGINE.md
        └── KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md
```

---

*This document is normative. All implementations must adhere to this architecture.*
