# KDSE Runtime Environment

**Document Version:** 1.0  
**Type:** Normative Runtime Specification  
**Effective Date:** 2026-07-10  

---

## Purpose

The KDSE Runtime Environment defines the local engineering environment that enables KDSE to execute consistently within any repository. The Runtime Environment transforms a repository into a self-contained, reproducible KDSE-enabled workspace.

---

## Design Philosophy

### The Git Analogy

Just as Git repositories contain `.git/` for version control, every KDSE-enabled repository contains `.kdse/` for engineering governance.

```
.git/                          .kdse/
├── objects/                   ├── standards/
├── refs/                     ├── runtime/
├── config                    ├── reports/
├── HEAD                      ├── history/
└── hooks/                    ├── cache/
                              ├── config.yaml
                              └── manifest.yaml
```

Both directories:
- Are hidden (prefixed with `.`)
- Contain repository-specific configuration
- Track versioned state
- Enable consistent local operations
- Do not require network connectivity for core operations

---

## Directory Structure

```
.kdse/
├── standards/                 # Pinned KDSE normative standards
│   ├── foundation/           # Core principles and definitions
│   ├── audit/               # Audit standards and templates
│   ├── execution/           # Execution model references
│   ├── glossary.md          # Terminology reference
│   └── templates/           # Audit and report templates
│
├── runtime/                 # Runtime-specific state
│   ├── state/               # Current session state
│   ├── logs/                # Execution logs
│   └── temp/                # Temporary runtime files
│
├── reports/                  # Generated reports (never overwritten)
│   ├── sessions/            # Runtime session reports
│   ├── audits/              # Compliance and foundation audits
│   └── reviews/             # Execution reviews
│
├── history/                  # Preserved historical data
│   ├── audit-history/        # Previous audit results
│   ├── session-history/      # Session execution records
│   └── sync-history/        # Standard update records
│
├── cache/                    # Optional cached data
│   └── artifacts/           # Cached artifact metadata
│
├── config.yaml              # Runtime configuration
├── manifest.yaml             # Standard version manifest
└── README.md                # Environment documentation
```

---

## Core Components

### 1. Standards Directory (.kdse/standards/)

The standards directory contains the pinned KDSE normative documents required for execution.

**Purpose:** Provides deterministic, offline-capable access to the engineering standard.

**Contents:**
- Foundation documents (principles, models, glossary)
- Audit standards (compliance, foundation, scoring)
- Execution references (execution model, session protocol)
- Templates (audit templates, report formats)

**NOT Included:**
- Case studies or examples
- Research documents
- Historical evolution records
- External references

**Rationale:** Only normative documents required for execution are included. Reference implementations and examples remain in the KDSE repository.

### 2. Runtime Directory (.kdse/runtime/)

The runtime directory contains transient state used during execution.

**Purpose:** Enables session continuity and logging.

**Contents:**
- `state/` - Current session state files
- `logs/` - Execution logs with timestamps
- `temp/` - Temporary working files

**Behavior:**
- Cleared between sessions
- Never committed to version control
- May be excluded via `.gitignore`

### 3. Reports Directory (.kdse/reports/)

The reports directory stores all generated reports.

**Purpose:** Preserves complete audit and session history.

**Structure:**
```
reports/
├── sessions/
│   └── YYYY-MM-DD-{session-id}.md
├── audits/
│   ├── compliance/
│   │   └── YYYY-MM-DD-{audit-id}.md
│   └── foundation/
│       └── YYYY-MM-DD-{audit-id}.md
└── reviews/
    └── YYYY-MM-DD-{review-id}.md
```

**Behavior:**
- Reports are never overwritten
- Each report has unique timestamp and ID
- Historical reports preserved indefinitely
- May be tracked or excluded per repository policy

### 4. History Directory (.kdse/history/)

The history directory maintains execution history for reproducibility.

**Purpose:** Enables audit trail and rollback capability.

**Structure:**
```
history/
├── audit-history/           # Audit execution records
│   └── {dimension}-{timestamp}.json
├── session-history/          # Session metadata
│   └── sessions.csv
└── sync-history/            # Standard update log
    └── updates.yaml
```

### 5. Cache Directory (.kdse/cache/)

The cache directory optionally stores reusable data.

**Purpose:** Improves performance for repeated operations.

**May Be Cached:**
- Artifact metadata (hashes, relationships)
- Audit intermediate results
- Repository structure scans
- Steward mappings

**Must Never Be Cached:**
- KDSE standard documents (use pinned versions)
- Human-authored content
- Session-specific data
- Credentials or secrets

**Behavior:**
- Cache may be regenerated from source
- Should not affect reproducibility
- May be cleared without data loss

---

## Configuration Files

### manifest.yaml

The manifest identifies the exact KDSE standard being used.

```yaml
# KDSE Manifest
# Generated: 2026-07-10T10:00:00Z

kdse:
  version: "1.3"              # KDSE standard version
  commit: "abc123..."         # Git commit of standard
  source: "github.com/..."    # KDSE repository URL

repository:
  version: "1.0"              # Repository's KDSE profile version
  initialized: "2026-07-10"   # When .kdse was created
  last-sync: "2026-07-10"     # Last standard update

profile:
  name: "default"             # Audit profile in use
  scope: ["full"]             # Dimensions included

runtime:
  version: "1.0"              # Runtime environment version
  mode: "standard"            # Execution mode
```

### config.yaml

The configuration defines runtime behavior.

```yaml
# KDSE Runtime Configuration

runtime:
  # Execution mode: standard, minimal, verbose
  mode: "standard"
  
  # Approval mode: manual, auto-approve, require-reason
  approval: "manual"

audit:
  # Default audit profile
  profile: "default"
  
  # Dimensions to include
  scope: ["full"]
  
  # Minimum severity for reporting
  min-severity: "medium"

reports:
  # Report location (relative to .kdse or absolute)
  location: "reports/"
  
  # Include timestamps in filenames
  timestamp: true
  
  # Generate summary reports
  summary: true

logging:
  # Log verbosity: debug, info, warn, error
  verbosity: "info"
  
  # Log location
  location: "runtime/logs/"
  
  # Retain logs for N days (0 = forever)
  retention-days: 30

cache:
  # Enable artifact caching
  enabled: true
  
  # Cache location
  location: "cache/"
  
  # Maximum cache size (MB)
  max-size: 100
```

---

## Design Principles

### 1. Deterministic Execution

The Runtime Environment ensures identical behavior across executions:
- Pinned standard version prevents drift
- Manifest captures exact dependencies
- Reports preserve historical state

### 2. Offline Capability

Core operations function without network access:
- Standards available locally
- Audits executable offline
- Reports generated locally

### 3. Version Pinning

Repositories control their own standard version:
- Not automatically updated
- Updates are intentional
- Migrations are documented

### 4. Reproducibility

Every session is reproducible:
- Manifest captures environment state
- Reports preserve findings
- History enables audit trail

### 5. Lightweight

The environment remains minimal:
- Only normative documents included
- Cache is optional
- Temporary files cleaned between sessions

---

## Runtime Discovery

When "Run KDSE" is executed, the Runtime:

1. **Discovers** `.kdse/` in current directory (or parent directories)
2. **Loads** manifest to identify standard version
3. **Validates** environment integrity
4. **Executes** assessment against pinned standard
5. **Reports** to `.kdse/reports/`
6. **Updates** history and state

---

## Environment Integrity

### Validation Checks

Before execution, the Runtime verifies:

| Check | Purpose |
|-------|---------|
| Manifest exists | Environment initialized |
| Standards present | Required documents available |
| Version compatible | Standard and runtime compatible |
| Config valid | Configuration is well-formed |
| Reports readable | Previous reports accessible |

### Recovery Actions

| Issue | Recovery |
|-------|---------|
| Missing manifest | Re-initialize required |
| Corrupt standards | Re-sync from source |
| Config invalid | Reset to defaults |
| Cache corrupted | Clear cache |

---

## Example Usage

### Initialize Repository

```bash
kdse init
# Creates .kdse/ with pinned standard
```

### Run Assessment

```bash
kdse run
# Discovers .kdse/, loads standard, runs assessment
```

### Generate Report

```bash
kdse report --latest
# Generates report from most recent assessment
```

### Update Standard

```bash
kdse sync --version 1.4
# Updates .kdse/standards/ to version 1.4
```

---

## Relationship to KDSE Repository

```
┌─────────────────────────────────────────────────────────────┐
│                    KDSE Repository                          │
│                                                             │
│  Contains:                                                  │
│  ├── docs/foundation/     (normative → copied to .kdse)   │
│  ├── docs/audit/          (normative → copied to .kdse)   │
│  ├── docs/execution/      (informative → NOT copied)       │
│  ├── docs/evolution/      (informative → NOT copied)       │
│  └── runtime/examples/     (informative → NOT copied)       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ kdse init / kdse sync
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   .kdse/standards/                          │
│                   (Local Repository)                        │
│                                                             │
│  Contains:                                                  │
│  ├── foundation/           (pinned normative)              │
│  ├── audit/               (pinned normative)               │
│  ├── glossary.md          (pinned normative)              │
│  └── templates/           (pinned normative)               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Versioning

| Component | Versioning | Policy |
|-----------|-----------|--------|
| Standard | Semantic | Pinned per repository |
| Runtime | Semantic | Per standard version |
| Config Schema | Semantic | Backward compatible |
| Reports | Timestamped | Never overwritten |

---

## Summary

The KDSE Runtime Environment (.kdse/) transforms any repository into a self-contained KDSE-enabled workspace:

| Feature | Benefit |
|---------|---------|
| Local standards | Offline execution |
| Version pinning | Reproducibility |
| Report history | Audit trail |
| Manifest | Environment identity |
| Configuration | Customizable behavior |
| Cache | Performance optimization |

The Runtime Environment makes KDSE as reliable and discoverable as Git—present in every repository, ready to execute, and independent of external dependencies.

---

*This document defines the KDSE Runtime Environment. It is normative for KDSE-enabled repositories.*
