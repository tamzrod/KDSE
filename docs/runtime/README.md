# KDSE Runtime Documentation

**Document Version:** 1.0  
**Effective Date:** 2026-07-10  

---

## Overview

This directory contains the normative specifications for the KDSE Runtime Environment (.kdse/). The Runtime Environment transforms any repository into a self-contained, reproducible KDSE-enabled workspace.

---

## Documents

| Document | Purpose |
|----------|---------|
| [KDSE_RUNTIME_ENVIRONMENT.md](KDSE_RUNTIME_ENVIRONMENT.md) | Overview and design philosophy of the .kdse/ directory |
| [KDSE_INITIALIZATION.md](KDSE_INITIALIZATION.md) | Process for creating and configuring .kdse/ |
| [KDSE_STANDARD_SYNC.md](KDSE_STANDARD_SYNC.md) | Process for updating pinned KDSE standards |
| [KDSE_RUNTIME_LAYOUT.md](KDSE_RUNTIME_LAYOUT.md) | Complete directory structure and specifications |

---

## Quick Reference

### Directory Structure

```
.kdse/
├── standards/              # Pinned KDSE normative standards
├── runtime/                # Transient execution state
├── reports/                # Generated reports (preserved)
├── history/                # Historical records (preserved)
├── cache/                  # Optional cache
├── config.yaml            # Runtime configuration
├── manifest.yaml           # Version manifest
└── README.md              # This file
```

### Key Concepts

| Concept | Description |
|---------|-------------|
| Version Pinning | Repositories pin exact KDSE version for reproducibility |
| Offline Execution | Standards available locally, no network required |
| Report Preservation | Reports never overwritten, unique timestamps |
| Sync Control | Updates are manual, not automatic |

### Common Commands

```bash
# Initialize repository
kdse init

# Run assessment
kdse run

# Generate report
kdse report --latest

# Check for updates
kdse sync --check

# Sync to new version
kdse sync --version 1.4

# Rollback
kdse sync --rollback

# Verify environment
kdse status
```

---

## Relationship to KDSE Repository

```
KDSE Repository                    .kdse/ (Repository)
     │                                   │
     ├── docs/foundation/ ──────────────► ├── standards/foundation/
     ├── docs/audit/ ───────────────────► ├── standards/audit/
     ├── docs/execution/ ───────────────► ├── standards/execution/
     └── docs/execution/templates/ ──────► ├── standards/templates/
                                               │
                                               ├── runtime/
                                               ├── reports/
                                               ├── history/
                                               ├── config.yaml
                                               └── manifest.yaml
```

Only **normative** documents are copied to `.kdse/`. Reference implementations and examples remain in the KDSE repository.

---

## Version Information

| Component | Version | Description |
|-----------|--------|-------------|
| Runtime Environment | 1.0 | Initial specification |
| Standards Format | 1.0 | Compatible with KDSE 1.x |
| Config Schema | 1.0 | Initial schema |

---

## Normative vs Informative

All documents in this directory are **normative**. They define the required structure and behavior of the KDSE Runtime Environment.

Reference implementations and examples remain in the main KDSE repository under `runtime/`.

---

*This directory contains the normative specifications for KDSE Runtime Environments.*
