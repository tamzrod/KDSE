# KDSE Runtime Environment

**Document Version:** 1.1  
**Type:** Normative Runtime Specification  
**Effective Date:** 2026-07-10  
**Supersedes:** Version 1.0

---

## Purpose

The KDSE Runtime Environment defines the local engineering environment that enables KDSE to execute consistently within any repository. The Runtime Environment transforms a repository into a self-contained, reproducible KDSE-enabled workspace.

This specification is **normative**. All implementations MUST follow this specification. Deviations MUST be documented and justified.

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
├── standards/                 # [MANDATORY] Pinned KDSE normative standards
│   ├── foundation/           # [MANDATORY] Core principles and definitions
│   ├── audit/               # [MANDATORY] Audit standards
│   └── glossary.md          # [MANDATORY] Terminology reference
│
├── runtime/                 # [MANDATORY] Runtime-specific state
│   ├── state/               # [MANDATORY] Current session state
│   ├── logs/                # [MANDATORY] Execution logs
│   └── temp/                # [RECOMMENDED] Temporary runtime files
│
├── reports/                  # [MANDATORY] Generated reports
│   ├── sessions/            # [MANDATORY] Runtime session reports
│   ├── audits/              # [MANDATORY] Compliance and foundation audits
│   │   ├── compliance/      # [MANDATORY] Compliance audit reports
│   │   └── foundation/      # [MANDATORY] Foundation audit reports
│   └── reviews/             # [MANDATORY] Execution reviews
│
├── history/                  # [MANDATORY] Preserved historical data
│   ├── audit-history/        # [MANDATORY] Previous audit results
│   ├── session-history/      # [MANDATORY] Session execution records
│   └── sync-history/        # [MANDATORY] Standard update records
│
├── cache/                    # [OPTIONAL] Cached data
│   └── artifacts/           # [OPTIONAL] Cached artifact metadata
│
├── config.yaml              # [MANDATORY] Runtime configuration
├── manifest.yaml             # [MANDATORY] Standard version manifest
└── README.md                # [MANDATORY] Environment documentation
```

### Directory Classification

| Directory | Status | Required for Valid Environment |
|-----------|--------|-------------------------------|
| `.kdse/standards/` | MANDATORY | YES |
| `.kdse/standards/foundation/` | MANDATORY | YES |
| `.kdse/standards/audit/` | MANDATORY | YES |
| `.kdse/standards/glossary.md` | MANDATORY | YES |
| `.kdse/runtime/` | MANDATORY | YES |
| `.kdse/runtime/state/` | MANDATORY | YES |
| `.kdse/runtime/logs/` | MANDATORY | YES |
| `.kdse/runtime/temp/` | RECOMMENDED | NO |
| `.kdse/reports/` | MANDATORY | YES |
| `.kdse/reports/sessions/` | MANDATORY | YES |
| `.kdse/reports/audits/` | MANDATORY | YES |
| `.kdse/reports/audits/compliance/` | MANDATORY | YES |
| `.kdse/reports/audits/foundation/` | MANDATORY | YES |
| `.kdse/reports/reviews/` | MANDATORY | YES |
| `.kdse/history/` | MANDATORY | YES |
| `.kdse/history/audit-history/` | MANDATORY | YES |
| `.kdse/history/session-history/` | MANDATORY | YES |
| `.kdse/history/sync-history/` | MANDATORY | YES |
| `.kdse/cache/` | OPTIONAL | NO |
| `.kdse/config.yaml` | MANDATORY | YES |
| `.kdse/manifest.yaml` | MANDATORY | YES |
| `.kdse/README.md` | MANDATORY | YES |

### Status Definitions

| Status | Meaning | Created During Init | Required for Runtime |
|--------|---------|--------------------|--------------------|
| MANDATORY | Must exist for valid environment | YES | YES |
| RECOMMENDED | Should exist, warnings if missing | YES | YES (with warnings) |
| OPTIONAL | May exist, no warnings if missing | NO | N/A |

---

## File Naming Conventions

### Report Files

```
{YYYY-MM-DD}-{type}-{id}.md

Components:
- YYYY-MM-DD: ISO 8601 date
- type: session | compliance | foundation | review
- id: Sequential identifier (001, 002, ...)

Examples:
- 2026-07-10-session-SES-001.md
- 2026-07-10-compliance-AUD-001.md
- 2026-07-10-foundation-AUD-002.md
```

### Log Files

```
{YYYY-MM-DD}-{session-id}.log

Examples:
- 2026-07-10-SES-001.log
- 2026-07-10-SES-002.log
```

### History Records

```
{dimension}-{YYYY-MM-DD}.json

Examples:
- knowledge-2026-07-10.json
- architecture-2026-07-10.json
```

---

---

## Normative Document Manifest

KDSE defines a canonical manifest of all normative documents. This manifest is the authoritative source for what MUST be installed during initialization.

### Canonical Normative Documents

#### Foundation Documents (docs/foundation/)

| ID | Filename | Mandatory | Purpose |
|----|----------|-----------|---------|
| F-001 | 000-what-is-kdse.md | YES | Canonical KDSE definition |
| F-002 | 001-why-kdse-exists.md | YES | Problem statement and rationale |
| F-003 | 002-scope.md | YES | Scope boundaries and limitations |
| F-004 | 003-core-principles.md | YES | Core principles of KDSE |
| F-005 | 004-engineering-model.md | YES | Engineering lifecycle model |
| F-006 | 005-engineering-artifacts.md | YES | Artifact type definitions |
| F-007 | 006-chain-of-authority.md | YES | Authority hierarchy |
| F-008 | 007-glossary.md | YES | Terminology definitions |
| F-009 | 008-future-vision.md | NO | Vision for future evolution |
| F-010 | 009-engineering-knowledge.md | YES | Knowledge artifact definition |
| F-011 | 010-knowledge-derivation.md | YES | Derivation rules |
| F-012 | 011-adoption-model.md | NO | Adoption guidance |
| F-013 | 012-traceability.md | YES | Traceability requirements |
| F-014 | 013-authority-resolution.md | YES | Resolution process |
| F-015 | 014-engineering-review-process.md | YES | Review process definition |

#### Audit Documents (docs/audit/)

| ID | Filename | Mandatory | Purpose |
|----|----------|-----------|---------|
| A-001 | FOUNDATION_AUDIT.md | YES | Foundation audit standard |
| A-002 | COMPLIANCE_AUDIT.md | YES | Compliance audit standard |
| A-003 | AUDIT_SCORING.md | YES | Scoring methodology |
| A-004 | AUDIT_MATURITY.md | YES | Maturity level definitions |
| A-005 | AUDIT_TEMPLATE.md | YES | Audit report template |
| A-006 | AUDIT_README.md | NO | Audit system documentation |

#### Glossary

| ID | Filename | Mandatory | Purpose |
|----|----------|-----------|---------|
| G-001 | 007-glossary.md | YES | Consolidated glossary (same as F-008) |

### Informative Documents (NOT Installed)

The following documents are informative and MUST NOT be installed to `.kdse/`:

| Source | Documents | Reason |
|--------|-----------|--------|
| docs/evolution/ | All documents | Historical evidence, not required |
| docs/execution/ | All documents | Reference implementations only |
| runtime/ | All documents | Example implementations only |
| docs/audit/ | KDSE_EXECUTION_MODEL_REVIEW.md | Audit reports, not standards |
| docs/audit/ | KDSE_FOUNDATION_AUDIT.md | Audit reports, not standards |

### Installation Requirement

During `kdse init`, implementations MUST install:
- All documents marked "Mandatory: YES" in the tables above
- Documents marked "Mandatory: NO" MAY be installed at implementer's discretion

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

### Discovery Algorithm

When "Run KDSE" is executed, the Runtime MUST follow this discovery algorithm:

```
1. START at current working directory

2. SEARCH for .kdse/:
   a. Check if .kdse/ exists in current directory
   b. If found and valid → USE THIS .kdse/
   c. If not found → CONTINUE to step 3

3. TRAVERSE PARENT DIRECTORIES (maximum 3 levels):
   a. For each parent directory (moving upward):
      - Check if .kdse/ exists
      - If found and valid → USE THIS .kdse/
      - If found but invalid → STOP and ERROR
   b. If no .kdse/ found after 3 levels → CONTINUE to step 4

4. HOME DIRECTORY CHECK:
   a. Check if .kdse/ exists in home directory (~/)
   b. If found and valid → WARN and USE HOME .kdse/
   c. If not found → CONTINUE to step 5

5. NO ENVIRONMENT FOUND:
   a. Return error: "No .kdse/ found"
   b. Suggest: "Run 'kdse init' to initialize"
```

### Discovery Validity

A `.kdse/` directory is valid when:
- `manifest.yaml` exists and is readable
- `standards/` directory exists
- At least one mandatory standard document exists

### Discovery Priority

When multiple `.kdse/` directories are found:
1. **Nearest to target wins**: Use the one closest to the current directory
2. **Current directory takes precedence** over parent directories
3. **Parent directories** checked in order (closest first)
4. **Home directory** checked last, with warning

### Discovery Scope

- Maximum parent traversal: 3 levels
- Starting point: Current working directory
- Search order: Current → Parent1 → Parent2 → Parent3 → Home

### Symlink Handling

- If `.kdse/` is a symlink, resolve to real path before validation
- Circular symlinks MUST be detected and rejected with error

### Nested Repository Handling

- If inside a git submodule, search within submodule only
- Parent repository's `.kdse/` is NOT automatically inherited
- Each repository requiring KDSE MUST have its own `.kdse/`

---

## Runtime Discovery Execution Sequence

When "Run KDSE" is executed, the Runtime:

1. **Discovers** `.kdse/` using the discovery algorithm above
2. **Loads** manifest.yaml to identify standard version
3. **Validates** environment integrity
4. **Executes** assessment against pinned standard
5. **Reports** to `.kdse/reports/`
6. **Updates** history and state

---

## Environment Integrity

### Integrity Verification

The Runtime MUST verify environment integrity before execution using these checks:

#### Mandatory Validation Checks

These checks are MANDATORY. Failure of any mandatory check BLOCKS execution:

| Check | Severity | Description | Failure Action |
|-------|----------|-------------|---------------|
| Manifest exists | BLOCKING | manifest.yaml must exist | Error: "Re-initialize required" |
| Manifest valid YAML | BLOCKING | manifest.yaml must parse | Error: "Manifest corrupted" |
| Standards directory exists | BLOCKING | standards/ must exist | Error: "Standards missing" |
| All mandatory standards present | BLOCKING | All F-001 to F-015 and A-001 to A-005 must exist | Error: "Standards incomplete" |
| At least one standard readable | BLOCKING | Standards must be readable | Error: "Standards unreadable" |

#### Recommended Validation Checks

These checks are RECOMMENDED. Failure generates a WARNING but does not block execution:

| Check | Severity | Description | Failure Action |
|-------|----------|-------------|---------------|
| Config exists | WARNING | config.yaml should exist | Warning: "Using defaults" |
| Config valid YAML | WARNING | config.yaml must parse | Warning: "Using defaults" |
| Directory structure complete | WARNING | All mandatory directories should exist | Warning: "Some directories missing" |
| Git repository detected | WARNING | Repository should be git-tracked | Warning: "Not a git repository" |

#### Version Compatibility Check

The Runtime MUST verify version compatibility between manifest and installed standards:

```
1. Read kdse.version from manifest.yaml
2. Verify version exists in KDSE version registry
3. Check runtime.version compatibility:
   - If runtime.version > required -> WARN: "Runtime may be outdated"
   - If runtime.version < minimum -> ERROR: "Runtime version incompatible"
4. Verify all standards match the declared version
```

### Recovery Actions

| Issue | Severity | Recovery Action |
|-------|----------|-----------------|
| Missing manifest | BLOCKING | Re-initialize required |
| Manifest invalid | BLOCKING | Delete manifest, re-initialize |
| Standards missing | BLOCKING | Run `kdse sync` |
| Standards corrupted | BLOCKING | Run `kdse sync --force` |
| Config invalid | WARNING | Reset to defaults |
| Cache corrupted | WARNING | Clear cache |
| Directory missing | WARNING | Create missing directory |

### Integrity Verification Procedure

Before any Runtime operation:

```
1. DISCOVER .kdse/ using discovery algorithm
2. VALIDATE manifest:
   a. Check manifest.yaml exists
   b. Parse manifest.yaml
   c. Verify required fields present
   d. Verify version compatibility
3. VALIDATE standards:
   a. Check standards/ directory exists
   b. Verify all mandatory documents present
   c. Verify at least one document readable
4. VALIDATE structure:
   a. Verify mandatory directories exist
   b. Warn about missing optional directories
5. REPORT status:
   a. If blocking errors -> FAIL with error details
   b. If warnings only -> PROCEED with warnings
```


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
