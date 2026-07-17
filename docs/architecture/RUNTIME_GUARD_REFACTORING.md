# Runtime Guard Architecture Refactoring - Impact Summary

**Document Version:** 1.0  
**Date:** 2026-07-17  
**Type:** Architecture Impact Report

---

## Executive Summary

This document summarizes the architectural refactoring of the KDSE runtime initialization system. The refactoring separates concerns into dedicated guards, improving maintainability and following the single responsibility principle.

### Key Changes

1. **Introduced Runtime Guard** - New orchestrator for guard execution
2. **Created Project Guard** - Separated project discovery and validation
3. **Created Workspace Guard** - Dedicated workspace validation
4. **Created Lifecycle Guard** - Dedicated lifecycle state management
5. **Refactored Session Guard** - Removed project/workspace logic
6. **Created Coordinator** - High-level initialization interface

---

## Architecture Overview

### Before (Problem)

```
Session Guard
├── Workspace validation ❌ (should be separate)
├── Project detection ❌ (should be separate)
├── Session management ✓
└── Lifecycle state ❌ (should be separate)
```

### After (Solution)

```
Runtime Guard (Orchestrator)
    │
    ├── Project Guard
    │      ├── Project Discovery ✓
    │      └── Project Validation ✓
    │
    ├── Workspace Guard
    │      ├── Workspace Discovery ✓
    │      └── Workspace Integrity ✓
    │
    ├── Session Guard (Validation)
    │      ├── Session Detection ✓
    │      ├── Session Validity ✓
    │      └── Session Expiration ✓
    │
    └── Lifecycle Guard
           ├── Current Phase ✓
           ├── Allowed Transitions ✓
           └── Lifecycle Integrity ✓
```

---

## New File Structure

```
internal/guard/
├── types.go                      # Shared types and error definitions
├── runtime_guard.go             # Orchestrator (new)
├── project_guard.go             # Project discovery (new)
├── workspace_guard.go           # Workspace validation (new)
├── session_validation_guard.go  # Session validation only (new)
├── lifecycle_guard.go           # Lifecycle management (new)
├── coordinator.go               # High-level initialization (new)
└── session_guard.go             # Legacy (refactored, deprecated)
```

---

## Guard Responsibilities

### Runtime Guard

| Responsibility | Description |
|---------------|-------------|
| Entry point | First guard called for every runtime command |
| Orchestration | Executes guards in deterministic order |
| Fail-fast | Stops on first failure |
| No validation | Delegates all validation to specialized guards |

### Project Guard

| Responsibility | Description |
|---------------|-------------|
| Project detection | Detects engineering projects |
| Project validation | Validates project eligibility |
| Artifact discovery | Identifies project artifacts |
| Git detection | Checks for git repository |

**Must NOT:**
- Initialize KDSE
- Create .kdse directory
- Manage sessions
- Handle lifecycle

### Workspace Guard

| Responsibility | Description |
|---------------|-------------|
| Workspace discovery | Finds existing KDSE workspace |
| Structure validation | Validates workspace directories |
| Version compatibility | Checks version compatibility |
| Integrity checks | Validates internal consistency |

**Assumes:** Project Guard already succeeded

**Must NOT:**
- Discover projects
- Manage sessions
- Handle lifecycle

### Session Validation Guard

| Responsibility | Description |
|---------------|-------------|
| Session detection | Finds active session |
| Session validity | Validates session state |
| Session expiration | Checks 24-hour expiry |
| Session recovery | Handles session state |

**Assumes:** Project Guard and Workspace Guard already succeeded

**Must NOT:**
- Discover projects
- Validate projects
- Create workspaces
- Initialize projects

### Lifecycle Guard

| Responsibility | Description |
|---------------|-------------|
| Phase management | Tracks current phase |
| Transition validation | Validates allowed transitions |
| History tracking | Records phase history |

**Must NOT:**
- Discover projects
- Manage workspaces
- Handle sessions

---

## State Model

```
No Project
    ↓
Project (validated by Project Guard)
    ↓
Workspace (validated by Workspace Guard)
    ↓
Session (validated by Session Guard)
    ↓
Lifecycle Ready (validated by Lifecycle Guard)
```

### State Definitions

| State | Guard | Description |
|-------|-------|-------------|
| `NO_PROJECT` | Project Guard | No valid project detected |
| `PROJECT` | Project Guard | Valid project detected |
| `WORKSPACE` | Workspace Guard | KDSE workspace exists |
| `SESSION` | Session Guard | Active session running |
| `LIFECYCLE_READY` | Lifecycle Guard | Ready for operations |

---

## Execution Flow

### Guard Execution Order

```
Runtime Guard
    │
    ├── Project Guard → Validate project
    │       │
    │       └── FAIL → Return error (NO_PROJECT)
    │
    ├── Workspace Guard → Validate workspace
    │       │
    │       └── FAIL → Return error (WORKSPACE_MISSING)
    │
    ├── Session Guard → Validate session
    │       │
    │       └── FAIL → Return error (NO_SESSION)
    │
    └── Lifecycle Guard → Validate lifecycle
            │
            └── FAIL → Return error (LIFECYCLE_INVALID)
```

### Stop-on-First-Failure

Each guard runs only if the previous guard succeeded. This ensures:
1. Clear error messages
2. Deterministic behavior
3. No wasted computation
4. Easy debugging

---

## Backward Compatibility

### Legacy Session Guard

The original `SessionGuard` has been refactored but preserved for backward compatibility:

```go
// OLD (still works but deprecated)
guard := NewSessionGuard(repoPath)
guard.EnforceInitialized()

// NEW (recommended)
guard := NewRuntimeGuard(repoPath)
result := guard.Validate(ctx)
```

### Migration Path

1. **Immediate**: Existing code continues to work
2. **Short-term**: Update imports to use new guards
3. **Long-term**: Remove legacy SessionGuard

---

## Error Handling

### Error Structure

All guard errors follow a consistent structure:

```go
type RuntimeGuardError struct {
    GuardType GuardType  // Which guard produced the error
    Code      string     // Machine-readable code
    Message   string     // Human-readable message
    Hint      string     // How to fix the issue
    State     RuntimeState // Current state when error occurred
}
```

### Error Codes by Guard Type

| Guard | Code | Description |
|-------|------|-------------|
| Project | `NO_PROJECT` | No engineering project detected |
| Project | `INVALID_LOCATION` | Project location invalid |
| Project | `GENERIC_DIRECTORY` | Not an engineering project |
| Workspace | `WORKSPACE_MISSING` | .kdse directory not found |
| Workspace | `WORKSPACE_CORRUPTED` | Workspace incomplete |
| Workspace | `VERSION_MISMATCH` | Incompatible version |
| Session | `NO_SESSION` | No active session |
| Session | `SESSION_EXPIRED` | Session expired |
| Session | `SESSION_INVALID` | Session state invalid |
| Lifecycle | `LIFECYCLE_INVALID` | Lifecycle state invalid |
| Lifecycle | `TRANSITION_INVALID` | Invalid phase transition |

---

## API Reference

### Runtime Guard

```go
// Create runtime guard
guard := NewRuntimeGuard(repoPath)

// Full validation (stops on first failure)
result := guard.Validate(ctx)

// Quick check (all guards run)
result := guard.QuickCheck(ctx)

// Command enforcement
err := guard.EnforceForCommand(ctx, "initialize")

// Get current state (no validation)
state := guard.GetCurrentState()
```

### Project Guard

```go
// Create project guard
guard := NewProjectGuard(repoPath)

// Validate project
result := guard.Validate(ctx)

// Quick existence check
exists := guard.Exists()
```

### Workspace Guard

```go
// Create workspace guard
guard := NewWorkspaceGuard(repoPath)

// Validate workspace
result := guard.Validate(ctx)

// Quick existence check
exists := guard.Exists()

// Get workspace info
info, err := guard.GetInfo()
```

### Session Validation Guard

```go
// Create session guard
guard := NewSessionValidationGuard(repoPath)

// Validate session
result := guard.Validate(ctx)

// Check for active session
hasSession := guard.HasActiveSession()

// Get session ID
sessionID, err := guard.GetSessionID()

// Create new session
state, err := guard.CreateSession()

// End session
err := guard.EndSession()
```

### Lifecycle Guard

```go
// Create lifecycle guard
guard := NewLifecycleGuard(repoPath)

// Validate lifecycle
result := guard.Validate(ctx)

// Check validity
isValid := guard.IsValid()

// Get current phase
phase, err := guard.GetCurrentPhase()

// Set phase
err := guard.SetPhase("architecture")
```

### Coordinator

```go
// Create coordinator
coord := NewCoordinator(repoPath)

// Full initialization
err := coord.Initialize(ctx)

// Validation
result := coord.Validate(ctx)

// Status
status := coord.Status()
```

---

## Testing Strategy

### Unit Tests

Each guard should have unit tests covering:

1. Happy path validation
2. Error cases for each error type
3. State transitions
4. Edge cases

### Integration Tests

Coordinator tests should verify:

1. Complete initialization flow
2. Error propagation
3. State machine transitions

### Backward Compatibility Tests

Ensure legacy SessionGuard still works correctly.

---

## Performance Considerations

1. **Lazy initialization**: Guards don't initialize resources until needed
2. **Quick checks**: `Exists()` methods for fast pre-checks
3. **Fail-fast**: Early termination on first failure
4. **Parallel evaluation**: Guards can be evaluated in parallel if needed

---

## Security Implications

1. **Separation of concerns**: Each guard has limited scope
2. **Fail-safe defaults**: Invalid state defaults to most restrictive
3. **Clear error messages**: Each error includes remediation hints
4. **Audit trail**: State transitions are logged

---

## Recommendations

### For New Code

Use the new guard architecture:

```go
guard := NewRuntimeGuard(repoPath)
result := guard.Validate(ctx)
if !result.Valid {
    // Handle error
}
```

### For Existing Code

The legacy SessionGuard continues to work. Consider migrating at next opportunity:

```go
// Before
guard := NewSessionGuard(repoPath)
guard.EnforceInitialized()

// After
guard := NewRuntimeGuard(repoPath)
if err := guard.EnforceForCommand(ctx, "your-command"); err != nil {
    // Handle error
}
```

---

## Conclusion

This refactoring achieves the following goals:

1. ✅ **Separation of concerns**: Each guard has a single responsibility
2. ✅ **Clear dependencies**: Guards form a deterministic chain
3. ✅ **Fail-fast behavior**: Errors stop execution immediately
4. ✅ **Backward compatibility**: Legacy SessionGuard preserved
5. ✅ **Testability**: Each guard can be tested in isolation
6. ✅ **Maintainability**: Changes to one guard don't affect others

### Files Changed

| File | Action |
|------|--------|
| `internal/guard/types.go` | Created |
| `internal/guard/runtime_guard.go` | Created |
| `internal/guard/project_guard.go` | Created |
| `internal/guard/workspace_guard.go` | Created |
| `internal/guard/session_validation_guard.go` | Created |
| `internal/guard/lifecycle_guard.go` | Created |
| `internal/guard/coordinator.go` | Created |
| `internal/guard/session_guard.go` | Refactored (deprecated) |

### Files Preserved

| File | Notes |
|------|-------|
| `internal/guard/session_guard.go` | Kept for backward compatibility |

---

*This document describes the refactored guard architecture. For implementation details, see source code documentation.*
