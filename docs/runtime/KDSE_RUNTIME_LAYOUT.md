# KDSE Runtime Layout

**Document Version:** 1.0  
**Type:** Normative Runtime Specification  
**Effective Date:** 2026-07-10  

---

## Purpose

This document defines the complete directory layout of the `.kdse/` Runtime Environment, explaining the purpose of each component and how they relate to each other.

---

## Complete Directory Structure

```
.kdse/                                    # Root - Runtime Environment
│
├── standards/                            # Normative Standards
│   ├── foundation/                       # Core KDSE Definitions
│   │   ├── 000-what-is-kdse.md         # What KDSE Is
│   │   ├── 001-why-kdse-exists.md       # Why KDSE Exists
│   │   ├── 002-scope.md                 # Scope Boundaries
│   │   ├── 003-core-principles.md       # Core Principles
│   │   ├── 004-engineering-model.md     # Engineering Model
│   │   ├── 005-engineering-artifacts.md # Artifact Definitions
│   │   ├── 006-chain-of-authority.md   # Authority Hierarchy
│   │   ├── 007-glossary.md              # Terminology
│   │   ├── 008-future-vision.md         # Vision (if normative)
│   │   ├── 009-engineering-knowledge.md # Knowledge Definition
│   │   ├── 010-knowledge-derivation.md  # Derivation Rules
│   │   ├── 011-adoption-model.md        # Adoption Guidance
│   │   ├── 012-traceability.md          # Traceability Requirements
│   │   ├── 013-authority-resolution.md  # Resolution Process
│   │   └── 014-engineering-review-process.md # Review Process
│   │
│   ├── audit/                           # Audit Standards
│   │   ├── FOUNDATION_AUDIT.md          # Foundation Audit
│   │   ├── COMPLIANCE_AUDIT.md          # Compliance Audit
│   │   ├── AUDIT_SCORING.md             # Scoring Methodology
│   │   ├── AUDIT_MATURITY.md            # Maturity Levels
│   │   └── AUDIT_TEMPLATE.md            # Audit Template
│   │
│   ├── execution/                        # Execution References
│   │   ├── EXECUTION_LOOP.md            # Execution Loop
│   │   ├── SESSION_PROTOCOL.md          # Session Protocol
│   │   ├── AGENT_SPECIFICATION.md       # Executor Spec
│   │   └── REPORT_FORMAT.md             # Report Format
│   │
│   ├── glossary.md                       # Unified Glossary
│   └── templates/                        # Standard Templates
│       ├── audit-template.md              # Audit Template
│       ├── report-template.md             # Report Template
│       └── finding-template.md            # Finding Template
│
├── runtime/                             # Runtime State
│   ├── state/                           # Session State
│   │   ├── current-session.yaml          # Active session info
│   │   ├── session-stack.yaml           # Nested session info
│   │   └── checkpoint.yaml              # Session checkpoint
│   │
│   ├── logs/                            # Execution Logs
│   │   └── {YYYY-MM-DD}-{session-id}.log # Timestamped logs
│   │
│   └── temp/                            # Temporary Files
│       ├── assessment/                   # Assessment artifacts
│       ├── analysis/                     # Analysis artifacts
│       └── working/                      # Working files
│
├── reports/                              # Generated Reports
│   ├── sessions/                        # Runtime Session Reports
│   │   └── {YYYY-MM-DD}-{session-id}.md
│   │
│   ├── audits/                          # Audit Reports
│   │   ├── compliance/                 # Compliance Audits
│   │   │   └── {YYYY-MM-DD}-{audit-id}-compliance.md
│   │   └── foundation/                 # Foundation Audits
│   │       └── {YYYY-MM-DD}-{audit-id}-foundation.md
│   │
│   └── reviews/                        # Execution Reviews
│       └── {YYYY-MM-DD}-{review-id}.md
│
├── history/                             # Historical Records
│   ├── audit-history/                   # Audit Execution Records
│   │   ├── {dimension}-{YYYY-MM-DD}.json
│   │   └── {dimension}-{YYYY-MM-DD}.json
│   │
│   ├── session-history/                 # Session Records
│   │   └── sessions.csv                 # All sessions
│   │
│   └── sync-history/                    # Sync Records
│       └── updates.yaml                 # All sync events
│
├── cache/                               # Optional Cache
│   ├── artifacts/                       # Artifact Cache
│   │   ├── manifest.json               # Cached artifact manifest
│   │   └── relationships.json          # Cached relationships
│   │
│   └── scans/                           # Repository Scans
│       └── {YYYY-MM-DD}-scan.json      # Structure scans
│
├── .backup/                             # Sync Backups
│   └── backup-{version}-{timestamp}/    # Versioned backups
│       ├── standards/
│       ├── manifest.yaml
│       └── config.yaml
│
├── config.yaml                          # Runtime Configuration
├── manifest.yaml                        # Standard Manifest
└── README.md                           # Environment README
```

---

## Component Details

### 1. Standards Directory (.kdse/standards/)

**Purpose:** Contains the pinned KDSE normative documents required for execution.

**Normative Documents Included:**

| Category | Documents | Purpose |
|----------|-----------|---------|
| Foundation | 000-014-*.md | Core principles, models, definitions |
| Audit | FOUNDATION_AUDIT.md, COMPLIANCE_AUDIT.md, etc. | Audit standards |
| Execution | EXECUTION_LOOP.md, SESSION_PROTOCOL.md, etc. | Execution references |
| Templates | audit-template.md, report-template.md, etc. | Standardized formats |
| Glossary | glossary.md | Unified terminology |

**NOT Included (Informative Only):**
- Case studies
- Examples
- Research documents
- Evolution records
- External references

**Standards Subdirectory Structure:**

```
standards/
├── foundation/           # Foundation documents (pinned)
├── audit/              # Audit documents (pinned)
├── execution/          # Execution references (pinned)
├── glossary.md         # Unified glossary (pinned)
└── templates/          # Standard templates (pinned)
```

### 2. Runtime Directory (.kdse/runtime/)

**Purpose:** Contains transient state used during execution.

**Important:** This directory is cleared between sessions and may be excluded from version control.

```
runtime/
├── state/              # Session state (preserved between pauses)
├── logs/              # Execution logs (rotated)
└── temp/              # Temporary files (cleared after session)
```

**State Subdirectory:**

```
runtime/state/
├── current-session.yaml    # Active session information
│                           # - Session ID
│                           # - Current phase
│                           # - Progress
│                           # - Timestamp
│
├── session-stack.yaml      # Nested session info (if applicable)
│                           # - Parent session ID
│                           # - Session depth
│
└── checkpoint.yaml         # Session checkpoint (if supported)
                            # - State for resumption
                            # - Last checkpoint time
```

**Logs Subdirectory:**

```
runtime/logs/
└── {YYYY-MM-DD}-{session-id}.log
    ├── Initialization timestamp
    ├── Phase transitions
    ├── Audit results
    ├── Errors and warnings
    └── Session completion
```

**Temp Subdirectory:**

```
runtime/temp/
├── assessment/          # Assessment artifacts
│   ├── artifacts.json # Discovered artifacts
│   ├── relationships.json # Mapped relationships
│   └── findings.json   # Raw findings
│
├── analysis/           # Analysis artifacts
│   ├── gaps.json      # Identified gaps
│   ├── priorities.json # Priority calculations
│   └── recommendations.json # Raw recommendations
│
└── working/           # Working files
    ├── temp-*.md     # Temporary documents
    └── scratch.*      # Scratch files
```

### 3. Reports Directory (.kdse/reports/)

**Purpose:** Stores all generated reports with historical preservation.

**Critical:** Reports are NEVER overwritten. Each report has a unique timestamp and ID.

```
reports/
├── sessions/          # Runtime Session Reports
│   └── {YYYY-MM-DD}-{session-id}.md
│       ├── Header (session metadata)
│       ├── Current Status
│       ├── Audit Summary
│       ├── Findings
│       ├── Recommendations
│       └── Session History
│
├── audits/           # Audit Reports
│   ├── compliance/   # Compliance Audit Reports
│   │   └── {YYYY-MM-DD}-{audit-id}-compliance.md
│   │       ├── Audit Metadata
│   │       ├── Dimension Scores
│   │       ├── Findings by Severity
│   │       ├── Evidence
│   │       └── Recommendations
│   │
│   └── foundation/   # Foundation Audit Reports
│       └── {YYYY-MM-DD}-{audit-id}-foundation.md
│           ├── Audit Metadata
│           ├── Dimension Scores
│           ├── Findings
│           └── Recommendations
│
└── reviews/          # Execution Review Reports
    └── {YYYY-MM-DD}-{review-id}.md
        ├── Review Metadata
        ├── Assessment Summary
        ├── Recommendations
        └── Approval
```

### 4. History Directory (.kdse/history/)

**Purpose:** Maintains historical records for auditability and reproducibility.

```
history/
├── audit-history/     # Audit Execution Records (JSON)
│   ├── knowledge-{YYYY-MM-DD}.json
│   ├── architecture-{YYYY-MM-DD}.json
│   └── traceability-{YYYY-MM-DD}.json
│       ├── Audit timestamp
│       ├── Dimension evaluated
│       ├── Score assigned
│       ├── Findings summary
│       └── Evidence references
│
├── session-history/   # Session Records (CSV)
│   └── sessions.csv
│       ├── session_id
│       ├── start_time
│       ├── end_time
│       ├── status
│       ├── score_before
│       ├── score_after
│       └── recommendations
│
└── sync-history/     # Sync Records (YAML)
    └── updates.yaml
        ├── All sync events
        ├── Version history
        ├── Rollbacks
        └── Migration records
```

**Audit History JSON Format:**

```json
{
  "audit_id": "AUD-2026-07-10-001",
  "dimension": "knowledge",
  "timestamp": "2026-07-10T10:00:00Z",
  "repository_version": "1.0",
  "kdse_version": "1.3",
  "score": 7.5,
  "findings_count": 3,
  "critical_findings": 0,
  "high_findings": 1,
  "medium_findings": 2,
  "low_findings": 0,
  "evidence_files": ["e001", "e002", "e003"],
  "session_id": "SES-2026-07-10-001"
}
```

**Session History CSV Format:**

```csv
session_id,start_time,end_time,duration_minutes,status,score_before,score_after,delta,recommendations_count,approved_count,rejected_count
SES-2026-07-10-001,2026-07-10T10:00:00Z,2026-07-10T10:45:00Z,45,completed,5.2,5.8,0.6,3,2,1
```

### 5. Cache Directory (.kdse/cache/)

**Purpose:** Optionally stores reusable data for improved performance.

**Behavior:** Cache is regenerated from source if cleared. Never affects reproducibility.

```
cache/
├── artifacts/         # Artifact Cache
│   ├── manifest.json  # Cached artifact inventory
│   ├── relationships.json # Cached relationships
│   └── stewards.json  # Cached steward mappings
│
└── scans/            # Repository Scans
    └── {YYYY-MM-DD}-scan.json
        ├── Directory structure
        ├── File inventory
        └── Metadata
```

**Cacheable Items:**

| Item | Cacheable | Rationale |
|------|-----------|----------|
| Artifact manifests | Yes | Derived from source |
| Relationship mappings | Yes | Derived from source |
| Steward assignments | Yes | Derived from source |
| Directory structures | Yes | Derived from source |
| KDSE Standards | No | Use pinned version |
| Human-authored content | No | Always from source |
| Session-specific data | No | Transient |

### 6. Backup Directory (.kdse/.backup/)

**Purpose:** Contains backups created during sync operations.

**Important:** Backups are created automatically before sync and can be used for rollback.

```
backup/
└── backup-{version}-{timestamp}/
    ├── standards/          # Backed-up standards
    ├── manifest.yaml       # Backed-up manifest
    └── config.yaml        # Backed-up config
```

### 7. Configuration Files

#### manifest.yaml

Located at: `.kdse/manifest.yaml`

**Purpose:** Identifies the exact KDSE standard in use.

```yaml
kdse:
  version: "1.3"
  commit: "a1b2c3d4e5f6..."
  source: "https://github.com/tamzrod/KDSE"
  previous-version: null

repository:
  version: "1.0"
  initialized: "2026-07-10T10:00:00Z"
  last-sync: "2026-07-10T10:00:00Z"

profile:
  name: "default"
  scope: ["full"]

runtime:
  version: "1.0"
  mode: "standard"
```

#### config.yaml

Located at: `.kdse/config.yaml`

**Purpose:** Defines runtime behavior.

```yaml
runtime:
  mode: "standard"          # standard, minimal, verbose
  approval: "manual"          # manual, auto-approve, require-reason

audit:
  profile: "default"
  scope: ["full"]
  min-severity: "medium"

reports:
  location: "reports/"
  timestamp: true
  summary: true

logging:
  verbosity: "info"          # debug, info, warn, error
  location: "runtime/logs/"
  retention-days: 30

cache:
  enabled: true
  location: "cache/"
  max-size: 100             # MB
```

---

## File Naming Conventions

### Reports

```
{YYYY-MM-DD}-{type}-{id}.{format}

Examples:
2026-07-10-session-SES-001.md
2026-07-10-compliance-AUD-001.md
2026-07-10-foundation-AUD-002.md
```

### Logs

```
{YYYY-MM-DD}-{session-id}.log

Examples:
2026-07-10-SES-001.log
2026-07-10-SES-002.log
```

### History Records

```
{dimension}-{YYYY-MM-DD}.{format}

Examples:
knowledge-2026-07-10.json
architecture-2026-07-10.json
sessions.csv
updates.yaml
```

---

## Retention Policies

### Reports
- **Retained indefinitely**
- Never overwritten
- Never deleted automatically
- May be manually archived after X years

### Logs
- **Retained for 30 days** (configurable)
- Rotated by date
- Compressed after rotation
- Deleted after retention period

### Temp Files
- **Cleared after each session**
- Never committed to version control
- May be preserved for debugging if requested

### Cache
- **May be cleared at any time**
- Regenerated from source
- Maximum size enforced

### History
- **Retained indefinitely**
- Append-only
- Never modified after creation

### Backups
- **Retained for 90 days**
- One per sync operation
- Auto-deleted after retention

---

## Version Control Integration

### Recommended .gitignore

```gitignore
# KDSE Runtime Environment

# Runtime state (may be committed if desired)
# .kdse/runtime/state/
# .kdse/runtime/logs/

# Temporary files
.kdse/runtime/temp/

# Cache (optional)
.kdse/cache/

# Backups
.kdse/.backup/

# Core files (recommended to commit)
.kdse/standards/
.kdse/reports/
.kdse/history/
.kdse/config.yaml
.kdse/manifest.yaml
.kdse/README.md
```

### Alternative (Full Exclusion)

```gitignore
# KDSE Runtime Environment (full exclusion)
.kdse/
```

### Alternative (Commit Everything)

```gitignore
# No exclusion - commit everything
```

---

## Directory Purpose Summary

| Directory | Purpose | Committable | Critical |
|-----------|---------|-------------|----------|
| `.kdse/standards/` | Normative documents | Yes | Yes |
| `.kdse/runtime/state/` | Session state | Optional | No |
| `.kdse/runtime/logs/` | Execution logs | Optional | No |
| `.kdse/runtime/temp/` | Temporary files | No | No |
| `.kdse/reports/` | Generated reports | Yes | Yes |
| `.kdse/history/` | Historical records | Yes | Yes |
| `.kdse/cache/` | Optional cache | No | No |
| `.kdse/.backup/` | Sync backups | No | No |
| `.kdse/config.yaml` | Configuration | Yes | Yes |
| `.kdse/manifest.yaml` | Version manifest | Yes | Yes |
| `.kdse/README.md` | Documentation | Yes | No |

---

## Summary

The `.kdse/` directory structure provides a complete, self-contained KDSE Runtime Environment:

| Component | Purpose | Key Feature |
|-----------|---------|-------------|
| standards/ | Pinned normative documents | Version-pinned, offline access |
| runtime/ | Transient execution state | Cleared between sessions |
| reports/ | Generated documentation | Never overwritten, preserved |
| history/ | Audit trail | Append-only, indefinitely retained |
| cache/ | Performance optimization | Regenerated from source |
| .backup/ | Sync recovery | Enables rollback |
| config.yaml | Runtime behavior | Customizable |
| manifest.yaml | Environment identity | Version-pinned |

The layout ensures:
- **Determinism**: Exact versions, reproducible environments
- **Auditability**: Complete history, never deleted
- **Performance**: Optional caching, temp files cleared
- **Reliability**: Backups before changes, rollback capability
- **Simplicity**: Clear structure, obvious purpose

---

*This document defines the KDSE Runtime Environment layout. It is normative for KDSE-enabled repositories.*
