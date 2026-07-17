# KDSE Architecture Migration Guide

**Document Version:** 1.0  
**Effective Date:** 2026-07-17

---

## Overview

This guide helps migrate from the current KDSE architecture to the Runtime-Centric Architecture (RCA).

---

## Breaking Changes

### Summary

| Category | Change | Impact |
|----------|--------|--------|
| CLI | Thin adapter only | Commands delegate to engine |
| State | Centralized in Workspace Engine | No direct state access |
| Initialization | Template-based | Different init flow |
| Verification | Required before actions | New behavior |

### Detailed Changes

#### 1. CLI Runtime Changes

**Before:**
```bash
kdse init              # Mixed business logic
kdse collect           # Direct artifact processing
kdse orchestrate       # State manipulation
```

**After:**
```bash
kdse init              # Thin adapter, calls engine
kdse verify            # Calls engine verification
kdse phase show        # Calls engine phase query
kdse phase advance     # Calls engine transition
```

**What Changed:**
- CLI no longer contains business logic
- All operations go through Workspace Engine
- Consistent interface for CLI and MCP

#### 2. State Management Changes

**Before:**
```go
// Direct state access
state.Load()
state.Save()
runtime.Manage()
```

**After:**
```go
// Through Workspace Engine only
engine := workspace.NewEngine(path)
ws, _ := engine.LoadWorkspace(ctx)
engine.AdvancePhase(ctx, target)
```

**What Changed:**
- State ownership unified in Workspace Engine
- No direct state access
- Clear boundaries

#### 3. Initialization Changes

**Before:**
```bash
kdse init
# Creates files with hardcoded values
# No verification
```

**After:**
```bash
kdse init --type cli --template default
# Template-based initialization
# Atomic operations
# Verification gate
```

**What Changed:**
- Templates define structure
- Atomic initialization with rollback
- Verification before completion

---

## Migration Steps

### Step 1: Backup Your Project

```bash
# Backup entire project
cp -r /path/to/project /path/to/project.backup

# Backup .kdse if exists
if [ -d ".kdse" ]; then
    cp -r .kdse .kdse.backup
fi
```

### Step 2: Update KDSE Binary

```bash
# Download new version
curl -LO https://github.com/kdse/runtime/releases/latest/kdse

# Verify checksum
sha256sum -c kdse.sha256

# Install
chmod +x kdse
sudo mv kdse /usr/local/bin/
```

### Step 3: Verify Current State

```bash
# Check if .kdse exists
ls -la .kdse

# Check current phase (if .kdse exists)
cat .kdse/phase.yaml 2>/dev/null || echo "No phase.yaml"
```

### Step 4: Run Migration (if needed)

For existing projects with .kdse:

```bash
# Migrate existing workspace
kdse migrate --from=v1 --to=v2

# Verify migration
kdse verify
```

### Step 5: Update Scripts and CI/CD

If you have scripts using KDSE:

**Before:**
```bash
#!/bin/bash
kdse collect
kdse orchestrate --phase Foundation
```

**After:**
```bash
#!/bin/bash
kdse verify
kdse phase advance
```

---

## New Directory Structure

### Before

```
project/
├── .kdse/
│   ├── runtime/
│   ├── foundation/
│   ├── knowledge/
│   ├── evidence/
│   └── reports/
```

### After

```
project/
├── .kdse/
│   ├── runtime.yaml          # Runtime configuration
│   ├── workspace.yaml         # Workspace state
│   ├── methodology.yaml       # Methodology reference
│   ├── phase.yaml              # Current phase
│   ├── session.yaml            # Session state
│   ├── metadata.yaml            # Runtime metadata
│   ├── phase-history.yaml       # Phase transitions
│   ├── knowledge/               # Knowledge artifacts
│   │   ├── requirements.md
│   │   ├── stakeholders.md
│   │   └── constraints.md
│   ├── architecture/           # Architecture artifacts
│   │   ├── architecture.md
│   │   └── decisions.md
│   ├── implementation/          # Implementation artifacts
│   │   └── implementation.md
│   ├── verification/            # Verification artifacts
│   │   └── verification.md
│   └── reports/                 # Generated reports
```

---

## Command Changes

### Removed Commands

| Old Command | Reason |
|-------------|--------|
| `kdse collect` | Knowledge collection via engine |
| `kdse orchestrate` | Orchestration via engine |
| `kdse normalize` | Normalization via engine |
| `kdse runtime verify` | Unified in `kdse verify` |

### New Commands

| New Command | Description |
|-------------|-------------|
| `kdse verify` | Verify workspace state |
| `kdse phase show` | Show current phase |
| `kdse phase advance` | Advance to next phase |
| `kdse artifacts <phase>` | List phase artifacts |
| `kdse report --type <type>` | Generate reports |

### Changed Commands

| Old | New | Change |
|-----|-----|--------|
| `kdse init` | `kdse init` | Template-based, atomic |
| `kdse status` | `kdse phase show` | Renamed for clarity |
| `kdse report` | `kdse report --type summary` | Explicit type |

---

## API Changes

### Before (Direct Access)

```go
// Old API
runtime := kdseruntime.New(repoPath)
runtime.Initialize()
state := state.Load()
```

### After (Engine-Based)

```go
// New API
engine := workspace.NewEngine(repoPath)
ws, _ := engine.InitializeWorkspace(ctx, opts)
phase, _ := engine.GetPhase(ctx)
engine.AdvancePhase(ctx, target)
```

---

## MCP Changes

### Before

```json
// Old MCP tools
{
  "tools": [
    {"name": "kdse_collect", "handler": collectHandler},
    {"name": "kdse_orchestrate", "handler": orchestrateHandler}
  ]
}
```

### After

```json
// New MCP tools
{
  "tools": [
    {"name": "kdse_init", "handler": initHandler},
    {"name": "kdse_verify", "handler": verifyHandler},
    {"name": "kdse_phase", "handler": phaseHandler},
    {"name": "kdse_artifacts", "handler": artifactsHandler},
    {"name": "kdse_report", "handler": reportHandler}
  ]
}
```

---

## Rollback Procedure

If migration fails:

```bash
# Restore backup
rm -rf .kdse
cp -r .kdse.backup .kdse

# Verify restored
kdse status
```

---

## Testing After Migration

```bash
# 1. Verify workspace
kdse verify

# 2. Check phase
kdse phase show

# 3. List artifacts
kdse artifacts knowledge

# 4. Generate report
kdse report --type summary

# 5. Check for warnings
kdse verify --verbose
```

---

## Common Issues

### Issue: "Runtime not found"

**Cause:** .kdse directory missing

**Solution:**
```bash
kdse init
```

### Issue: "Invalid phase transition"

**Cause:** Attempting to skip phases

**Solution:**
```bash
# Check current phase
kdse phase show

# Advance to next phase only
kdse phase advance
```

### Issue: "Verification failed"

**Cause:** Missing required artifacts

**Solution:**
```bash
# Check what's missing
kdse verify

# Create missing artifacts
# Then verify again
kdse verify
```

---

## Benefits of Migration

1. **Single Source of Truth**: .kdse is authoritative
2. **Consistent Behavior**: CLI and MCP behave identically
3. **Evidence-Based**: Verification before actions
4. **Atomic Operations**: No partial states
5. **Clear Boundaries**: Methodology independence
6. **Future-Proof**: Easy to add new runtimes

---

## Related Documents

| Document | Description |
|----------|-------------|
| [KAE-001](KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md) | Architecture evolution details |
| [RUNTIME_ARCHITECTURE.md](../architecture/RUNTIME_ARCHITECTURE.md) | Target architecture |
| [PRINCIPLES.md](../architecture/PRINCIPLES.md) | Core principles |

---

*This document is informative. For normative guidance, see the architecture documents.*
