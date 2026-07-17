# KDSE Methodology

**Document Version:** 1.0  
**Type:** Normative  
**Effective Date:** 2026-07-17

---

## Purpose

This document defines the KDSE Methodology—the engineering rules and phase definitions that govern all KDSE projects. The methodology is independent of runtime implementations.

---

## Methodology Definition

### What is KDSE Methodology?

The KDSE Methodology is the set of engineering rules that define:
- How engineering work proceeds
- What phases must be completed
- What artifacts must exist
- How artifacts must be validated
- What evidence is required

### What the Methodology Is NOT

The methodology does NOT include:
- Runtime implementations
- CLI commands
- MCP tool definitions
- User interfaces
- Platform-specific code

### Dependency Rule

```
Methodology packages MUST NEVER import runtime packages.
Runtime packages MAY import methodology packages.
```

---

## Phase Model

### Phase Definitions

| Phase | Description | Entry Criteria | Exit Criteria |
|-------|-------------|----------------|----------------|
| Initialization | Set up runtime and workspace | Project exists | .kdse verified |
| Knowledge | Gather and document requirements | Runtime verified | Knowledge artifacts verified |
| Architecture | Design system architecture | Knowledge verified | Architecture artifacts verified |
| Implementation | Build the system | Architecture verified | Implementation artifacts verified |
| Verification | Test and validate | Implementation verified | Verification artifacts verified |
| Reports | Document results | Verification verified | Reports generated |

### Phase Dependencies

```
┌──────────────┐
│ Initialization │ ← No dependencies (starts here)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Knowledge   │ ← Requires Initialization
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Architecture │ ← Requires Knowledge
└──────┬───────┘
       │
       ▼
┌──────────────┐
│Implementation│ ← Requires Architecture
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Verification │ ← Requires Implementation
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Reports    │ ← Requires Verification
└──────────────┘
```

### Phase Transition Rules

```go
// Valid transitions
var validTransitions = map[Phase][]Phase{
    PhaseInitialization: {PhaseKnowledge},
    PhaseKnowledge:      {PhaseArchitecture},
    PhaseArchitecture:   {PhaseImplementation},
    PhaseImplementation: {PhaseVerification},
    PhaseVerification:   {PhaseReports},
    PhaseReports:        {}, // Terminal phase
}

// Transition validation
func IsValidTransition(from, to Phase) bool {
    allowed := validTransitions[from]
    for _, p := range allowed {
        if p == to {
            return true
        }
    }
    return false
}
```

---

## Artifact Specifications

### Knowledge Phase Artifacts

| Artifact | Required | Description |
|----------|----------|-------------|
| requirements.md | Yes | System requirements |
| stakeholders.md | Yes | Stakeholder list |
| constraints.md | Yes | Project constraints |
| glossary.md | No | Terminology definitions |

### Architecture Phase Artifacts

| Artifact | Required | Description |
|----------|----------|-------------|
| architecture.md | Yes | System architecture |
| decisions.md | Yes | Architectural decisions |
| components.md | Yes | Component specifications |
| interfaces.md | No | API specifications |

### Implementation Phase Artifacts

| Artifact | Required | Description |
|----------|----------|-------------|
| implementation.md | Yes | Implementation details |
| testing.md | Yes | Testing approach |
| deployment.md | No | Deployment procedures |

### Verification Phase Artifacts

| Artifact | Required | Description |
|----------|----------|-------------|
| verification.md | Yes | Verification procedures |
| test-results.md | Yes | Test execution results |
| coverage.md | Yes | Coverage analysis |

### Reports Phase Artifacts

| Artifact | Required | Description |
|----------|----------|-------------|
| summary.md | Yes | Executive summary |
| findings.md | Yes | Detailed findings |
| recommendations.md | Yes | Recommended actions |

---

## Evidence Requirements

### Evidence Hierarchy

1. **Primary Evidence:** Filesystem artifacts in .kdse/
2. **Secondary Evidence:** Source code, configs, tests
3. **Tertiary Evidence:** Documentation, comments
4. **Assertions:** Claims without supporting evidence (NOT ACCEPTABLE)

### Verification Requirements

For each phase:

```go
type PhaseVerification struct {
    Phase          Phase
    Required       []Artifact
    Verified       []Artifact
    Missing        []Artifact
    Invalid        []Artifact
    AllVerified    bool
    Timestamp      time.Time
    VerifierEvidence []string // Evidence of verification
}
```

### Evidence Collection

```go
func (m *Methodology) CollectEvidence(ctx context.Context, phase Phase) (*Evidence, error) {
    artifacts := m.GetRequiredArtifacts(phase)
    
    var verified []Artifact
    for _, artifact := range artifacts {
        if exists, checksum := m.VerifyArtifact(ctx, artifact); exists {
            verified = append(verified, Artifact{
                Path:     artifact.Path,
                Verified: true,
                Checksum: checksum,
            })
        }
    }
    
    return &Evidence{
        Phase:    phase,
        Artifacts: verified,
        Timestamp: time.Now(),
    }, nil
}
```

---

## Lifecycle Management

### Lifecycle Interface

```go
// Lifecycle defines the methodology's phase lifecycle management
type Lifecycle interface {
    // GetPhases returns all phases in order
    GetPhases() []Phase
    
    // GetPhaseName returns the display name for a phase
    GetPhaseName(phase Phase) string
    
    // GetPhaseDescription returns the description for a phase
    GetPhaseDescription(phase Phase) string
    
    // GetRequiredArtifacts returns required artifacts for a phase
    GetRequiredArtifacts(phase Phase) []ArtifactSpec
    
    // IsValidTransition checks if a phase transition is valid
    IsValidTransition(from, to Phase) bool
    
    // GetNextPhase returns the valid next phase
    GetNextPhase(current Phase) (Phase, error)
    
    // GetPreviousPhase returns the previous phase
    GetPreviousPhase(current Phase) (Phase, error)
    
    // GetCompletionCriteria returns criteria for phase completion
    GetCompletionCriteria(phase Phase) []CompletionCriterion
}
```

### Lifecycle Implementation

```go
// lifecycle implements Lifecycle interface
type lifecycle struct {
    phases []Phase
    specs  map[Phase]PhaseSpec
}

func (l *lifecycle) GetPhases() []Phase {
    return l.phases
}

func (l *lifecycle) GetPhaseName(phase Phase) string {
    return l.specs[phase].Name
}

func (l *lifecycle) GetPhaseDescription(phase Phase) string {
    return l.specs[phase].Description
}

func (l *lifecycle) GetRequiredArtifacts(phase Phase) []ArtifactSpec {
    return l.specs[phase].Artifacts
}

func (l *lifecycle) IsValidTransition(from, to Phase) bool {
    allowed := l.specs[from].AllowedTransitions
    for _, p := range allowed {
        if p == to {
            return true
        }
    }
    return false
}

func (l *lifecycle) GetNextPhase(current Phase) (Phase, error) {
    allowed := l.specs[current].AllowedTransitions
    if len(allowed) == 0 {
        return "", ErrNoNextPhase
    }
    return allowed[0], nil
}

func (l *lifecycle) GetCompletionCriteria(phase Phase) []CompletionCriterion {
    return l.specs[phase].CompletionCriteria
}
```

---

## Phase Specifications

### Initialization Phase

**Purpose:** Establish the runtime and verify workspace readiness.

**Required Artifacts:**
- `.kdse/runtime.yaml`
- `.kdse/workspace.yaml`
- `.kdse/methodology.yaml`
- `.kdse/phase.yaml`
- `.kdse/session.yaml`
- `.kdse/metadata.yaml`

**Completion Criteria:**
- All required files exist
- All files have valid format
- Runtime version is supported
- Workspace path is valid

### Knowledge Phase

**Purpose:** Gather and document requirements.

**Required Artifacts:**
- `knowledge/requirements.md`
- `knowledge/stakeholders.md`
- `knowledge/constraints.md`

**Completion Criteria:**
- All required files exist
- Files are non-empty
- Content meets quality standards

### Architecture Phase

**Purpose:** Design the system architecture.

**Required Artifacts:**
- `architecture/architecture.md`
- `architecture/decisions.md`
- `architecture/components.md`

**Completion Criteria:**
- All required files exist
- Architecture addresses all requirements
- Decisions are documented with rationale
- Components are clearly defined

### Implementation Phase

**Purpose:** Build the system.

**Required Artifacts:**
- `implementation/implementation.md`
- `implementation/testing.md`

**Completion Criteria:**
- All required files exist
- Implementation matches architecture
- Testing approach is defined
- Code follows project standards

### Verification Phase

**Purpose:** Test and validate the implementation.

**Required Artifacts:**
- `verification/verification.md`
- `verification/test-results.md`
- `verification/coverage.md`

**Completion Criteria:**
- All required files exist
- All tests pass
- Coverage meets target
- Verification is complete

### Reports Phase

**Purpose:** Document results and recommendations.

**Required Artifacts:**
- `reports/summary.md`
- `reports/findings.md`
- `reports/recommendations.md`

**Completion Criteria:**
- All required files exist
- Summary is clear
- Findings are specific
- Recommendations are actionable

---

## Authority Model

### Chain of Authority

```
┌─────────────────────────────────────────────────────────────────────┐
│                        KDSE Methodology                              │
│                            (Source of Rules)                         │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Defines
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       Workspace Engine                                │
│                         (Enforces Rules)                              │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ Verifies
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Runtime (.kdse/)                              │
│                          (Evidence Store)                             │
└─────────────────────────────────────────────────────────────────────┘
```

### Enforcement Points

| Point | Enforcement | Action on Failure |
|-------|-------------|------------------|
| Workspace Load | Verify .kdse exists | Return error, require init |
| Phase Transition | Check valid transitions | Return error, list valid |
| Artifact Required | Check file exists | Return error, list missing |
| Artifact Valid | Validate content | Return error, list invalid |

---

## Knowledge Management

### Knowledge Lifecycle

```
┌─────────────┐
│   Create    │ ← Initial knowledge capture
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Verify    │ ← Validate knowledge
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Promote    │ ← Move to architecture
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Archive   │ ← Keep for reference
└─────────────┘
```

### Knowledge Interface

```go
// Knowledge defines knowledge management
type Knowledge interface {
    // CreateKnowledge creates new knowledge artifact
    CreateKnowledge(ctx context.Context, knowledge KnowledgeInput) (*KnowledgeArtifact, error)
    
    // VerifyKnowledge verifies knowledge artifact
    VerifyKnowledge(ctx context.Context, path string) (*VerificationResult, error)
    
    // PromoteKnowledge moves knowledge to next phase
    PromoteKnowledge(ctx context.Context, path string) error
    
    // ListKnowledge returns all knowledge artifacts
    ListKnowledge(ctx context.Context) ([]KnowledgeArtifact, error)
}
```

---

## Architecture Management

### Architecture Lifecycle

```
┌─────────────┐
│   Design    │ ← Create architecture
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Review     │ ← Review against requirements
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Approve    │ ← Approve decisions
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Implement  │ ← Implement from architecture
└─────────────┘
```

### Architecture Interface

```go
// Architecture defines architecture management
type Architecture interface {
    // CreateArchitecture creates architecture artifact
    CreateArchitecture(ctx context.Context, input ArchitectureInput) (*ArchitectureArtifact, error)
    
    // RecordDecision records architectural decision
    RecordDecision(ctx context.Context, decision DecisionInput) (*Decision, error)
    
    // VerifyArchitecture verifies architecture against requirements
    VerifyArchitecture(ctx context.Context) (*VerificationResult, error)
    
    // ListDecisions returns all architectural decisions
    ListDecisions(ctx context.Context) ([]Decision, error)
}
```

---

## Implementation Guidelines

### Dependency Management

```go
// internal/methodology/lifecycle/lifecycle.go
package lifecycle

// This file MUST NOT import runtime packages

import (
    "context"
    "errors"
    
    "kdse/internal/types"  // Only internal types
)

// Errors specific to lifecycle
var (
    ErrNoNextPhase      = errors.New("no next phase available")
    ErrInvalidTransition = errors.New("invalid phase transition")
)

// Lifecycle implementation
type Lifecycle struct {
    // Internal state
}

// This package only depends on:
// - context
// - internal/types
// - standard library
```

### Package Structure

```
internal/
├── methodology/
│   ├── lifecycle/           # Phase lifecycle management
│   │   ├── lifecycle.go     # Lifecycle interface + implementation
│   │   ├── lifecycle_test.go
│   │   └── phases.go        # Phase definitions
│   ├── phases/              # Phase-specific logic
│   │   ├── knowledge/
│   │   ├── architecture/
│   │   ├── implementation/
│   │   ├── verification/
│   │   └── reports/
│   ├── authority/           # Authority and verification
│   │   ├── authority.go
│   │   └── verification.go
│   ├── verification/        # Artifact verification
│   │   ├── verification.go
│   │   └── validators/
│   └── knowledge/           # Knowledge management
│       ├── knowledge.go
│       └── promotion.go
```

---

## Document Relationships

```
METHODOLOGY.md
    │
    ├── Defines: Phases, artifacts, evidence, authority
    │
    ├── Referenced By:
    │   ├── WORKSPACE_ENGINE.md
    │   ├── CLI_RUNTIME.md
    │   ├── MCP_RUNTIME.md
    │   └── All runtime implementations
    │
    ├── References:
    │   ├── PRINCIPLES.md
    │   └── RUNTIME_ARCHITECTURE.md
    │
    └── Related Documents:
        ├── KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md
        └── ADR-002-ONE-METHODOLOGY-MULTIPLE-RUNTIMES.md
```

---

*This document is normative. The methodology defines the engineering rules that all KDSE implementations must follow.*
