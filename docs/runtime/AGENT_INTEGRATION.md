# KDSE Runtime Agent Integration Guide

**Document Version:** 1.0  
**Type:** AI Agent Integration Specification  
**Effective Date:** 2026-07-10

---

## Purpose

This document defines how AI agents should discover, resolve, and invoke KDSE Runtime engineering commands. The goal is to enable AI agents to dispatch directly to the appropriate command without reasoning about implementation details.

---

## Overview

The KDSE Runtime exposes a deterministic Engineering Command Interface. AI agents should use this interface rather than:

- Inspecting runtime scripts
- Reasoning about file locations
- Building execution plans
- Manually running shell commands

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         AI Agent                                      │
├─────────────────────────────────────────────────────────────────────┤
│  1. Receive natural language request                                  │
│  2. Resolve to command using registry                                 │
│  3. Invoke command                                                   │
│  4. Return result                                                    │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    KDSE Runtime                                        │
├─────────────────────────────────────────────────────────────────────┤
│  .kdse/                                                              │
│  ├── kdse              # Command interface                            │
│  ├── runtime/commands.yaml  # Command registry                        │
│  └── ...                                                           │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Command Discovery

### Step 1: Locate Command Registry

The command registry is located at:

```
~/.kdse/runtime/commands.yaml
```

Or via the Runtime installation path:

```bash
${KDSE_HOME}/.kdse/runtime/commands.yaml
```

### Step 2: Load Registry

AI agents should load the command registry to discover available commands:

```bash
# Via command
kdse commands

# Or directly
cat ~/.kdse/runtime/commands.yaml
```

### Step 3: Parse Commands

Each command in the registry defines:

| Field | Description |
|-------|-------------|
| `name` | Command identifier |
| `aliases` | Alternative invocations |
| `purpose` | What the command does |
| `category` | Command category (information, maintenance, session, administration) |
| `inputs` | Required and optional inputs |
| `outputs` | Expected output fields |
| `exit_codes` | Success and failure codes |
| `examples` | Usage examples |
| `natural_language_patterns` | Patterns that resolve to this command |

---

## Natural Language Resolution

### Resolution Process

When the AI receives a request, it should:

1. **Parse the request** for intent
2. **Match against patterns** in `natural_language_patterns`
3. **Select the command** with the best match
4. **Invoke the command** directly

### Resolution Map

The registry includes a `resolution_map` that maps common requests to commands:

| Request | Command |
|---------|---------|
| "Update KDSE" | `kdse update` |
| "KDSE status" | `kdse status` |
| "Run KDSE" | `kdse run` |
| "Resume KDSE" | `kdse resume` |
| "Verify KDSE" | `kdse verify` |
| "KDSE doctor" | `kdse doctor` |
| "Show reports" | `kdse report` |
| "Show history" | `kdse history` |

### Pattern Matching Examples

| User Request | Matched Pattern | Command |
|--------------|-----------------|---------|
| "Update KDSE" | "update KDSE" | `kdse update` |
| "check if KDSE is working" | "is KDSE working" | `kdse status` |
| "start a new session" | "start" | `kdse run` |
| "continue where I left off" | "continue" | `kdse resume` |
| "verify the installation" | "verify" | `kdse verify` |
| "something is wrong" | "problems" | `kdse doctor` |

---

## When to Invoke the Runtime

### Always Invoke For:

| Category | Examples |
|----------|----------|
| **Runtime Operations** | status, version, update, verify, doctor |
| **Session Management** | run, resume, audit |
| **Information** | history, report |
| **Administration** | install, uninstall |

### Never Reason About:

| Instead of this | Do this |
|-----------------|---------|
| Inspecting shell scripts | Invoke `kdse status` |
| Checking for git repository | Invoke `kdse verify` |
| Building update commands | Invoke `kdse update` |
| Manually running sync | Invoke `kdse update` |
| Inspecting manifest files | Invoke `kdse status` |
| Creating execution plans | Invoke appropriate command |

---

## Command Invocation Pattern

### Correct Pattern

```python
# 1. Receive request
request = "Update KDSE"

# 2. Resolve to command (from registry)
command = "kdse update"

# 3. Invoke directly
result = run_command(command)

# 4. Return result to user
return result
```

### Incorrect Pattern

```python
# 1. Receive request
request = "Update KDSE"

# 2. INCORRECT: Reason about implementation
if not os.path.exists("~/.kdse"):
    # Try to clone repository
    run("git clone ...")
    
# Check manifest
manifest = read_file("~/.kdse/manifest.json")

# Determine if update needed
if needs_update(manifest):
    # Run git pull
    run("cd ~/.kdse && git pull")
    
# Verify
run("cd ~/.kdse && ./verify.sh")

return "Updated successfully"
```

---

## Output Handling

### Standard Output

Commands produce deterministic, machine-readable output:

```
============================================================
 KDSE Runtime - Runtime Status
============================================================

Installation:
  Path:        ~/.kdse
  Format:      json
  Version:     abc12345
  ...

============================================================
```

### Exit Codes

| Code | Name | Action |
|------|------|--------|
| 0 | `EXIT_OK` | Success - operation completed |
| 1 | `EXIT_ERROR` | Failed - check error message |
| 2 | `EXIT_NOT_INSTALLED` | Runtime not installed - invoke `kdse install` |
| 3 | `EXIT_INVALID_ARGS` | Invalid arguments - check usage |

### Error Handling

```python
result = run_command("kdse status")
if result.exit_code == 2:
    # Runtime not installed
    run_command("kdse install")
elif result.exit_code != 0:
    # Error occurred
    show_error(result.stderr)
else:
    show_result(result.stdout)
```

---

## Command Reference

### Information Commands

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `kdse status` | Show runtime health | Any status check |
| `kdse version` | Show version | Version inquiry |
| `kdse commands` | List commands | Discovery |
| `kdse history` | Show session history | Past session inquiry |
| `kdse report` | List reports | Report inquiry |

### Maintenance Commands

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `kdse update` | Update runtime | Refresh installation |
| `kdse verify` | Verify integrity | Installation check |
| `kdse doctor` | Diagnose problems | Troubleshooting |

### Session Commands

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `kdse run` | Start session | New audit session |
| `kdse resume` | Resume session | Continue previous |
| `kdse audit` | Run assessment | Repository assessment |

### Administration Commands

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `kdse install` | Install runtime | Fresh installation |
| `kdse uninstall` | Remove runtime | Clean removal |

---

## Self-Discovery Protocol

The AI agent should:

1. **On startup**: Check if KDSE is installed
   ```bash
   kdse status
   ```

2. **If not installed**: Prompt user or install
   ```bash
   kdse install
   ```

3. **On command unknown**: List available commands
   ```bash
   kdse commands
   ```

4. **On error**: Use doctor for diagnosis
   ```bash
   kdse doctor
   ```

---

## Example Flows

### Flow 1: Update Request

```
User: "Update KDSE"

AI:
  1. Match "Update KDSE" to command: kdse update
  2. Invoke: kdse update
  3. Return result
```

### Flow 2: Status Check

```
User: "Is KDSE working?"

AI:
  1. Match "working" to status command
  2. Invoke: kdse status
  3. Return health status
```

### Flow 3: Session Management

```
User: "Run KDSE"

AI:
  1. Match "Run KDSE" to run command
  2. Invoke: kdse run
  3. Return session ID
```

### Flow 4: Problem Diagnosis

```
User: "KDSE is having issues"

AI:
  1. Match "issues" to doctor command
  2. Invoke: kdse doctor
  3. Present diagnosis and recommendations
```

---

## Registry Format

The command registry is YAML for human readability:

```yaml
commands:
  - name: status
    aliases:
      - "kdse status"
      - "show status"
    purpose: |
      Display runtime health, version, and required action.
    category: information
    exit_codes:
      - code: 0
        meaning: Success
      - code: 2
        meaning: Not installed
    natural_language_patterns:
      - "status"
      - "health"
```

---

## Versioning

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2026-07-10 | Initial agent integration specification |

---

## Related Documents

| Document | Description |
|----------|-------------|
| [RUNTIME_COMMAND_INTERFACE.md](RUNTIME_COMMAND_INTERFACE.md) | Command interface reference |
| [COMMANDS.md](../../COMMANDS.md) | Session command interface |
| [SESSION_PROTOCOL.md](../../SESSION_PROTOCOL.md) | Session lifecycle |

---

*This document defines how AI agents should integrate with KDSE Runtime. The agent invokes commands; the Runtime performs work.*
