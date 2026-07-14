# Context Handoff Protocol

**Document Version:** 1.0  
**Type:** Normative Runtime Specification  
**Effective Date:** 2026-07-14

---

## Purpose

The Context Handoff Protocol enables seamless continuity between AI sessions by maintaining a structured handoff artifact that captures completed work and directs subsequent actions.

Without context handoff, every new sandbox session must re-discover completed work by scanning the repository. This wastes compute and introduces inconsistencies as the AI may reach different conclusions about project state.

---

## Overview

```
.kdse/
├── context.json          # Primary handoff artifact
├── reports/              # Generated reports (KDSE-RUN-*.md)
└── evidence/
    ├── screenshots/      # Visual evidence
    ├── tests/            # Test results and logs
    └── benchmarks/       # Performance data
```

---

## context.json Schema

```json
{
  "project": "forge",
  "project_version": "1.0.0",
  "schema": "https://kdse.dev/schemas/context-handoff/v1",
  "current_stage": "UX Review",
  "previous_stage": "Visual Acceptance Test",
  "stage_history": [
    {"stage": "Concept", "completed_at": "2026-07-10T10:00:00Z"},
    {"stage": "Architecture", "completed_at": "2026-07-11T14:30:00Z"},
    {"stage": "Visual Acceptance Test", "completed_at": "2026-07-13T09:15:00Z"}
  ],
  "evidence": [
    "docs/screenshots/dashboard-v2.png",
    ".kdse/reports/KDSE-RUN-20260713.md"
  ],
  "next_action": "Perform UX Audit",
  "allowed_context": [
    "docs/screenshots",
    "docs/design/",
    ".kdse/context.json"
  ],
  "restricted_paths": [
    "vendor/",
    "node_modules/",
    ".cache/"
  ],
  "session": {
    "session_id": "KDSE-RT-2026-07-14-080123",
    "started_at": "2026-07-14T08:01:23Z",
    "last_updated": "2026-07-14T08:45:00Z",
    "sandbox_id": "sandbox-abc123"
  },
  "artifacts": {
    "reports": ".kdse/reports/",
    "screenshots": ".kdse/evidence/screenshots/",
    "tests": ".kdse/evidence/tests/",
    "benchmarks": ".kdse/evidence/benchmarks/"
  },
  "metadata": {
    "initialized_at": "2026-07-10T10:00:00Z",
    "last_transition": "2026-07-13T09:15:00Z",
    "transitions_count": 3
  }
}
```

### Field Definitions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `project` | string | Yes | Project identifier |
| `project_version` | string | No | Semantic version of project |
| `schema` | string | Yes | JSON Schema URI for validation |
| `current_stage` | string | Yes | Current KDSE stage name |
| `previous_stage` | string | No | Previously completed stage |
| `stage_history` | array | No | Array of {stage, completed_at} objects |
| `evidence` | array | No | Paths to evidence files |
| `next_action` | string | Yes | Explicit next action for new sessions |
| `allowed_context` | array | No | Paths the AI may read without re-scanning |
| `restricted_paths` | array | No | Paths to skip during discovery |
| `session.*` | object | No | Current session metadata |
| `artifacts.*` | object | No | Artifact directory paths |
| `metadata.*` | object | No | Lifecycle metadata |

---

## Stage Definitions

KDSE defines the following major stages:

| Stage | Description |
|-------|-------------|
| `Concept` | Initial project definition |
| `Knowledge Gathering` | Domain knowledge collection |
| `Architecture` | System design and structure |
| `Implementation` | Code development |
| `Verification` | Testing and validation |
| `Visual Acceptance Test` | UI/UX validation |
| `UX Review` | User experience audit |
| `Documentation` | Documentation completion |
| `Deployment` | Release preparation |

---

## Session Lifecycle

### Session Start

1. **Read context.json** (first action, not repository scan)
2. **Validate schema** against `https://kdse.dev/schemas/context-handoff/v1`
3. **Update session fields**:
   - `session.session_id` = new ID
   - `session.started_at` = current time
   - `session.last_updated` = current time
   - `session.sandbox_id` = current sandbox ID
4. **Read allowed_context paths** only
5. **Execute next_action**

### Session End

1. **Update context.json**:
   - Set `previous_stage` = `current_stage` (if changed)
   - Append to `stage_history` (if stage completed)
   - Update `evidence` with new artifacts
   - Set `next_action` for next session
   - Update `session.last_updated`
   - Increment `metadata.transitions_count`

---

## API Operations

### Read Context

```bash
# Read current handoff context
kdse context read

# Output: Full context.json content
```

### Update Stage

```bash
# Transition to new stage
kdse context stage --to "Architecture" --evidence "docs/arch.md"

# Actions:
# 1. Set previous_stage = current_stage
# 2. Set current_stage = "Architecture"
# 3. Append to stage_history
# 4. Add evidence files
# 5. Update metadata.last_transition
```

### Set Next Action

```bash
# Define next action for new sessions
kdse context next-action "Review architectural decisions"

# Sets next_action field
```

### Add Evidence

```bash
# Add evidence file reference
kdse context add-evidence "docs/screenshots/dashboard.png"

# Appends to evidence array
```

### Generate Handoff Report

```bash
# Create session handoff report
kdse context report --stage "Architecture"

# Output: .kdse/reports/KDSE-RUN-YYYYMMDD.md
```

---

## Acceptance Criteria

| # | Criterion | Verification |
|---|-----------|--------------|
| 1 | New sandbox resumes immediately | context.json read in < 1 second |
| 2 | AI does not rediscover completed work | No repository scan on start |
| 3 | Repository scanning is minimized | Only allowed_context scanned |
| 4 | Next action is explicit | `next_action` field populated |

---

## Schema Validation

Context files MUST validate against the JSON Schema at `https://kdse.dev/schemas/context-handoff/v1`.

Required fields:
- `project` (string)
- `schema` (string, must match URI)
- `current_stage` (string)
- `next_action` (string)

---

## Example Usage

### Initialize New Project

```bash
# Initialize project context
kdse context init --project "myapp" --stage "Concept"

# Creates .kdse/context.json with:
# - project: "myapp"
# - current_stage: "Concept"
# - next_action: "Define domain knowledge"
```

### Transition Stages

```bash
# Complete Concept, move to Architecture
kdse context stage --to "Architecture" --evidence "docs/domain.md,docs/requirements.md"

# Result:
# - previous_stage: "Concept"
# - stage_history: [{"stage": "Concept", "completed_at": "2026-07-14T10:00:00Z"}]
# - current_stage: "Architecture"
# - evidence: ["docs/domain.md", "docs/requirements.md"]
```

### Resume Interrupted Work

```bash
# New sandbox starts
# Reads context.json
# Finds current_stage: "Architecture"
# Finds next_action: "Review domain model"
# Directly reads allowed_context paths
# Executes next_action
```

---

## Implementation

See `internal/context/context.go` for the Go implementation of context operations.

---

*This document is part of the KDSE Runtime Specification.*
