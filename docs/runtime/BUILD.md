# KDSE Runtime Binary

**Version:** 1.0.0  
**Status:** Production Ready

---

## Overview

The KDSE Runtime Binary transforms KDSE from a documentation-driven methodology into a binary-driven Engineering Runtime. The binary provides a standardized operational workflow for applying KDSE principles during day-to-day engineering work.

## Quick Start

### Build

```bash
go build -o kdse ./cmd/kdse/
```

### Commands

```bash
# Install runtime configuration
kdse install

# Start a KDSE session (detects repository automatically)
kdse run

# View session status
kdse status

# Generate runtime report
kdse report

# Check for updates
kdse update
```

## Runtime Responsibilities

The runtime answers four engineering questions:

| Question | Responsibility |
|----------|----------------|
| Where am I? | Detect repository, determine engineering phase |
| What do I know? | Read documentation, build engineering context |
| What should happen next? | Determine current objective, produce phase guidance |
| What happened? | Persist state, produce reports |

## Generated Artifacts

The runtime creates the following artifacts in `.kdse/`:

| File | Description |
|------|-------------|
| `state.json` | Current session state |
| `context.json` | Engineering context |
| `reports/` | Generated runtime reports |

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      KDSE Standard                           │
│                   (Normative - Markdown/YAML)               │
└───────────────────────────┬─────────────────────────────────┘
                            │ Consumes
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    KDSE Runtime Binary                      │
│                   (Informative - Go CLI)                   │
│                                                             │
│  • Repository Detection                                     │
│  • Phase Identification                                     │
│  • Context Building                                        │
│  • State Management                                        │
│  • Report Generation                                       │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ Orchestrates
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  Engineering Participants                  │
│                                                             │
│  • Human Engineer                                          │
│  • OpenHands                                               │
│  • Claude Code                                             │
│  • Cursor                                                  │
│  • Codex                                                   │
│  • CI/CD                                                   │
└─────────────────────────────────────────────────────────────┘
```

## Phase Detection

The runtime detects the current engineering phase based on repository artifacts:

| Phase | Score Range | Indicators |
|-------|-------------|------------|
| Concept | 0-2 | Empty repository |
| Defined | 2-4 | Git initialized, basic files |
| Structured | 4-6 | Documentation, source code |
| Usable | 6-8 | Tests, requirements |
| Validated | 8-9 | Architecture, governance |
| Proven | 9-10 | Full KDSE compliance |

## Session State

Runtime states follow the KDSE execution model:

```
Idle → Loading → Verifying → Assessing → Reporting → Pending Approval
                                                          ↓
                                    ┌─────────────────────┼─────────────────────┐
                                    │                     │                     │
                                    ▼                     ▼                     ▼
                              Implement             Reject               Defer
                                    │                                         │
                                    ▼                                         │
                             Verify Results                                   │
                                    │                                         │
                                    ├─────────────────────┐                   │
                                    │                     │                   │
                                    ▼                     ▼                   ▼
                               Complete              Report               Paused
```

## Executor Model

The runtime is executor-agnostic. Supported executors include:

- **Human** - Manual operation
- **OpenHands** - AI agent orchestration
- **Claude Code** - Anthropic Claude integration
- **Cursor** - Cursor IDE integration
- **Codex** - OpenAI Codex integration
- **CI/CD** - Automated pipelines

## Build Requirements

- Go 1.21+
- Cross-platform compatible (Linux, macOS, Windows)

## Installation

### Binary Installation

```bash
# Download the appropriate binary for your platform
# Extract and place in PATH

# Verify installation
kdse install
kdse --version
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/kdse/runtime.git
cd runtime

# Build
go build -o kdse ./cmd/kdse/

# Install configuration
./kdse install

# Test
./kdse --version
```

## Verification

```bash
# Test 1: Verify startup
kdse --version

# Test 2: Start a session
kdse run

# Test 3: View status
kdse status

# Test 4: Generate report
kdse report

# Test 5: Verify artifacts created
ls -la .kdse/
cat .kdse/state.json
cat .kdse/context.json
ls -la .kdse/reports/
```

## Implementation Philosophy

The runtime follows these principles:

1. **Small footprint** - Minimal code, no framework dependencies
2. **Simple functions** - Prefer functions over services
3. **No premature optimization** - Simple is better than clever
4. **No unnecessary interfaces** - Only add interfaces when necessary
5. **Documentation is truth** - Markdown/YAML are authoritative

## Source of Truth

> Markdown and YAML remain authoritative.
> The runtime never replaces documentation.
> It only interprets documentation.
> Generated files are cache only.

## Related Documents

- [KDSE Runtime Architecture](../ARCHITECTURE.md)
- [KDSE Execution Model](../EXECUTION_MODEL.md)
- [KDSE Session Protocol](../SESSION_PROTOCOL.md)
- [KDSE Report Specification](../REPORT_SPEC.md)
- [KDSE Commands](../COMMANDS.md)

---

*This runtime is an informative reference implementation. It demonstrates how KDSE can be applied in practice while preserving the normative authority of the KDSE Standard.*
