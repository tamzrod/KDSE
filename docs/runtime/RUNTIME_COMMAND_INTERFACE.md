# KDSE Runtime Command Interface

**Document Version:** 1.0  
**Type:** Engineering Operation Interface  
**Effective Date:** 2026-07-10

---

## Purpose

This document defines the KDSE Runtime Engineering Command Interface. The interface provides deterministic engineering operations that transform the Runtime from a collection of scripts into an authoritative execution layer.

The interface ensures:
- **Consistent Operations**: Same commands produce same results
- **AI Independence**: The AI merely invokes commands; the Runtime performs work
- **Operator Simplicity**: No implementation knowledge required
- **Extensibility**: New commands can be added without breaking existing ones

---

## Design Philosophy

### Current State (Incorrect)

```
Operator
    ↓
AI reasons about implementation
    ↓
AI runs multiple shell commands
    ↓
AI builds response
    ↓
Result (with AI narration)
```

### Target State (Correct)

```
Operator
    ↓
Engineering Command
    ↓
KDSE Runtime
    ↓
Deterministic Result
```

The Runtime becomes the authoritative execution layer. The AI invokes the Runtime, not individual scripts.

---

## Command Philosophy

### Principles

1. **Deterministic**: Same inputs produce same outputs
2. **Idempotent**: Running multiple times produces same result
3. **No AI Logic**: Commands only perform file operations
4. **No Engineering Decisions**: All decisions are made by the operator
5. **Machine-Readable**: Output is structured and parseable

### Contract

Each command provides:

| Element | Description |
|---------|-------------|
| **Purpose** | What the command does |
| **Inputs** | Required and optional parameters |
| **Outputs** | Expected result structure |
| **Exit Status** | Success/failure codes |
| **Messages** | Human-readable status messages |
| **Failure Modes** | Known failure conditions |

---

## Quick Reference

### Command Summary

| Command | Purpose | Category |
|---------|---------|----------|
| `kdse status` | Runtime health and status | Information |
| `kdse version` | Display runtime version | Information |
| `kdse update` | Update runtime to latest | Maintenance |
| `kdse verify` | Verify installation integrity | Maintenance |
| `kdse doctor` | Diagnose problems | Diagnostics |
| `kdse install` | Install runtime | Administration |
| `kdse uninstall` | Remove runtime | Administration |
| `kdse run` | Start new session | Session |
| `kdse resume` | Resume previous session | Session |
| `kdse audit` | Execute assessment | Session |
| `kdse history` | Show session history | Information |
| `kdse report` | List available reports | Information |

---

## Command Reference

### kdse status

**Purpose:** Display runtime health, version, and required action.

**Inputs:** None

**Outputs:**
```
============================================================
 KDSE Runtime - Runtime Status
============================================================

Installation:
  Path:        ~/.kdse
  Format:      json
  Version:     abc12345
  Installed:   2026-07-10T12:00:00Z
  Last Sync:   2026-07-10T12:00:00Z
  Repository:   https://github.com/tamzrod/KDSE.git
  Branch:      main
  Platform:    linux

Directories:
  ✓ reports/
  ✓ history/
  ✓ runtime/
  ✓ cache/
  ✓ configuration/
  ✓ standards/normative/
  ✓ standards/informative/

Required Action:
  None - runtime healthy
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Runtime not installed |

**Failure Modes:**
- Runtime not installed → Exit 2, prompt to run `kdse install`
- Missing directories → List missing items
- Legacy format → Recommend `kdse update`

---

### kdse version

**Purpose:** Display runtime version information.

**Inputs:** None

**Outputs:**
```
KDSE Runtime
Version:    abc12345
Installed:  2026-07-10T12:00:00Z
Interface:  1.0
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Runtime not installed |

---

### kdse update

**Purpose:** Update runtime to latest version from repository.

**Inputs:** None

**Outputs:**
```
============================================================
 KDSE Runtime - Runtime Update
============================================================

[INFO] Running synchronization...
[INFO] Fetching latest changes...
[INFO] Detected installation format: JSON (current)
[INFO] Pre-synchronization checks passed
[OK] Synchronization complete

============================================================
 KDSE Runtime - Synchronization Complete
============================================================

Installation Path: ~/.kdse
Previous Version:  abc12345
Current Version:   def67890
Synchronized:      2026-07-10T12:30:00Z

Installation Format: json
Preserved Data:
  reports/       - Session reports
  history/       - Session history
  runtime/       - Runtime state
  cache/         - Cached data
  configuration/ - User configuration

============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Update failed |
| 2 | Runtime not installed |

**Failure Modes:**
- Not a git repository → Exit 1, error message
- Sync conflict → Preserve user data, report conflict

---

### kdse verify

**Purpose:** Verify runtime installation integrity.

**Inputs:**
| Option | Description |
|--------|-------------|
| `-v, --verbose` | Verbose output |
| `-q, --quiet` | Quiet mode (PASS/FAIL only) |
| `-j, --json` | JSON output format |

**Outputs (standard):**
```
============================================================
 KDSE Runtime - KDSE Runtime Verification
============================================================

Detected Installation Format: json

============================================================
 KDSE Runtime - Verification Summary
============================================================

Installation Path: ~/.kdse

Results:
  PASSED:  7
  FAILED:  0
  WARNINGS: 0

[OK] All verification checks passed

============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | All checks passed |
| 1 | One or more checks failed |
| 2 | Runtime not installed |

**Checks Performed:**
1. Directory layout
2. Mandatory documents
3. Manifest integrity
4. Configuration validity
5. Version consistency
6. Preserved directories
7. Standards structure

---

### kdse doctor

**Purpose:** Diagnose runtime problems and recommend corrective action.

**Inputs:** None

**Outputs:**
```
============================================================
 KDSE Runtime - Runtime Diagnostics
============================================================

Checking Installation...
[OK] reports/
[OK] history/
[OK] runtime/
[OK] cache/
[OK] configuration/
[OK] Normative documents: 6
[OK] Manifest format: JSON (current)
[OK] Git repository: present

Disk space available: 50G

Summary:
  Issues:   0
  Warnings: 0

[OK] Runtime is healthy
============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | No issues |
| 1 | Issues found |
| 2 | Runtime not installed |

**Failure Modes:**
- Missing directory → Report as issue, increment issue count
- Legacy format → Report as warning, recommend `kdse update`
- No git repository → Report as warning

---

### kdse install

**Purpose:** Install or reinstall runtime.

**Inputs:**
| Option | Description |
|--------|-------------|
| `-f, --force` | Force reinstall |
| `-r, --repo URL` | Repository URL |
| `-b, --branch NAME` | Branch name |
| `-p, --path PATH` | Installation path |

**Outputs:**
```
============================================================
 KDSE Runtime - Installation Complete
============================================================

Installation Path: ~/.kdse
KDSE Repository:   https://github.com/tamzrod/KDSE.git
Branch:           main
Version:          def67890
Installed:        2026-07-10T12:30:00Z

Installation Format: json

Directory Structure:
  .kdse/
    reports/
    history/
    runtime/
    cache/
    configuration/
    standards/normative/
    standards/informative/
  manifest.json
  config.sh
  kdse

Quick Start:
  ~/.kdse/kdse status    # Check runtime health
  ~/.kdse/kdse update    # Update runtime
  ~/.kdse/kdse run       # Start session

  Add to PATH for convenience:
  export PATH="~/.kdse:$PATH"

============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Installation failed |
| 2 | Already installed (use --force) |

---

### kdse uninstall

**Purpose:** Remove runtime installation.

**Inputs:**
| Option | Description |
|--------|-------------|
| `-f, --force` | Force uninstall without confirmation |
| `-k, --keep-reports` | Keep reports directory |
| `-a, --keep-all` | Keep all user data |

**Outputs:**
```
============================================================
 KDSE Runtime - Uninstallation Complete
============================================================

Installation Path:  ~/.kdse
Removed:           2026-07-10T12:30:00Z

Preserved Data:
  reports/       -> ~/.kdse-backup/
  history/       -> ~/.kdse-backup/

Removed:
  KDSE Runtime
  Normative documents
  Informative documents
  Git repository
  Cache

============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Uninstallation failed |
| 2 | Runtime not installed |

---

### kdse run

**Purpose:** Start a new KDSE session.

**Inputs:** None

**Outputs:**
```
============================================================
 KDSE Runtime - Start KDSE Session
============================================================

[INFO] Session initialization...

[OK] Session started: KDSE-20260710-123000

Session ID: KDSE-20260710-123000
State:     active

The Runtime will now execute the session workflow.
Use 'kdse resume' to continue if interrupted.

============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Session already active |
| 2 | Runtime not installed |

**Failure Modes:**
- Active session exists → Exit 1, prompt to use `kdse resume`

---

### kdse resume

**Purpose:** Resume a previous or active session.

**Inputs:** None

**Outputs:**
```
============================================================
 KDSE Runtime - Resume KDSE Session
============================================================

[OK] Session found: KDSE-20260710-123000

Session ID: KDSE-20260710-123000
State:     active
Started:   2026-07-10T12:30:00Z

The Runtime will continue the session workflow.

============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | No active session |

---

### kdse audit

**Purpose:** Execute repository assessment.

**Inputs:** None

**Outputs:**
```
============================================================
 KDSE Runtime - Repository Assessment
============================================================

[INFO] Executing repository assessment...
[INFO] This would execute the Compliance Audit against the repository.
[INFO] Assessment dimensions:
  - Knowledge
  - Process
  - Artifact
  - Verification
  - Traceability
[INFO] Use 'kdse run' for a complete session including assessment.

============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Runtime not installed |

---

### kdse history

**Purpose:** Display session history.

**Inputs:** None

**Outputs:**
```
============================================================
 KDSE Runtime - Session History
============================================================

Recent Sessions:
  2026-07-10  KDSE-20260710-123000.json
  2026-07-09  KDSE-20260709-100000.json
  2026-07-08  KDSE-20260708-090000.json

Total sessions: 3

============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | History not found |

---

### kdse report

**Purpose:** List available reports.

**Inputs:** None

**Outputs:**
```
============================================================
 KDSE Runtime - Available Reports
============================================================

Reports:
  session-20260710-123000.md (2KB)
  session-20260709-100000.md (2KB)
  session-20260708-090000.md (2KB)

Total reports: 3

============================================================
```

**Exit Codes:**
| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | No reports found |

---

## Exit Status Codes

| Code | Name | Description |
|------|------|-------------|
| 0 | `EXIT_OK` | Success |
| 1 | `EXIT_ERROR` | General error |
| 2 | `EXIT_NOT_INSTALLED` | Runtime not installed |
| 3 | `EXIT_INVALID_ARGS` | Invalid arguments |

---

## Operator Messages

All commands provide consistent operator messages:

| Prefix | Meaning |
|--------|---------|
| `[INFO]` | Informational message |
| `[OK]` | Success message |
| `[WARN]` | Warning message |
| `[ERROR]` | Error message |

---

## Backward Compatibility

### Legacy Script Compatibility

The following legacy scripts are still supported:

| Legacy Script | Equivalent Command |
|--------------|-------------------|
| `install.sh` | `kdse install` |
| `sync.sh` | `kdse update` |
| `verify.sh` | `kdse verify` |
| `uninstall.sh` | `kdse uninstall` |

### Manifest Format Compatibility

- `manifest.json` - Current format (supported)
- `manifest.yaml` - Legacy format (auto-migrated)

---

## Shell Integration

### Adding to PATH

Add to your shell profile (~/.bashrc, ~/.zshrc):

```bash
# KDSE Runtime
export PATH="~/.kdse:$PATH"
```

### Using Configuration

```bash
# Source the configuration
source ~/.kdse/config.sh

# Now kdse is in PATH
kdse status
```

### Shell Completion

Source the completion script for tab completion:

```bash
# Bash
source ~/.kdse/kdse-completion.sh

# Zsh
autoload bashcompinit
bashcompinit
source ~/.kdse/kdse-completion.sh
```

---

## AI Integration

### Correct Usage

The AI should invoke the Runtime, not reason about implementation:

```
# CORRECT
User: "Update KDSE"
AI: Invokes "kdse update"
Runtime: Performs update
Result: Displayed to user

# INCORRECT
User: "Update KDSE"
AI: Reasons about git pull, checks files, runs multiple commands
AI: Builds response explaining what it did
Result: AI narration mixed with result
```

### Command Mapping

| Operator Request | Command |
|-----------------|---------|
| "Update KDSE" | `kdse update` |
| "Check status" | `kdse status` |
| "Verify installation" | `kdse verify` |
| "Diagnose problems" | `kdse doctor` |
| "Start session" | `kdse run` |
| "Resume session" | `kdse resume` |
| "Show history" | `kdse history` |
| "List reports" | `kdse report` |
| "Run audit" | `kdse audit` |

---

## Versioning

The command interface follows semantic versioning:

- **Major Version**: Breaking changes to commands
- **Minor Version**: New commands added
- **Patch Version**: Bug fixes, documentation updates

Current interface version: 1.0

---

## Related Documents

| Document | Description |
|----------|-------------|
| [COMMANDS.md](COMMANDS.md) | Session command interface |
| [SESSION_PROTOCOL.md](SESSION_PROTOCOL.md) | Session lifecycle |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Runtime architecture |
| [INSTALL_README.md](../runtime/install/README.md) | Installation guide |

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-07-10 | Initial command interface definition |

---

*This document defines the KDSE Runtime Engineering Command Interface. The Runtime provides deterministic operations; the AI merely invokes commands.*
