# Phase 0: Runtime Initialization

**Document Version:** 1.0  
**Type:** Runtime Specification  
**Effective Date:** 2026-07-11  

---

## Purpose

Phase 0 ensures the KDSE Runtime automatically loads the KDSE methodology into AI working context before any engineering activity begins.

**Engineering Principle:**
> "The Runtime—not the operator prompt—shall own AI initialization."

---

## Overview

Phase 0 is the mandatory first phase executed before any engineering activity. It establishes the complete KDSE context for AI agents, eliminating the need for manual bootstrap prompts.

```
┌─────────────────────────────────────────────────────────────┐
│                    PHASE 0: INITIALIZATION                   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Every KDSE Runtime Session begins with Phase 0.            │
│  The operator should only need to issue:                    │
│                                                              │
│      Run KDSE.                                              │
│                                                              │
│  All methodology loading occurs automatically.               │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Bootstrap Sequence

Phase 0 executes in eight sequential steps:

```
┌─────────────────────────────────────────────────────────────┐
│                    BOOTSTRAP SEQUENCE                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Step 1. Discover Installation                              │
│      └─ Verify .kdse/ directory exists                      │
│                                                              │
│  Step 2. Load Manifest                                       │
│      └─ Parse .kdse/knowledge/manifest.yaml                 │
│                                                              │
│  Step 3. Verify Versions                                     │
│      └─ Check Runtime/Standard compatibility                │
│                                                              │
│  Step 4. Load Knowledge (in order)                           │
│      └─ Load required, then optional documents               │
│                                                              │
│  Step 5. Verify Integrity                                    │
│      └─ Generate and verify knowledge fingerprint            │
│                                                              │
│  Step 6. Discover Capabilities                               │
│      └─ Build capability registry                           │
│                                                              │
│  Step 7. Generate AI Context                                 │
│      └─ Create kdse-ai.json for AI context                 │
│                                                              │
│  Step 8. Produce Initialization Summary                     │
│      └─ Display initialization report                       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Knowledge Manifest

The Knowledge Manifest (`.kdse/knowledge/manifest.yaml`) defines:

### Required Knowledge

Must be loaded before any engineering activity begins:

| Order | Document | Source |
|-------|----------|--------|
| 1 | Core Principles | docs/foundation/003-core-principles.md |
| 2 | Engineering Model | docs/foundation/004-engineering-model.md |
| 3 | Chain of Authority | docs/foundation/006-chain-of-authority.md |
| 4 | Glossary | docs/foundation/007-glossary.md |
| 5 | Session Protocol | runtime/SESSION_PROTOCOL.md |
| 6 | Command Registry | runtime/install/commands.yaml |
| 7 | Runtime Configuration | runtime/VERSIONING.md |

### Optional Knowledge

Should be loaded for complete KDSE context:

| Order | Document | Source |
|-------|----------|--------|
| 8 | Engineering Knowledge Definition | docs/foundation/009-engineering-knowledge.md |
| 9 | Traceability Framework | docs/foundation/012-traceability.md |
| 10 | Engineering Artifacts | docs/foundation/005-engineering-artifacts.md |
| 11 | Audit Standards | docs/audit/COMPLIANCE_AUDIT.md |

---

## Knowledge Loading Order

The AI shall not guess which documents to load. The Runtime explicitly defines the initialization order:

```
Loading Order
═════════════

  1. Core Principles
     ↓
  2. Engineering Model
     ↓
  3. Chain of Authority
     ↓
  4. Glossary
     ↓
  5. Session Protocol
     ↓
  6. Command Registry
     ↓
  7. Runtime Configuration
     ↓
  [Optional: 8-11]
```

---

## Knowledge Fingerprint

The Knowledge Fingerprint is a SHA-256 hash that verifies knowledge integrity:

```
Fingerprint = SHA256(
  sorted([
    "source:path" + SHA256(content),
    ...
  ])
)
```

**Purpose:**
- Detects unauthorized knowledge changes
- Ensures reproducibility of initialization
- Provides audit trail

---

## Capability Discovery

After Phase 0 initialization, the following capabilities are available:

| Capability | Description | Entrypoint |
|-----------|-------------|------------|
| Assessment | Repository compliance assessment | Run KDSE |
| Architecture | Architecture design and review | Engineering Phase |
| Implementation | Implementation guidance | Engineering Phase |
| Verification | Implementation verification | Verification Phase |
| Evolution | Methodology evolution | Evolution Phase |
| Feedback | Feedback collection | Session Protocol |

---

## AI Knowledge Artifact

The AI Knowledge Artifact (`.kdse/knowledge/kdse-ai.json`) contains machine-readable engineering knowledge:

```json
{
  "$schema": "https://kdse.dev/schemas/kdse-ai/v1.0",
  "version": "1.0.0",
  "runtime": {
    "version": "1.0.0",
    "name": "KDSE Runtime",
    "compatible_standard": ">= 1.0.0"
  },
  "knowledge": {
    "version": "1.0.0",
    "fingerprint": "sha256:...",
    "loaded": ["core-principles", "engineering-model", ...],
    "status": "READY"
  },
  "capabilities": {
    "assessment": { ... },
    "architecture": { ... },
    "verification": { ... },
    "evolution": { ... },
    "feedback": { ... }
  },
  "principles": {
    "core": [...]
  },
  "status": "READY"
}
```

---

## Initialization Summary

Phase 0 produces a human-readable initialization summary:

```
═══════════════════════════════════════════════════════════════
                    KDSE Runtime Initialized
═══════════════════════════════════════════════════════════════

Runtime Version:    1.0.0
Knowledge Version:  1.0.0
Knowledge Fingerprint: b84085605a4477f8...

Capabilities Loaded:
  ✓ Assessment
  ✓ Architecture
  ✓ Verification
  ✓ Evolution
  ✓ Feedback

Knowledge Loaded:
  7 documents

Repository:
  Path: /workspace/project/KDSE
  Lifecycle: Active

Status: READY

═══════════════════════════════════════════════════════════════
```

---

## Runtime State

After initialization, the Runtime state is persisted to `.kdse/runtime/state.json`:

```json
{
  "runtime_version": "1.0.0",
  "knowledge_version": "1.0.0",
  "knowledge_fingerprint": "sha256:...",
  "compatible_standard": ">= 1.0.0",
  "initialized_at": "2026-07-11T00:00:00Z",
  "repository_path": "/workspace/project/KDSE",
  "knowledge_loaded": 7,
  "status": "READY"
}
```

---

## Failure Modes

### Missing Installation

```
ERROR: KDSE Runtime not installed
Hint: Run ./runtime/install/install.sh
```

### Invalid Manifest

```
ERROR: Invalid Knowledge Manifest
Hint: manifest.yaml is corrupted or missing required fields
```

### Version Incompatibility

```
ERROR: Version incompatibility detected
  Runtime Version: 1.0.0
  Standard Version: 0.9.0
  Required: >= 1.0.0
Hint: Run 'kdse update' to sync with compatible version
```

### Missing Required Knowledge

```
ERROR: Required knowledge missing
  Missing: docs/foundation/003-core-principles.md
  Required by: manifest.yaml
Hint: Restore missing file from KDSE repository
```

---

## Usage

### Automatic Initialization

Phase 0 runs automatically when starting a KDSE session:

```
Run KDSE.
```

### Manual Initialization

To manually run Phase 0:

```bash
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
| `.kdse/phase0/runtime-state.sh` | Runtime state management |
| `.kdse/phase0/load-knowledge.sh` | Knowledge loading script |
| `.kdse/phase0/generate-fingerprint.sh` | Fingerprint generation |
| `.kdse/knowledge/manifest.yaml` | Knowledge Manifest |
| `.kdse/knowledge/kdse-ai.json` | AI Knowledge Artifact |
| `.kdse/runtime/state.json` | Runtime state |

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-07-11 | Initial Phase 0 specification |

---

*This document defines Phase 0: Runtime Initialization for the KDSE Runtime.*
