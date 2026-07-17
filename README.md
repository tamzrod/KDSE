# Knowledge-Driven Software Engineering (KDSE)

**Lean KDSE** is a practical engineering methodology where structured knowledge serves as the authoritative source from which all software artifacts are derived, maintained, and verified.

## Session Guard (Enforcement Layer)

KDSE implements a **Session Guard** that strictly enforces initialization. This ensures:

1. **Workspace Validation**: All operations require a valid `.kdse/` directory
2. **Session State Verification**: Session state is verified before any operation proceeds
3. **Auto-Recovery**: Missing workspace is automatically initialized on first run
4. **Thread-Safe**: Guard operations are synchronized for concurrent access

### Initialization Workflow

KDSE v2.0 uses **GitHub Bootstrap Initialization** - workspace templates are downloaded from the official KDSE repository, making GitHub the authoritative source for methodology artifacts.

```bash
# First run - downloads official workspace template from GitHub
kdse init

# Subsequent runs - session guard validates automatically
kdse execute "Build a login system"
```

### Guard Enforcement Points

The following operations are blocked if workspace is not initialized:

- `execute` - Primary orchestration tool
- `collect` - Artifact collection
- `foundation` - Foundation documentation
- `audit` - Audit report generation
- `migrate` - Legacy directory migration

Operations that don't require initialization:

- `initialize` - Sets up the workspace
- `status` - Shows current state (with guard status)
- `help` - Shows available tools

## Quick Start

```bash
# Initialize KDSE workspace (downloads template from GitHub)
kdse init

# Initialize with specific template (future)
kdse init --template=web

# Add knowledge to notebook
kdse notebook add "Users need password reset" --source "customer-feedback.md"

# Promote to candidate
kdse promote submit <entry-id>

# View status
kdse status
```

## Workspace Structure

The `.kdse/` workspace contains:

```
.kdse/
├── session.yaml         # Session metadata
├── authority.yaml       # Authority rules
├── phase.yaml          # Phase configuration
├── project.yaml        # Project metadata
├── knowledge/          # Knowledge artifacts
├── architecture/       # System design
├── implementation/     # Development tracking
├── verification/       # Testing artifacts
├── reports/            # Generated reports
└── docs/              # Documentation
```

## Core Principles

| # | Principle | What It Means |
|---|-----------|---------------|
| 1 | Knowledge precedes architecture | Derive, don't assume |
| 2 | Architecture precedes implementation | Follow the design |
| 3 | Authority flows downward | Lower can't contradict higher |
| 4 | Traceability enables authority | Every decision traces to knowledge |
| 5 | Repository first | Analyze artifacts before asking |

## Bootstrap Architecture

The bootstrap layer downloads workspace templates from GitHub:

1. **Download**: Fetches official template from GitHub archive
2. **Extract**: Unpacks template into workspace
3. **Verify**: Ensures all required artifacts exist
4. **Rollback**: Cleans up on failure

Templates are stored in the [kdse/workspace-templates](https://github.com/kdse/workspace-templates) repository.

## Foundation Documents

| Document | Purpose |
|----------|---------|
| [001-principles.md](docs/foundation/001-principles.md) | Core principles and philosophy |
| [002-knowledge-engine.md](docs/foundation/002-knowledge-engine.md) | Evidence → Knowledge pipeline |
| [003-authority-traceability.md](docs/foundation/003-authority-traceability.md) | Authority hierarchy and traceability |
| [004-workspace.md](docs/foundation/004-workspace.md) | `.kdse/` workspace structure |
| [005-adoption.md](docs/foundation/005-adoption.md) | Getting started guide |

## Architecture

```
Evidence → Derivation → Knowledge Artifact → Architecture → Implementation → Verification
     ↑                                                      ↓
     └──────────── Evidence Strength (confidence) ───────────┘
```

## Key Concepts

- **Evidence**: Raw domain information from any source
- **Knowledge Artifact**: Structured understanding with Evidence Strength
- **Derivation Pipeline**: Evidence → Derivation → Knowledge
- **Agreement**: Compact project state for delta communication

## CLI Commands

| Command | Description | Requires Initialization |
|---------|-------------|------------------------|
| `kdse init` | Initialize `.kdse/` workspace | No |
| `kdse initialize` | Full runtime initialization | No |
| `kdse status` | Show current state | No |
| `kdse notebook add` | Add insight to notebook | Yes |
| `kdse promote submit` | Submit candidate for review | Yes |
| `kdse promote review` | Accept/reject with rationale | Yes |
| `kdse agreement init` | Initialize project agreement | Yes |

## MCP Server

The KDSE MCP Server provides Model Context Protocol access to KDSE capabilities:

```bash
# Start MCP server (stdio mode)
kdse-mcp

# Start MCP server (HTTP mode)
kdse-mcp --transport http
```

### MCP Tools

| Tool | Description | Guard Enforced |
|------|-------------|----------------|
| `execute` | Primary orchestration tool | Yes |
| `initialize` | Initialize workspace | No |
| `status` | Show status | No |
| `collect` | Collect artifacts | Yes |
| `foundation` | Foundation docs | Yes |
| `audit` | Generate audit report | Yes |
| `migrate` | Migrate legacy dirs | Yes |

## License

Apache 2.0