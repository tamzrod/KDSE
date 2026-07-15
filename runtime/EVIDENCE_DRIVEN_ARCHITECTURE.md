# KDSE Evidence-Driven Runtime Architecture

## Fundamental Principle

**KDSE does not automate engineering.**
**KDSE proves engineering.**

Every command reports reality, not assumptions.

---

## Architectural Changes

### Before (Broken)
```
kdse initialize
> KDSE initialized.  ← LIES
```

### After (Verified)
```
kdse initialize
> ╔═══════════════════════════════════════════════════════════════╗
> ║     KDSE Evidence-Driven Runtime Initialization               ║
> ╠═══════════════════════════════════════════════════════════════╣
> ║ PHASE 1: Execute                                              ║
> ║ Creating: runtime/                                            ║
> ║ Creating: foundation/                                         ║
> ║ ...                                                          ║
> ╠═══════════════════════════════════════════════════════════════╣
> ║ VERIFICATION RESULTS                                         ║
> ║ ✓ PASS runtime           /path/.kdse/runtime                 ║
> ║ ✓ PASS foundation        /path/.kdse/foundation              ║
> ║ ...                                                          ║
> ╠═══════════════════════════════════════════════════════════════╣
> ║ Confidence: 1.00                                              ║
> ║ Status: OPERATIONAL                                          ║
```

---

## Command Execution Pattern

Every KDSE command now follows:

```
Execute → Verify → Report
```

### Execute
- Create/update artifacts
- Modify state

### Verify
- Check every artifact created
- Verify existence and readability
- Calculate confidence

### Report
- Return structured evidence
- Include PASS/FAIL for each component
- Include confidence score

---

## New Commands

### `kdse initialize`
Creates a full operational KDSE runtime with verification.

**Required directories:**
- `runtime/` - Runtime execution state
- `foundation/` - Project foundation documents
- `knowledge/` - Collected engineering knowledge
- `laboratory/` - Testing and experimentation
- `evidence/` - Evidence artifacts
- `references/` - Reference materials
- `traceability/` - Requirement traceability
- `reports/` - Generated reports
- `config/` - Runtime configuration
- `state/` - Session and runtime state
- `artifacts/` - Artifact inventory
- `sessions/` - Session history
- `normalized/` - Normalized documentation
- `cache/` - Cached computations

**Required files:**
- `manifest.yaml` - Runtime manifest
- `session-state.yaml` - Session state
- `runtime.yaml` - Runtime configuration
- `knowledge-index.yaml` - Knowledge artifact index
- `artifact-index.yaml` - Artifact inventory

### `kdse runtime verify`
Performs a self-audit of the runtime.

```
╔═══════════════════════════════════════════════════════════════╗
║              KDSE Runtime Self-Audit                         ║
╠═══════════════════════════════════════════════════════════════╣
║ Workspace     PASS     /path/.kdse                           ║
║ Runtime       PASS     /path/.kdse/runtime                   ║
║ Foundation    PASS     /path/.kdse/foundation                ║
║ Knowledge     PASS     /path/.kdse/knowledge                 ║
║ Laboratory    PASS     /path/.kdse/laboratory                ║
║ Config        PASS     /path/.kdse/config                    ║
║ State         PASS     /path/.kdse/state                     ║
║ Manifest      PASS     /path/.kdse/manifest.yaml             ║
║ Session       PASS     /path/.kdse/session-state.yaml        ║
╠═══════════════════════════════════════════════════════════════╣
║ Confidence: 1.00                                               ║
║ Status: OPERATIONAL                                           ║
╚═══════════════════════════════════════════════════════════════╝
```

### `kdse runtime invariant --phase <PHASE>`
Checks if a phase transition is allowed.

**Example:**
```bash
kdse runtime invariant --phase Implementation
```

---

## Runtime Invariants

A project cannot advance phases unless runtime invariants hold.

| Phase | Requires | Description |
|-------|----------|-------------|
| Problem | Runtime initialized | Runtime must be initialized before problem phase |
| Foundation | Runtime initialized | Foundation requires initialized runtime |
| Knowledge | Foundation exists | Knowledge collection requires foundation |
| Architecture | Knowledge collected | Architecture requires knowledge |
| Implementation | Architecture approved | Implementation requires approved architecture |
| Verification | Implementation complete | Verification requires implementation |
| Documentation | Verification complete | Documentation requires verification |
| Audit | All phases complete | Audit requires all phases complete |

---

## Evidence Requirements

Every tool output becomes evidence.

**Not accepted:**
- "I created screenshots."
- "Done."

**Required:**
- Absolute path
- Existence verified
- Count verified
- Manifest generated

---

## Session State Format

```yaml
{
  "version": "1.0.0",
  "session_id": "KDSE-SESSION-20240115-143022",
  "status": "Initialized",
  "phase": "Problem",
  "confidence": 0.0,
  "evidence": [],
  "created_at": "2024-01-15T14:30:22Z",
  "last_verified": null
}
```

---

## Confidence Calculation

Confidence is calculated as:

```
confidence = passed_verifications / total_verifications
```

- 1.00 = All artifacts verified
- 0.00 = No artifacts verified

---

## Phase Transitions

Phase transitions require proof:

1. **Problem → Foundation**: Requires `Runtime initialized`
2. **Foundation → Knowledge**: Requires `Foundation exists`
3. **Knowledge → Architecture**: Requires `Knowledge collected`
4. **Architecture → Implementation**: Requires `Architecture approved`
5. **Implementation → Verification**: Requires `Implementation complete`
6. **Verification → Documentation**: Requires `Verification complete`
7. **Documentation → Audit**: Requires `All phases complete`

---

## Implementation Notes

### Files Created

- `/internal/runtime/runtime.go` - Core runtime implementation
- `/internal/runtime/invariant.go` - Invariant engine
- `/cmd/kdse/main.go` - Updated CLI with new commands

### Key Types

```go
// VerificationResult represents verification status
type VerificationResult struct {
    Artifact  string
    Path      string
    Status    string  // "PASS" or "FAIL"
    Evidence  string
    Error     string
    Timestamp string
}

// InitializeResult contains initialization results
type InitializeResult struct {
    Success       bool
    WorkspacePath string
    Confidence    float64
    Verification  []VerificationResult
    Evidence      []string
    Errors        []string
    Timestamp     string
}

// VerificationReport is the result of runtime verification
type VerificationReport struct {
    Success      bool
    Confidence   float64
    Components   []VerificationResult
    Missing      []string
    Failed       []string
    Timestamp    string
}
```

---

## Migration Path

For existing projects:

1. Run `kdse initialize` to create new structure
2. Run `kdse runtime verify` to confirm
3. Migrate existing artifacts to new directories
4. Update phase via `kdse context stage --to <PHASE>`

---

## Backward Compatibility

The new evidence-driven commands are additive:

- `kdse initialize` - NEW (evidence-driven)
- `kdse runtime verify` - NEW (self-audit)
- `kdse runtime invariant` - NEW (phase checks)
- Existing commands unchanged
