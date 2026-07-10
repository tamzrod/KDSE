# KDSE Initialization

**Document Version:** 1.0  
**Type:** Normative Runtime Specification  
**Effective Date:** 2026-07-10  

---

## Purpose

KDSE Initialization is the process that transforms a repository into a KDSE-enabled workspace by creating the `.kdse/` directory structure and installing the required standards.

---

## Initialization Overview

### What Initialization Does

1. Creates the `.kdse/` directory structure
2. Installs the KDSE normative standards
3. Generates the manifest identifying the installed standard
4. Creates the initial configuration
5. Verifies the installation integrity

### What Initialization Does NOT Do

- Does not modify source code
- Does not create artifacts in the repository
- Does not perform any assessment
- Does not require network after initial setup

---

## Initialization Command

### Command Interface

```bash
kdse init [OPTIONS]

OPTIONS:
  --source URL        KDSE source repository (default: official KDSE)
  --version VERSION   KDSE version to install (default: latest stable)
  --profile PROFILE   Audit profile to use (default: default)
  --no-verify        Skip post-install verification
  --force           Overwrite existing .kdse/ (requires confirmation)
```

### Examples

```bash
# Initialize with latest stable version
kdse init

# Initialize with specific version
kdse init --version 1.3

# Initialize from custom source
kdse init --source https://github.com/myorg/kdse --version 2.0

# Initialize with specific profile
kdse init --profile security-focused
```

---

## Initialization Process

### Phase 1: Discovery

```
┌─────────────────────────────────────────────────────────────┐
│ 1. DISCOVER                                                 │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ Check current directory:                                      │
│   ├── .kdse/ exists?                                        │
│   │     └── YES → Error: Already initialized                 │
│   │     └── NO → Continue                                   │
│   │                                                          │
│ Check parent directories:                                     │
│   ├── .kdse/ found in parent?                              │
│   │     └── YES → Offer to use parent OR create local       │
│   │     └── NO → Continue                                  │
│   │                                                          │
│ Check for .git/:                                             │
│   ├── .git/ exists?                                         │
│   │     └── NO → Warning: Not a git repository              │
│   │     └── YES → Continue                                  │
│   │                                                          │
└─────────────────────────────────────────────────────────────┘
```

**Output:**
- Confirms initialization is appropriate
- Identifies any existing .kdse/ environments
- Warns about non-git repositories

### Phase 2: Preparation

```
┌─────────────────────────────────────────────────────────────┐
│ 2. PREPARE                                                  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ Determine source:                                             │
│   ├── --source flag provided?                                │
│   │     └── YES → Use specified source                      │
│   │     └── NO → Use default KDSE repository               │
│   │                                                          │
│ Determine version:                                            │
│   ├── --version flag provided?                               │
│   │     └── YES → Use specified version                     │
│   │     └── NO → Fetch latest stable version                │
│   │                                                          │
│ Check network connectivity:                                   │
│   └── Required for initial installation only                │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Output:**
- Source repository identified
- Version determined
- Network access verified (if needed)

### Phase 3: Directory Creation

```
┌─────────────────────────────────────────────────────────────┐
│ 3. CREATE DIRECTORIES                                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ Create .kdse/                                                │
│   ├── .kdse/standards/                                      │
│   │     ├── .kdse/standards/foundation/                     │
│   │     ├── .kdse/standards/audit/                          │
│   │     ├── .kdse/standards/execution/                       │
│   │     ├── .kdse/standards/glossary.md                      │
│   │     └── .kdse/standards/templates/                       │
│   ├── .kdse/runtime/                                        │
│   │     ├── .kdse/runtime/state/                             │
│   │     ├── .kdse/runtime/logs/                             │
│   │     └── .kdse/runtime/temp/                            │
│   ├── .kdse/reports/                                       │
│   │     ├── .kdse/reports/sessions/                        │
│   │     ├── .kdse/reports/audits/                          │
│   │     └── .kdse/reports/reviews/                         │
│   ├── .kdse/history/                                        │
│   │     ├── .kdse/history/audit-history/                   │
│   │     ├── .kdse/history/session-history/                 │
│   │     └── .kdse/history/sync-history/                    │
│   ├── .kdse/cache/                                          │
│   │     └── .kdse/cache/artifacts/                          │
│   └── .kdse/README.md                                       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Output:**
- Complete directory structure created
- Permissions set appropriately
- Gitignore patterns recommended

### Phase 4: Standards Installation

```
┌─────────────────────────────────────────────────────────────┐
│ 4. INSTALL STANDARDS                                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ Fetch normative documents from source:                        │
│   ├── docs/foundation/                                     │
│   │     └── Copy to .kdse/standards/foundation/            │
│   ├── docs/audit/                                           │
│   │     └── Copy to .kdse/standards/audit/                 │
│   ├── docs/execution/                                       │
│   │     └── Copy to .kdse/standards/execution/              │
│   ├── docs/foundation/007-glossary.md                       │
│   │     └── Copy to .kdse/standards/glossary.md             │
│   └── templates/                                             │
│         └── Copy to .kdse/standards/templates/               │
│                                                              │
│ Exclude (NOT copied):                                        │
│   ├── docs/evolution/                                       │
│   ├── docs/execution/examples/                              │
│   └── runtime/examples/                                     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Standards Installed:**

| Category | Documents | Purpose |
|----------|-----------|---------|
| Foundation | 000-014-*.md | Core principles, models, definitions |
| Audit | FOUNDATION_AUDIT.md, COMPLIANCE_AUDIT.md, AUDIT_SCORING.md, AUDIT_MATURITY.md, AUDIT_TEMPLATE.md | Audit standards |
| Execution | EXECUTION_LOOP.md, SESSION_PROTOCOL.md, AGENT_SPECIFICATION.md, REPORT_FORMAT.md | Execution references |
| Templates | Audit templates, Report templates | Standardized formats |
| Glossary | glossary.md | Terminology reference |

### Phase 5: Manifest Generation

```
┌─────────────────────────────────────────────────────────────┐
│ 5. GENERATE MANIFEST                                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ Create .kdse/manifest.yaml:                                 │
│   ├── kdse.version: "X.Y"                                 │
│   ├── kdse.commit: "..."                                   │
│   ├── kdse.source: "URL"                                   │
│   ├── repository.version: "1.0"                           │
│   ├── repository.initialized: "ISO-8601"                   │
│   ├── repository.last-sync: "ISO-8601"                     │
│   ├── profile.name: "default"                              │
│   ├── runtime.version: "1.0"                               │
│   └── runtime.mode: "standard"                             │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Example manifest.yaml:**
```yaml
kdse:
  version: "1.3"
  commit: "a1b2c3d4e5f6..."
  source: "https://github.com/tamzrod/KDSE"

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

### Phase 6: Configuration Generation

```
┌─────────────────────────────────────────────────────────────┐
│ 6. GENERATE CONFIGURATION                                    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ Create .kdse/config.yaml:                                   │
│   ├── runtime.mode: "standard"                             │
│   ├── runtime.approval: "manual"                            │
│   ├── audit.profile: "default"                             │
│   ├── audit.scope: ["full"]                                │
│   ├── audit.min-severity: "medium"                        │
│   ├── reports.location: "reports/"                         │
│   ├── reports.timestamp: true                               │
│   ├── reports.summary: true                                 │
│   ├── logging.verbosity: "info"                            │
│   ├── logging.location: "runtime/logs/"                     │
│   ├── logging.retention-days: 30                           │
│   ├── cache.enabled: true                                  │
│   ├── cache.location: "cache/"                             │
│   └── cache.max-size: 100                                  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Phase 7: Documentation

```
┌─────────────────────────────────────────────────────────────┐
│ 7. CREATE DOCUMENTATION                                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ Create .kdse/README.md:                                      │
│   ├── What is .kdse/                                       │
│   ├── Directory structure overview                          │
│   ├── Quick start guide                                     │
│   ├── Common commands                                       │
│   └── Troubleshooting                                        │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Phase 8: Verification

```
┌─────────────────────────────────────────────────────────────┐
│ 8. VERIFY INSTALLATION                                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ Run integrity checks:                                        │
│   ├── Directory structure complete?                          │
│   ├── All standards installed?                              │
│   ├── Manifest valid?                                       │
│   ├── Configuration valid?                                   │
│   └── Standards readable?                                   │
│                                                              │
│ Verify standards integrity:                                  │
│   ├── Foundation documents present?                          │
│   ├── Audit documents present?                               │
│   ├── Cross-references valid?                               │
│   └── No corruption detected?                               │
│                                                              │
│ Generate verification report:                                 │
│   └── Display summary to user                               │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Verification Output:**
```
KDSE Initialization Complete
============================

Source:     github.com/tamzrod/KDSE
Version:    1.3
Commit:     a1b2c3d4e5f6...
Installed:  2026-07-10T10:00:00Z

Standards Installed:
  Foundation: 15 documents
  Audit: 5 documents
  Execution: 4 documents
  Templates: 3 documents

Status: ✓ Ready

Run 'kdse status' to verify environment
Run 'kdse run' to start your first session
```

---

## Post-Initialization Steps

### Recommended Git Configuration

Add `.kdse/` to `.gitignore`:

```gitignore
# KDSE Runtime Environment
.kdse/
```

**Note:** The `.kdse/` directory may optionally be committed to track environment history. This decision is repository-specific.

### First Session

After initialization:

```bash
# Verify environment
kdse status

# Run first assessment
kdse run

# Generate first report
kdse report --latest
```

---

## Re-Initialization

### When Required

Re-initialization may be required when:
- `.kdse/` is corrupted or missing
- Switching to a different KDSE source
- Upgrading to a new major version

### Force Re-Initialization

```bash
kdse init --force
```

**Warning:** Force re-initialization will:
- Overwrite existing `.kdse/` directory
- Clear all local configuration
- NOT delete reports or history (if using --force-with-history)
- Require confirmation before proceeding

### Selective Re-Initialization

```bash
# Reinstall standards only
kdse init --standards-only

# Regenerate manifest only
kdse init --manifest-only

# Reset configuration only
kdse init --config-only
```

---

## Initialization Profiles

### Default Profile

Standard installation with all normative documents.

### Minimal Profile

Installs only essential documents:
- Core principles
- Compliance audit
- Audit scoring
- Glossary

```bash
kdse init --profile minimal
```

### Custom Profile

Create a custom profile in `kdse.profiles/`:
```bash
mkdir -p kdse.profiles/my-profile
kdse init --profile my-profile
```

---

## Troubleshooting

### Network Errors

| Error | Solution |
|-------|----------|
| Cannot reach source | Check network connectivity |
| Source not found | Verify source URL |
| Version not found | Check version exists |

### Integrity Errors

| Error | Solution |
|-------|----------|
| Checksum mismatch | Re-download standards |
| Missing documents | Re-run initialization |
| Corrupt manifest | Delete manifest, regenerate |

### Permission Errors

| Error | Solution |
|-------|----------|
| Cannot create directory | Check write permissions |
| Cannot read standards | Check file permissions |

---

## Initialization Acceptance Criteria

The initialization process is considered successful when ALL mandatory acceptance criteria are met.

### Mandatory Acceptance Criteria

These criteria MUST be met for initialization to succeed:

| ID | Criterion | Verification Method | Failure Action |
|----|-----------|-------------------|--------------|
| AC-01 | `.kdse/` directory created | Check directory exists | Create directory |
| AC-02 | `manifest.yaml` created | Check file exists | Create manifest |
| AC-03 | `manifest.yaml` valid YAML | Parse YAML | Regenerate manifest |
| AC-04 | `standards/` directory created | Check directory exists | Create directory |
| AC-05 | All mandatory foundation documents installed | Check F-001 to F-015 | Install missing |
| AC-06 | All mandatory audit documents installed | Check A-001 to A-005 | Install missing |
| AC-07 | At least one standard document readable | Open and read file | Re-download standards |
| AC-08 | `config.yaml` created with valid YAML | Parse YAML | Regenerate config |
| AC-09 | All mandatory subdirectories created | Check structure | Create missing |

### Verification Procedure

After Phase 8 (Verification), the implementation MUST verify:

```
FOR EACH mandatory acceptance criterion:
    1. Execute verification check
    2. IF check fails:
       a. Log failure with criterion ID
       b. Execute recovery action
       c. Re-verify
    3. IF check still fails:
       a. Report criterion ID and failure reason
       b. STOP initialization with error
    4. IF all mandatory criteria pass:
       a. Report success
       b. Proceed to completion
```

### Success Criteria Summary

Initialization is successful when:

- [ ] All mandatory directories created
- [ ] All mandatory documents installed
- [ ] manifest.yaml valid and complete
- [ ] config.yaml valid and complete
- [ ] No blocking errors

### Error Reporting

When initialization fails, the report MUST include:

```
Initialization Failed
====================

Failed Criterion: AC-XX [Criterion Name]
Failure Reason: [Specific reason]
Recovery Attempted: [What was tried]
Recovery Result: [Success/Failure]

Suggestion: [How to resolve]
```

---
## Summary

KDSE Initialization transforms any repository into a KDSE-enabled workspace:

| Phase | Action | Output |
|-------|--------|--------|
| 1. Discover | Check environment | Validation status |
| 2. Prepare | Determine source/version | Source, version |
| 3. Create | Build directory structure | `.kdse/` |
| 4. Install | Copy standards | Normative documents |
| 5. Generate | Create manifest | manifest.yaml |
| 6. Configure | Create config | config.yaml |
| 7. Document | Create README | README.md |
| 8. Verify | Check integrity | Verification report |

After initialization, the repository is ready for KDSE execution without requiring network access.

---

*This document defines the KDSE initialization process. It is normative for KDSE-enabled repositories.*
