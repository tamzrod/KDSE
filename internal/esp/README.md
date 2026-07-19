# KDSE Engineering Session Protocol (ESP) Implementation

This package implements the KDSE Engineering Session Protocol as defined in the ESP specification (version 1.0).

## Overview

The ESP protocol establishes the complete contractual framework for engineering sessions, encompassing:

- **Identity establishment** — Determination of participant roles and authorities
- **Context establishment** — Creation and verification of engineering context
- **Session lifecycle** — Definition of all valid session states and transitions
- **Failure handling** — Deterministic recovery procedures for all failure modes
- **Completion criteria** — Unambiguous conditions for session activation

## Components

### Session State Machine

The protocol defines 9 session states:

| State | Description |
|-------|-------------|
| `IDLE` | No active engineering session exists |
| `DISCOVERING` | Runtime is locating and identifying the workspace |
| `AUTHENTICATING` | Runtime is verifying workspace ownership |
| `BOOTSTRAPPING` | Runtime is establishing engineering context |
| `VERIFYING` | Runtime is validating engineering context completeness |
| `ACTIVE` | Engineering session is established |
| `SUSPENDED` | Session is paused but may be resumed |
| `TERMINATED` | Session has ended and cannot be resumed |
| `FAILED` | An unrecoverable error has occurred |

### Handshake Protocol

The handshake protocol consists of 12 message types that establish session identity and authority:

1. `CONNECT` - AI Agent initiates session
2. `DISCOVER` / `DISCOVER_ACK` - Workspace discovery
3. `AUTHENTICATE` / `AUTHENTICATE_ACK` - Authentication
4. `CONTEXT_REQUEST` / `CONTEXT_QUERY` / `CONTEXT_RESPONSE` - Context exchange
5. `BOOTSTRAP` / `BOOTSTRAP_ACK` - Context establishment
6. `VERIFY` / `VERIFY_ACK` - Context verification

### Bootstrap Protocol

The bootstrap process consists of 8 steps:

1. Discover Workspace
2. Verify Runtime
3. Load Manifest
4. Determine Context
5. Resolve Conflicts
6. Assemble Context
7. Compute Checksum
8. Record Evidence

### Authority Hierarchy

The protocol enforces the following authority hierarchy:

1. KDSE Standard (highest)
2. Runtime Manifest
3. Project Manifest
4. Session State
5. AI Reasoning (lowest)

### Failure Codes

The protocol defines 20 failure codes covering:

- Discovery failures (1001-1099)
- Authentication failures (1101-1199)
- Bootstrap failures (1201-1299)
- Verification failures (1301-1399)
- Version failures (1501-1599)

## Usage

### Creating a Session

```go
import "github.com/kdse/runtime/internal/esp"

// Create a new session
session := esp.NewSession("workspace-123")

// Transition through states
if err := session.TransitionTo(esp.StateDiscovering, "Starting session"); err != nil {
    log.Fatalf("Invalid transition: %v", err)
}
```

### Bootstrapping Context

```go
// Bootstrap the engineering context
protocol := esp.NewBootstrapProtocol("/path/to/repo")
result, err := protocol.Execute(session.ID)
if err != nil {
    log.Fatalf("Bootstrap failed: %v", err)
}

fmt.Printf("Bootstrap successful: %v\n", result.Success)
```

### Verifying Completion Criteria

```go
// Verify all conditions for ACTIVE state
verifier := esp.NewCompletionVerifier("/path/to/repo")
criteria := verifier.VerifyActiveState(result.Context)

fmt.Printf("Ready for AI reasoning: %v\n", criteria.AllPassed)
fmt.Printf("Confidence: %.0f%%\n", criteria.Confidence*100)
```

## Compliance

This implementation satisfies all mandatory requirements defined in Section 15.2 of the ESP specification:

- C-001: All handshake messages implemented
- C-002: All bootstrap steps executed
- C-003: Authority hierarchy enforced
- C-004: All state transitions implemented
- C-005: All failure codes handled
- C-006: Completion criteria verified
- C-007: No reasoning before ACTIVE state
- C-008: Deterministic context generation
- C-009: Authoritative sources take precedence
- C-010: Implementation independence maintained

## Files

| File | Description |
|------|-------------|
| `esp.go` | Package documentation and version |
| `protocol.go` | Core protocol types and session states |
| `session.go` | Session management and state machine |
| `bootstrap.go` | Bootstrap protocol implementation |
| `failure.go` | Failure codes and recovery procedures |
| `completion.go` | Completion criteria verification |
| `compliance_test.go` | Conformance test suite |

## References

- [KDSE Engineering Session Protocol Specification](.agents_tmp/PLAN.md)
- [Context Handoff Protocol](../runtime/CONTEXT_HANDOFF.md)
