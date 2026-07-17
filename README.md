# Knowledge-Driven Software Engineering (KDSE)

**KDSE** is a Runtime-Centric Engineering Platform where the KDSE Runtime (`.kdse/`) is the authoritative foundation of every engineering project.

## Core Principles

| # | Principle | What It Means |
|---|-----------|---------------|
| P-001 | Evidence First | No engineering work without evidence |
| P-002 | Runtime is the Authority | `.kdse/` is the authoritative state |
| P-003 | One Methodology | Only one engineering methodology |
| P-004 | Multiple Runtimes | CLI, MCP, IDE are execution adapters |
| P-005 | Runtime Independence | Methodology never depends on runtime |
| P-006 | Filesystem is Evidence | Without `.kdse/`, project is NOT KDSE |

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                       KDSE Methodology                                │
│                        (Engineering Rules)                            │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Workspace Engine                                 │
│                       (State Owner)                                  │
│                                                                     │
│  • Create/Load/Verify Runtime                                       │
│  • Phase Management                                                  │
│  • Artifact Validation                                              │
│  • Report Generation                                                │
└─────────────────────────────────────────────────────────────────────┘
                    ┌───────────────────┬───────────────────┐
                    │                   │                   │
                    ▼                   ▼                   ▼
            ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
            │ CLI Runtime │     │ MCP Runtime │     │ IDE Runtime │
            │(Thin Adapter│     │(Thin Adapter│     │  (Future)  │
            └─────────────┘     └─────────────┘     └─────────────┘
```

## Quick Start

```bash
# Initialize KDSE workspace
kdse init

# Verify workspace state
kdse verify

# Show current phase
kdse phase show

# Advance to next phase
kdse phase advance

# Generate report
kdse report --type summary
```

## KDSE Runtime (.kdse/)

The `.kdse/` directory is **mandatory**. Without it, the project is NOT a KDSE project.

```
.kdse/
├── runtime.yaml           # Runtime configuration (type, version)
├── workspace.yaml         # Workspace state
├── methodology.yaml       # Methodology reference
├── phase.yaml              # Current phase
├── session.yaml            # Session state
├── metadata.yaml            # Runtime metadata
├── knowledge/               # Knowledge artifacts
│   ├── requirements.md
│   ├── stakeholders.md
│   └── constraints.md
├── architecture/           # Architecture artifacts
│   ├── architecture.md
│   └── decisions.md
├── implementation/          # Implementation artifacts
├── verification/            # Verification artifacts
└── reports/                 # Generated reports
```

## Engineering Phases

```
┌──────────────┐
│ Initialization │ ← Start here (.kdse created)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Knowledge   │ ← Document requirements
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Architecture │ ← Design system
└──────┬───────┘
       │
       ▼
┌──────────────┐
│Implementation│ ← Build system
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Verification │ ← Test and validate
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Reports    │ ← Document results
└──────────────┘
```

**Skipping phases is prohibited.**

## Evidence First

No engineering action proceeds without verified evidence:

```
Before ANY operation:
1. Verify .kdse/ exists
2. Load runtime configuration
3. Load current phase
4. Validate required artifacts
5. Continue (only if all pass)

If verification fails → STOP, return error
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `kdse init` | Initialize KDSE runtime |
| `kdse verify` | Verify workspace state |
| `kdse phase show` | Show current phase |
| `kdse phase advance` | Advance to next phase |
| `kdse artifacts <phase>` | List phase artifacts |
| `kdse report --type <type>` | Generate report |

## MCP Tools

| Tool | Description |
|------|-------------|
| `kdse_init` | Initialize KDSE runtime |
| `kdse_verify` | Verify workspace state |
| `kdse_phase` | Show/advance phase |
| `kdse_artifacts` | List artifacts |
| `kdse_report` | Generate report |

## Architecture Documents

| Document | Description |
|----------|-------------|
| [docs/architecture/PRINCIPLES.md](docs/architecture/PRINCIPLES.md) | Core principles |
| [docs/architecture/RUNTIME_ARCHITECTURE.md](docs/architecture/RUNTIME_ARCHITECTURE.md) | Runtime architecture |
| [docs/architecture/METHODOLOGY.md](docs/architecture/METHODOLOGY.md) | Engineering methodology |
| [docs/architecture/WORKSPACE_ENGINE.md](docs/architecture/WORKSPACE_ENGINE.md) | Workspace Engine |
| [docs/architecture/CLI_RUNTIME.md](docs/architecture/CLI_RUNTIME.md) | CLI Runtime |
| [docs/architecture/MCP_RUNTIME.md](docs/architecture/MCP_RUNTIME.md) | MCP Runtime |

## Architecture Evolution

| Document | Description |
|----------|-------------|
| [docs/architecture-evolution/README.md](docs/architecture-evolution/README.md) | Evolution index |
| [docs/architecture-evolution/KAE-001](docs/architecture-evolution/KAE-001-RUNTIME-CENTRIC-ARCHITECTURE.md) | Runtime-Centric Architecture |
| [docs/architecture-evolution/KAE-002](docs/architecture-evolution/KAE-002-CLI-MCP-RUNTIME-SEPARATION.md) | CLI/MCP Separation |
| [docs/architecture-evolution/KAE-003](docs/architecture-evolution/KAE-003-RUNTIME-BOOTSTRAP.md) | Runtime Bootstrap |
| [docs/architecture-evolution/MIGRATION.md](docs/architecture-evolution/MIGRATION.md) | Migration Guide |

## ADRs (Architecture Decision Records)

| ADR | Decision |
|-----|----------|
| [ADR-001](docs/decisions/ADR-001-RUNTIME-IS-THE-AUTHORITY.md) | Runtime is the Authority |
| [ADR-002](docs/decisions/ADR-002-ONE-METHODOLOGY-MULTIPLE-RUNTIMES.md) | One Methodology, Multiple Runtimes |
| [ADR-003](docs/decisions/ADR-003-WORKSPACE-ENGINE.md) | Workspace Engine |

## License

Apache 2.0