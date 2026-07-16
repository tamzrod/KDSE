# Knowledge-Driven Software Engineering (KDSE)

**Lean KDSE** is a practical engineering methodology where structured knowledge serves as the authoritative source from which all software artifacts are derived, maintained, and verified.

## Quick Start

```bash
# Initialize KDSE workspace
kdse init

# Add knowledge to notebook
kdse notebook add "Users need password reset" --source "customer-feedback.md"

# Promote to candidate
kdse promote submit <entry-id>

# View status
kdse status
```

## Core Principles

| # | Principle | What It Means |
|---|-----------|---------------|
| 1 | Knowledge precedes architecture | Derive, don't assume |
| 2 | Architecture precedes implementation | Follow the design |
| 3 | Authority flows downward | Lower can't contradict higher |
| 4 | Traceability enables authority | Every decision traces to knowledge |
| 5 | Repository first | Analyze artifacts before asking |

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

| Command | Description |
|---------|-------------|
| `kdse init` | Initialize `.kdse/` workspace |
| `kdse status` | Show current state |
| `kdse notebook add` | Add insight to notebook |
| `kdse promote submit` | Submit candidate for review |
| `kdse promote review` | Accept/reject with rationale |
| `kdse agreement init` | Initialize project agreement |

## License

Apache 2.0