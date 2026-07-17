# Workspace Engine

**Document Version:** 1.0  
**Type:** Normative  
**Effective Date:** 2026-07-17

---

## Purpose

This document defines the Workspace Engine—the central component that owns all engineering state and enforces the KDSE methodology.

---

## Overview

### What is the Workspace Engine?

The Workspace Engine is the **single source of truth** for KDSE project state. It:
- Owns all workspace state management
- Enforces methodology rules
- Coordinates runtime operations
- Validates artifacts and phases

### Core Principle

```
The Workspace Engine is the ONLY owner of project state.
All other components (CLI, MCP) MUST go through the Workspace Engine.
```

### Architecture Position

```
┌─────────────────────────────────────────────────────────────────────┐
│                       KDSE Methodology                                │
│                      (Engineering Rules)                              │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Workspace Engine                                │
│                       (State Owner)                                  │
│                                                                     │
│  ✓ Owns .kdse/ directory                                            │
│  ✓ Owns workspace state                                             │
│  ✓ Owns phase state                                                 │
│  ✓ Enforces methodology                                             │
│  ✓ Validates artifacts                                              │
│  ✓ Generates reports                                                │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                    ┌───────────────┴───────────────┐
                    ▼                               ▼
┌───────────────────────────────┐   ┌───────────────────────────────┐
│         CLI Runtime            │   │        MCP Runtime             │
│        (Thin Adapter)          │   │       (Thin Adapter)           │
└───────────────────────────────┘   └───────────────────────────────┘
```

---

## Responsibilities

### Complete Responsibility List

| Responsibility | Description | Priority |
|----------------|-------------|----------|
| Create Runtime | Initialize .kdse directory and files | Critical |
| Load Runtime | Load runtime configuration | Critical |
| Verify Runtime | Verify runtime integrity | Critical |
| Load Phase | Load current phase state | Critical |
| Persist Phase | Save phase state | Critical |
| Validate Workspace | Validate workspace structure | Critical |
| Manage Metadata | Store/retrieve runtime metadata | High |
| Manage Sessions | Track session state | High |
| Generate Reports | Produce verification reports | High |
| Collect Knowledge | Aggregate knowledge artifacts | Medium |
| Generate Architecture | Produce architecture artifacts | Medium |
| Verify Artifacts | Validate artifact completeness | Critical |
| Advance Phase | Handle phase transitions | Critical |

### Non-Responsibilities

The Workspace Engine does NOT:
- Parse command-line arguments
- Handle MCP protocol
- Format output for display
- Implement business logic (delegates to methodology)
- Make platform-specific decisions

---

## Interface Definition

### Main Interface

```go
// Engine is the main interface for the Workspace Engine
type Engine interface {
    // Workspace lifecycle
    InitializeWorkspace(ctx context.Context, opts InitOptions) (*Workspace, error)
    VerifyWorkspace(ctx context.Context) (*VerificationResult, error)
    LoadWorkspace(ctx context.Context) (*Workspace, error)
    DestroyWorkspace(ctx context.Context) error
    
    // Phase management
    GetPhase(ctx context.Context) (*Phase, error)
    AdvancePhase(ctx context.Context, target Phase) (*Transition, error)
    GetPhaseHistory(ctx context.Context) ([]PhaseTransition, error)
    
    // Artifact management
    GetArtifacts(ctx context.Context, phase Phase) ([]Artifact, error)
    VerifyArtifacts(ctx context.Context, phase Phase) (*VerificationResult, error)
    ValidateArtifact(ctx context.Context, artifact Artifact) (*ValidationResult, error)
    
    // Reporting
    GenerateReport(ctx context.Context, opts ReportOptions) (*Report, error)
    ListReports(ctx context.Context) ([]ReportSummary, error)
    
    // Session management
    CreateSession(ctx context.Context, opts SessionOptions) (*Session, error)
    GetSession(ctx context.Context) (*Session, error)
    EndSession(ctx context.Context) error
    
    // Knowledge management
    CollectKnowledge(ctx context.Context) (*KnowledgeCollection, error)
    PromoteKnowledge(ctx context.Context, artifact Artifact) error
}
```

### Supporting Types

```go
// InitOptions contains options for workspace initialization
type InitOptions struct {
    Path      string
    Type      RuntimeType  // cli, mcp
    Version   string
    Template  string
    Metadata  map[string]string
}

// VerificationResult contains the result of workspace verification
type VerificationResult struct {
    Valid         bool
    Phase         Phase
    Errors        []VerificationError
    Warnings      []VerificationWarning
    Timestamp     time.Time
    RuntimeInfo   RuntimeInfo
}

// Phase represents a KDSE engineering phase
type Phase string

const (
    PhaseInitialization  Phase = "initialization"
    PhaseKnowledge       Phase = "knowledge"
    PhaseArchitecture    Phase = "architecture"
    PhaseImplementation  Phase = "implementation"
    PhaseVerification    Phase = "verification"
    PhaseReports         Phase = "reports"
)

// Transition represents a phase transition
type Transition struct {
    From         Phase
    To           Phase
    Timestamp    time.Time
    Verified     bool
    Evidence     []string
}
```

---

## Workspace Structure

### Directory Layout

```
.kdse/
├── runtime.yaml           # Runtime configuration
├── workspace.yaml         # Workspace state
├── methodology.yaml        # Methodology reference
├── phase.yaml              # Current phase state
├── session.yaml             # Session state
├── metadata.yaml            # Runtime metadata
├── knowledge/               # Knowledge artifacts
│   ├── requirements.md
│   ├── stakeholders.md
│   └── constraints.md
├── architecture/           # Architecture artifacts
│   ├── architecture.md
│   └── decisions.md
├── implementation/          # Implementation artifacts
│   └── implementation.md
├── verification/            # Verification artifacts
│   ├── verification.md
│   └── test-results.md
└── reports/                 # Generated reports
    └── summary.md
```

### Workspace State

```go
// Workspace represents a KDSE workspace
type Workspace struct {
    Path       string
    Root       string
    Config     *WorkspaceConfig
    State      *WorkspaceState
    Runtime    *RuntimeConfig
    Session    *Session
    Metadata   *Metadata
}

// WorkspaceState contains current workspace state
type WorkspaceState struct {
    CurrentPhase      Phase
    PhaseHistory      []PhaseTransition
    LastVerification  time.Time
    ArtifactCount     int
    ReportCount       int
    SessionCount      int
}
```

---

## Verification Flow

### Verification Sequence

```go
func (e *Engine) VerifyWorkspace(ctx context.Context) (*VerificationResult, error) {
    // Step 1: Check .kdse exists
    if !e.exists(".kdse") {
        return &VerificationResult{
            Valid:  false,
            Errors: []VerificationError{{
                Code:    "KDSE_MISSING",
                Message: ".kdse directory not found",
            }},
        }, ErrRuntimeMissing
    }
    
    // Step 2: Load runtime configuration
    runtime, err := e.loadRuntimeConfig()
    if err != nil {
        return &VerificationResult{
            Valid:  false,
            Errors: []VerificationError{{
                Code:    "RUNTIME_INVALID",
                Message: err.Error(),
            }},
        }, err
    }
    
    // Step 3: Verify runtime version
    if !e.isVersionSupported(runtime.Version) {
        return &VerificationResult{
            Valid:  false,
            Errors: []VerificationError{{
                Code:    "VERSION_UNSUPPORTED",
                Message: fmt.Sprintf("runtime version %s not supported", runtime.Version),
            }},
        }, ErrVersionUnsupported
    }
    
    // Step 4: Load current phase
    phase, err := e.loadPhase()
    if err != nil {
        return &VerificationResult{
            Valid:  false,
            Errors: []VerificationError{{
                Code:    "PHASE_INVALID",
                Message: err.Error(),
            }},
        }, err
    }
    
    // Step 5: Validate required artifacts
    artifacts, err := e.validatePhaseArtifacts(phase)
    if err != nil {
        return &VerificationResult{
            Valid:  false,
            Phase:  phase,
            Errors: err.Failures,
        }, err
    }
    
    return &VerificationResult{
        Valid:       true,
        Phase:       phase,
        RuntimeInfo: runtime,
        Timestamp:   time.Now(),
    }, nil
}
```

---

## Initialization Flow

### Initialization Sequence

```go
func (e *Engine) InitializeWorkspace(ctx context.Context, opts InitOptions) (*Workspace, error) {
    // Step 1: Create .kdse directory
    if err := e.createRuntimeDirectory(opts.Path); err != nil {
        return nil, err
    }
    
    // Step 2: Download runtime template (if remote)
    template, err := e.downloadTemplate(opts.Template)
    if err != nil {
        e.rollback() // Rollback on failure
        return nil, err
    }
    
    // Step 3: Extract runtime files
    if err := e.extractRuntimeFiles(template); err != nil {
        e.rollback()
        return nil, err
    }
    
    // Step 4: Generate metadata
    if err := e.generateMetadata(opts); err != nil {
        e.rollback()
        return nil, err
    }
    
    // Step 5: Initialize phase to initialization
    if err := e.initializePhase(); err != nil {
        e.rollback()
        return nil, err
    }
    
    // Step 6: Verify runtime (MUST succeed)
    result, err := e.VerifyWorkspace(ctx)
    if err != nil || !result.Valid {
        e.rollback()
        return nil, ErrVerificationFailed
    }
    
    return e.LoadWorkspace(ctx)
}

// Rollback undoes partial initialization
func (e *Engine) rollback() {
    os.RemoveAll(".kdse")
}
```

### Atomic Initialization

```go
// InitializeWorkspace MUST be atomic
// If any step fails, ALL previous steps are rolled back
// No partial initialization is acceptable

func (e *Engine) InitializeWorkspace(ctx context.Context, opts InitOptions) (*Workspace, error) {
    // Create a temporary directory for staging
    tmpDir, err := os.MkdirTemp("", "kdse-init-*")
    if err != nil {
        return nil, err
    }
    defer os.RemoveAll(tmpDir) // Cleanup on any failure
    
    // Stage all operations in tmpDir
    // ...
    
    // Only move to final location if ALL operations succeed
    if err := atomicMove(tmpDir, opts.Path); err != nil {
        return nil, err
    }
    
    // Now .kdse exists with complete state
    return e.LoadWorkspace(ctx)
}
```

---

## Phase Management

### Phase Advancement

```go
func (e *Engine) AdvancePhase(ctx context.Context, target Phase) (*Transition, error) {
    // Step 1: Verify current state
    current, err := e.GetPhase(ctx)
    if err != nil {
        return nil, err
    }
    
    // Step 2: Validate transition
    if !e.lifecycle.IsValidTransition(current, target) {
        return nil, &InvalidTransitionError{
            Current:   current,
            Target:     target,
            ValidNext:  e.lifecycle.GetNextPhase(current),
        }
    }
    
    // Step 3: Verify current phase completion
    artifacts, err := e.VerifyArtifacts(ctx, current)
    if err != nil {
        return nil, err
    }
    if !artifacts.AllVerified {
        return nil, &IncompletePhaseError{
            Phase:   current,
            Missing: artifacts.Missing,
        }
    }
    
    // Step 4: Create transition record
    transition := &Transition{
        From:      current,
        To:        target,
        Timestamp: time.Now(),
        Verified:  true,
    }
    
    // Step 5: Persist new phase
    if err := e.persistPhase(target); err != nil {
        return nil, err
    }
    
    // Step 6: Record transition history
    if err := e.recordTransition(transition); err != nil {
        return nil, err
    }
    
    return transition, nil
}
```

---

## Artifact Management

### Artifact Verification

```go
func (e *Engine) VerifyArtifacts(ctx context.Context, phase Phase) (*VerificationResult, error) {
    // Get required artifacts for phase
    required := e.lifecycle.GetRequiredArtifacts(phase)
    
    var verified []Artifact
    var missing []ArtifactSpec
    var invalid []Artifact
    
    for _, spec := range required {
        artifact, exists := e.loadArtifact(spec.Path)
        if !exists {
            missing = append(missing, spec)
            continue
        }
        
        if err := e.validateArtifact(artifact, spec); err != nil {
            invalid = append(invalid, artifact)
            continue
        }
        
        verified = append(verified, artifact)
    }
    
    return &VerificationResult{
        Phase:       phase,
        Verified:    verified,
        Missing:     missing,
        Invalid:     invalid,
        AllVerified: len(missing) == 0 && len(invalid) == 0,
    }, nil
}

func (e *Engine) validateArtifact(artifact Artifact, spec ArtifactSpec) error {
    // Check file exists
    if !e.exists(artifact.Path) {
        return ErrArtifactMissing
    }
    
    // Check non-empty
    if artifact.Size == 0 {
        return ErrArtifactEmpty
    }
    
    // Check format (if specified)
    if spec.Format != "" && !e.validateFormat(artifact, spec.Format) {
        return ErrArtifactInvalidFormat
    }
    
    // Check content (if validators specified)
    for _, validator := range spec.Validators {
        if err := validator(artifact); err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## Report Generation

### Report Types

```go
type ReportType string

const (
    ReportTypePhase       ReportType = "phase"
    ReportTypeVerification ReportType = "verification"
    ReportTypeSummary      ReportType = "summary"
    ReportTypeProgress     ReportType = "progress"
)

// GenerateReport creates a report of the specified type
func (e *Engine) GenerateReport(ctx context.Context, opts ReportOptions) (*Report, error) {
    switch opts.Type {
    case ReportTypePhase:
        return e.generatePhaseReport(ctx)
    case ReportTypeVerification:
        return e.generateVerificationReport(ctx)
    case ReportTypeSummary:
        return e.generateSummaryReport(ctx)
    case ReportTypeProgress:
        return e.generateProgressReport(ctx)
    default:
        return nil, ErrUnknownReportType
    }
}
```

---

## Session Management

### Session Lifecycle

```go
type Session struct {
    ID        string
    Created   time.Time
    Ended     time.Time
    Phase     Phase
    Actions   []SessionAction
    Metadata  map[string]string
}

type SessionAction struct {
    Timestamp  time.Time
    Action     string
    Phase      Phase
    Result     ActionResult
}

// CreateSession starts a new engineering session
func (e *Engine) CreateSession(ctx context.Context, opts SessionOptions) (*Session, error) {
    session := &Session{
        ID:       generateUUID(),
        Created:  time.Now(),
        Phase:     e.currentPhase,
        Metadata:  opts.Metadata,
        Actions:   []SessionAction{},
    }
    
    if err := e.persistSession(session); err != nil {
        return nil, err
    }
    
    return session, nil
}

// EndSession terminates the current session
func (e *Engine) EndSession(ctx context.Context) error {
    session, err := e.GetSession(ctx)
    if err != nil {
        return err
    }
    
    session.Ended = time.Now()
    
    return e.persistSession(session)
}
```

---

## Error Handling

### Error Types

```go
var (
    ErrRuntimeMissing       = errors.New("runtime not found: .kdse directory missing")
    ErrRuntimeInvalid       = errors.New("runtime configuration invalid")
    ErrVersionUnsupported   = errors.New("runtime version not supported")
    ErrPhaseInvalid         = errors.New("invalid phase state")
    ErrInvalidTransition    = errors.New("invalid phase transition")
    ErrIncompletePhase      = errors.New("phase has incomplete artifacts")
    ErrVerificationFailed   = errors.New("workspace verification failed")
    ErrArtifactMissing      = errors.New("required artifact missing")
    ErrArtifactInvalid      = errors.New("artifact validation failed")
    ErrSessionNotFound      = errors.New("session not found")
)
```

### Error Response Format

```go
type EngineError struct {
    Code    string
    Message string
    Details map[string]interface{}
    Remediation string
}

func (e *EngineError) Error() string {
    return e.Message
}
```

---

## Implementation Structure

### Package Layout

```
internal/
├── workspace/
│   ├── engine/
│   │   ├── engine.go           # Main Engine interface
│   │   ├── engine_test.go
│   │   ├── verify.go           # Verification logic
│   │   ├── init.go             # Initialization logic
│   │   ├── phase.go            # Phase management
│   │   ├── artifact.go         # Artifact management
│   │   ├── session.go          # Session management
│   │   └── report.go           # Report generation
│   ├── loader/
│   │   ├── loader.go           # Workspace loading
│   │   └── config.go           # Configuration loading
│   ├── validator/
│   │   ├── validator.go        # General validation
│   │   └── phase.go            # Phase validation
│   └── state/
│       ├── state.go            # State management
│       └── persistence.go      # State persistence
```

### Dependency Graph

```
workspace/engine
├── workspace/loader
├── workspace/validator
├── workspace/state
├── methodology/lifecycle
└── methodology/authority

workspace/loader
├── workspace/state
└── types

workspace/validator
├── methodology/verification
└── types

workspace/state
└── types
```

---

## Document Relationships

```
WORKSPACE_ENGINE.md
    │
    ├── Defines: Engine interface, responsibilities, flows
    │
    ├── Referenced By:
    │   ├── CLI_RUNTIME.md
    │   ├── MCP_RUNTIME.md
    │   └── All runtime implementations
    │
    ├── References:
    │   ├── RUNTIME_ARCHITECTURE.md
    │   ├── METHODOLOGY.md
    │   └── PRINCIPLES.md
    │
    └── Related Documents:
        └── ADR-003-WORKSPACE-ENGINE.md
```

---

*This document is normative. The Workspace Engine is the authoritative state manager for all KDSE projects.*
