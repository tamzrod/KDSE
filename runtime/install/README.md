# KDSE Runtime Installation Framework

**Document Version:** 1.0  
**Type:** Runtime Infrastructure  
**Effective Date:** 2026-07-10

---

## Purpose

This framework provides deterministic scripts for installing, synchronizing, verifying, and uninstalling the KDSE Runtime. The scripts perform only filesystem operations and do not contain any AI logic, audit logic, or engineering decisions.

---

## Overview

The installation framework manages the KDSE Runtime lifecycle:

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Install   │───▶│    Sync     │───▶│   Verify    │───▶│  Uninstall  │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

**Install** initializes the Runtime in a repository.  
**Sync** updates the Runtime with latest standards.  
**Verify** confirms installation integrity.  
**Uninstall** removes the Runtime cleanly.

---

## Directory Structure

```
runtime/install/
├── README.md          # This file
├── common.sh          # Shared utilities and functions
├── install.sh         # Initialize KDSE Runtime
├── sync.sh            # Synchronize with repository
├── verify.sh          # Verify installation integrity
└── uninstall.sh       # Remove KDSE Runtime
```

**Installed Structure:**

```
.kdse/                     # Installation root
├── manifest.json          # Installation metadata
├── config.sh              # Runtime configuration
├── reports/               # Session reports
├── history/               # Session history
├── runtime/               # Runtime state
├── cache/                 # Cached data
├── configuration/         # User configuration
└── standards/
    ├── normative/         # Core KDSE documents
    └── informative/        # Reference documents
```

---

## Supported Installation Formats

The Runtime supports multiple manifest formats for backward compatibility:

| Format | File | Status | Notes |
|--------|------|--------|-------|
| JSON | `manifest.json` | Current | Primary format |
| YAML | `manifest.yaml` | Legacy | Auto-migrated to JSON |

### Automatic Migration

When the Runtime detects a legacy YAML manifest, it automatically migrates to the current JSON format:

1. **Detection**: Identifies legacy YAML manifest
2. **Backup**: Creates `manifest.yaml.backup` for safety
3. **Migration**: Converts to JSON while preserving metadata
4. **Verification**: Confirms successful migration

The migration process:
- Preserves all user data (reports, history, runtime, cache, configuration)
- Is idempotent (safe to run multiple times)
- Requires no operator intervention
- Records migration metadata in the manifest

---

## Scripts

### install.sh

Initializes KDSE Runtime in a repository.

**Responsibilities:**
- Create `.kdse/` directory structure
- Clone KDSE repository
- Install normative documents
- Install informative documents
- Generate manifest
- Generate configuration
- Verify installation

**Usage:**
```bash
./install.sh                    # Interactive install
./install.sh -v                 # Verbose output
./install.sh -f                 # Force reinstall
./install.sh -r <repo-url>     # Specify repository
./install.sh -b <branch>        # Specify branch
./install.sh -p <path>          # Installation path
```

### sync.sh

Synchronizes installed standards with KDSE repository.

**Responsibilities:**
- Read installed manifest
- Compare with repository version
- Synchronize normative documents
- Preserve user data (reports, history, runtime, cache, configuration)
- Update manifest
- Verify integrity

**Usage:**
```bash
./sync.sh                       # Sync with latest
./sync.sh -v                    # Verbose output
./sync.sh -f                    # Force sync even if up-to-date
```

**Preserved Data:**
| Directory | Description |
|----------|-------------|
| `reports/` | Session reports |
| `history/` | Session history |
| `runtime/` | Runtime state |
| `cache/` | Cached data |
| `configuration/` | User configuration |

### verify.sh

Verifies installation integrity.

**Checks Performed:**
1. Directory layout
2. Mandatory documents
3. Manifest integrity
4. Configuration validity
5. Version consistency
6. Preserved directories
7. Standards structure

**Usage:**
```bash
./verify.sh                     # Full verification
./verify.sh -v                 # Verbose output
./verify.sh -q                 # Quiet mode (PASS/FAIL only)
./verify.sh -j                 # JSON output
```

**Exit Codes:**
| Code | Status |
|------|--------|
| 0 | PASS - All checks passed |
| 1 | FAIL - One or more checks failed |
| 2 | ERROR - Cannot verify (not installed) |

### uninstall.sh

Removes KDSE Runtime installation.

**Responsibilities:**
- Confirm uninstallation (unless forced)
- Preserve user data by default
- Remove installation files
- Clean up old backups

**Usage:**
```bash
./uninstall.sh                  # Interactive uninstall
./uninstall.sh -f              # Force uninstall
./uninstall.sh -k              # Keep reports (default)
./uninstall.sh -a              # Keep all user data
```

**Preserved by Default:**
- `reports/` - Session reports
- `history/` - Session history

**Removed:**
- KDSE Runtime
- Normative documents
- Informative documents
- Git repository
- Cache

---

## Design Principles

### Deterministic

Same inputs always produce same outputs. No random behavior.

### Idempotent

Running multiple times produces same result. Safe to re-run.

### Non-Interactive

Scripts complete without user input unless required (e.g., uninstall confirmation).

### Verbose on Failure

Quiet during success, informative on failure.

### No AI Logic

Scripts contain only filesystem operations. No reasoning, no LLM calls.

### No Audit Logic

Audit functionality is handled by KDSE Runtime, not these scripts.

### No Engineering Decisions

Scripts only manage files. All decisions are made elsewhere.

---

## Supported Platforms

| Platform | Status |
|----------|--------|
| Linux | ✅ Supported |
| macOS | ✅ Supported |

---

## Requirements

### System Dependencies

Scripts require only basic Unix utilities:
- `bash` (4.0+)
- `mkdir`, `cp`, `rm`, `cat`, `grep`, `sed`
- `date`, `git` (for installation/sync)

### Checksum Tools (optional)

For integrity verification, one of:
- `sha256sum` (GNU coreutils)
- `shasum` (macOS)
- `openssl` (fallback)

---

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `KDSE_HOME` | `$HOME` | Installation parent directory |
| `KDSE_DIR` | `.kdse` | Installation directory name |
| `KDSE_REPO` | GitHub URL | KDSE repository URL |
| `KDSE_BRANCH` | `main` | KDSE branch |

### Configuration File

After installation, `config.sh` is generated with paths:

```bash
source ~/.kdse/config.sh   # Load configuration
```

---

## Future CLI Relationship

These scripts provide the underlying filesystem operations for a future KDSE CLI. The CLI will wrap these scripts with user interaction, progress reporting, and error handling.

Current relationship:
```
┌─────────────────┐
│   User/Operator │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     uses      ┌─────────────────┐
│   Future CLI     │──────────────▶│  Install Scripts │
└─────────────────┘               └─────────────────┘
```

The scripts are intentionally minimal and focused on deterministic file operations. The CLI will add:
- Interactive prompts
- Progress bars
- Error recovery
- Command aliases
- Shell integration

---

## Exit Codes

| Code | Name | Description |
|------|------|-------------|
| 0 | `EXIT_OK` | Success |
| 1 | `EXIT_ERROR` | General error |
| 2 | `EXIT_ALREADY_INSTALLED` | Already installed |
| 3 | `EXIT_NOT_INSTALLED` | Not installed |
| 4 | `EXIT_INVALID_ARGS` | Invalid arguments |
| 5 | `EXIT_MISSING_DEPS` | Missing dependencies |

---

## Examples

### Basic Installation

```bash
cd ~/my-project
/path/to/kdse/runtime/install/install.sh
source ~/.kdse/config.sh
```

### Regular Synchronization

```bash
# Add to CI/CD or cron
/path/to/kdse/runtime/install/sync.sh
```

### Pre-Deployment Verification

```bash
if /path/to/kdse/runtime/install/verify.sh -q; then
    echo "KDSE Runtime ready"
else
    echo "KDSE Runtime needs attention"
    exit 1
fi
```

### Clean Uninstall

```bash
/path/to/kdse/runtime/install/uninstall.sh -f
```

---

## Troubleshooting

### "Manifest not found"

Run `install.sh` first to initialize the Runtime.

### "Directory not empty"

Use `uninstall.sh -f` to force removal, or manually remove the directory.

### "Missing dependencies"

Install required utilities:
- Debian/Ubuntu: `apt-get install coreutils git`
- macOS: `brew install coreutils git` (git usually pre-installed)

### "Permission denied"

Ensure write access to `KDSE_HOME` directory.

---

## Version

- **Framework Version:** 1.1
- **KDSE Standard Version:** Referenced from manifest
- **Last Updated:** 2026-07-10

---

*This framework provides deterministic filesystem operations for KDSE Runtime management. No AI logic, no audit logic, no engineering decisions.*
