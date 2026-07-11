# Phase 0: Runtime Initialization

**Document Version:** 2.0  
**Type:** Runtime Specification  
**Effective Date:** 2026-07-11  

---

## Purpose

Phase 0 ensures the KDSE Runtime automatically loads the KDSE methodology into AI working context before any engineering activity begins.

**Engineering Principle:**
> "The Runtime—not the operator prompt—shall own AI initialization."

---

## Overview

Phase 0 is the mandatory first phase executed automatically on `kdse run`. It establishes the complete KDSE context for AI agents.

```
┌─────────────────────────────────────────────────────────────┐
│                 RUNTIME INITIALIZATION                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Operator Command:                                           │
│                                                              │
│      kdse run                                               │
│                                                              │
│                          │                                   │
│                          ▼                                   │
│  ┌─────────────────────────────────────────────────┐     │
│  │              PHASE 0: INITIALIZE                  │     │
│  │                                                       │     │
│  │  1. Verify Runtime Integrity                       │     │
│  │  2. Verify Runtime Version                         │     │
│  │  3. Load Knowledge Manifest                       │     │
│  │  4. Load Capability Registry                      │     │
│  │  5. Load Command Registry                        │     │
│  │  6. Load Runtime Limitations                      │     │
│  │  7. Generate AI Working Context                  │     │
│  │  8. Generate Runtime Fingerprint                  │     │
│  │  9. Produce Initialization Summary               │     │
│  │                                                       │     │
│  └─────────────────────────────────────────────────┘     │
│                          │                                   │
│                          ▼                                   │
│  ┌─────────────────────────────────────────────────┐     │
│  │              PHASE 1: ASSESS                      │     │
│  │                                                       │     │
│  │  Repository Assessment begins...                     │     │
│  │                                                       │     │
│  └─────────────────────────────────────────────────┘     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Bootstrap Sequence

Phase 0 executes in nine sequential steps:

| Step | Action | Description |
|------|--------|-------------|
| 1 | Verify Runtime Integrity | Check `.kdse/` and `.kdse/bootstrap/` exist |
| 2 | Verify Runtime Version | Confirm version compatibility |
| 3 | Load Knowledge Manifest | Parse `knowledge.yaml` |
| 4 | Load Capability Registry | Parse `capabilities.yaml` |
| 5 | Load Command Registry | Parse `commands.yaml` |
| 6 | Load Runtime Limitations | Parse `limitations.yaml` |
| 7 | Generate AI Working Context | Update `kdse-ai.json` |
| 8 | Generate Runtime Fingerprint | Create integrity hash |
| 9 | Produce Initialization Summary | Display report |

---

## Bootstrap Artifacts

```
.kdse/
├── bootstrap/                     # Phase 0 artifacts
│   ├── knowledge.yaml            # Knowledge Manifest
│   ├── capabilities.yaml         # Capability Registry
│   ├── commands.yaml             # Command Registry
│   ├── limitations.yaml          # Runtime Limitations
│   ├── kdse-ai.json            # AI Working Context
│   └── fingerprints/             # Fingerprint storage
└── runtime/
    └── state.json               # Runtime state
```

---

## Knowledge Manifest

### Required Knowledge

| Order | ID | Source |
|-------|-----|--------|
| 1 | core-principles | docs/foundation/003-core-principles.md |
| 2 | engineering-model | docs/foundation/004-engineering-model.md |
| 3 | chain-of-authority | docs/foundation/006-chain-of-authority.md |
| 4 | glossary | docs/foundation/007-glossary.md |
| 5 | session-protocol | runtime/SESSION_PROTOCOL.md |
| 6 | command-registry | runtime/install/commands.yaml |
| 7 | runtime-configuration | runtime/VERSIONING.md |

### Optional Knowledge

| Order | ID | Source |
|-------|-----|--------|
| 8 | engineering-knowledge | docs/foundation/009-engineering-knowledge.md |
| 9 | traceability | docs/foundation/012-traceability.md |
| 10 | engineering-artifacts | docs/foundation/005-engineering-artifacts.md |
| 11 | audit-standards | docs/audit/COMPLIANCE_AUDIT.md |

---

## Capability Registry

| Capability | Description | Dependencies |
|-----------|-------------|--------------|
| assessment | Repository compliance assessment | - |
| recommendation_engine | Action recommendation | assessment |
| architecture | Architecture design/review | assessment |
| verification | Implementation verification | architecture |
| evolution | Methodology evolution | verification |
| feedback | Feedback collection | - |

---

## Runtime Limitations

| ID | Severity | Description |
|----|----------|-------------|
| no_implementation | info | Runtime does not implement code changes |
| human_approval_required | info | All changes require human approval |
| no_real_time_audit | warning | Audits require explicit invocation |
| session_state_persistence | warning | State not persisted between shells |
| no_code_generation | info | Runtime provides guidance only |
| limited_verification | warning | Static analysis only |
| knowledge_dependency | warning | Requires knowledge documents |

---

## Initialization Summary

Phase 0 produces a human-readable initialization summary:

```
----------------------------------------------------
KDSE Runtime Initialization

Runtime Version:    1.0.0
Knowledge Version:  1.0.0
Runtime Fingerprint: sha256:abc123...def456

Capabilities:
✓ Assessment
✓ Recommendation Engine
✓ Architecture
✓ Verification
✓ Evolution
✓ Feedback

Known Limitations:
• No code implementation - requires human action
• Human approval required for all changes
• Audits require explicit invocation

Knowledge Loaded: 7 documents

Initialization Complete
----------------------------------------------------
```

---

## Runtime State

After initialization, the Runtime state is persisted:

```json
{
  "runtime_version": "1.0.0",
  "knowledge_version": "1.0.0",
  "runtime_fingerprint": "sha256:abc123...",
  "initialized_at": "2026-07-11T00:00:00Z",
  "knowledge_loaded": 7,
  "status": "INITIALIZED"
}
```

---

## Failure Modes

### Runtime Integrity Check Failed

```
ERROR: Runtime integrity check failed
Hint: Run 'kdse install' to reinstall
```

### Version Incompatibility

```
ERROR: Runtime version incompatible
  Current: 1.0.0
  Required: >= 1.0.0
Hint: Run 'kdse update' to upgrade
```

### Missing Required Knowledge

```
ERROR: Required knowledge missing
  Missing: docs/foundation/003-core-principles.md
Hint: Restore missing file from KDSE repository
```

---

## Usage

### Automatic (Recommended)

Phase 0 runs automatically on `kdse run`:

```bash
kdse run
```

### Manual

To manually run Phase 0:

```bash
python3 .kdse/phase0/phase0-init.py
python3 .kdse/phase0/phase0-init.py --verbose
```

### Check Status

```bash
cat .kdse/runtime/state.json
```

---

## Files

| File | Purpose |
|------|---------|
| `.kdse/phase0/phase0-init.py` | Phase 0 bootstrap script |
| `.kdse/bootstrap/knowledge.yaml` | Knowledge Manifest |
| `.kdse/bootstrap/capabilities.yaml` | Capability Registry |
| `.kdse/bootstrap/commands.yaml` | Command Registry |
| `.kdse/bootstrap/limitations.yaml` | Runtime Limitations |
| `.kdse/bootstrap/kdse-ai.json` | AI Working Context |
| `.kdse/runtime/state.json` | Runtime state |

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-07-11 | Initial Phase 0 specification |
| 2.0 | 2026-07-11 | Updated with bootstrap directory structure, capabilities, limitations |

---

*This document defines Phase 0: Runtime Initialization for the KDSE Runtime.*
