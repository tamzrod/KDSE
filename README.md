# Knowledge-Driven Software Engineering (KDSE)

**KDSE** is an Engineering Runtime that **augments** software engineering practices. It never replaces them.

## Core Principles

| # | Principle | What It Means |
|---|-----------|---------------|
| P-001 | Evidence First | No engineering work without evidence |
| P-002 | Project is the Foundation | The software project owns all deliverables |
| P-003 | Runtime Augments | KDSE runtime supports, not replaces engineering |
| P-004 | Ownership Boundaries | Project vs Runtime layers are strictly separated |
| P-005 | Standard Practices | KDSE enhances, not bypasses, standard engineering |
| P-006 | Runtime is the Engineering Runtime | `.kdse/` is for engineering artifacts, not project deliverables |

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

The `.kdse/` directory contains the **KDSE Engineering Runtime** - it is NOT the project workspace.

**Key Rule:** A KDSE-enabled repository must remain a standard software project. Any engineer unfamiliar with KDSE should still be able to understand, build, and maintain the project.

```
Project/                     # ← Project Layer (owned by software project)
├── README.md               # Project documentation
├── docs/                   # Project documentation
├── src/                    # Source code
├── tests/                  # Test code
└── .kdse/                 # ← Runtime Layer (owned by KDSE)
    ├── runtime/            # Runtime state
    ├── sessions/          # Engineering sessions
    ├── reports/           # Engineering reports
    ├── evidence/          # Engineering evidence
    ├── traceability/      # Traceability links
    ├── references/        # External standards (IEC 61850, Modbus, etc.)
    ├── knowledge/         # Extracted knowledge
    └── laboratory/        # Engineering laboratory
```

### Ownership Boundaries

| Layer | Owner | Contains | Location |
|-------|-------|----------|----------|
| **Project Layer** | Software Project | Deliverables, docs, code | Project root |
| **Runtime Layer** | KDSE Runtime | State, sessions, evidence | `.kdse/` |
| **Reference Layer** | KDSE Runtime | External standards | `.kdse/references/` |
| **Knowledge Layer** | KDSE Runtime | Extracted knowledge | `.kdse/knowledge/` |

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
| [docs/architecture/OWNERSHIP_MODEL.md](docs/architecture/OWNERSHIP_MODEL.md) | **Ownership model** - Project vs Runtime boundaries |
| [docs/architecture/AUDIT-ARCHITECTURE-DRIFT.md](docs/architecture/AUDIT-ARCHITECTURE-DRIFT.md) | Architecture audit report |
| [docs/architecture/PRINCIPLES.md](docs/architecture/PRINCIPLES.md) | Core principles |
| [docs/architecture/RUNTIME_ARCHITECTURE.md](docs/architecture/RUNTIME_ARCHITECTURE.md) | Runtime architecture |
| [docs/architecture/METHODOLOGY.md](docs/architecture/METHODOLOGY.md) | Engineering methodology |

## KDSE Runtime Documentation

| Document | Description |
|----------|-------------|
| [docs/runtime/README.md](docs/runtime/README.md) | KDSE Runtime usage guide |
| [docs/runtime/COMMANDS.md](docs/runtime/COMMANDS.md) | CLI commands reference |
| [docs/runtime/INITIALIZATION.md](docs/runtime/KDSE_INITIALIZATION.md) | Initialization guide |

## ADRs (Architecture Decision Records)

| ADR | Decision |
|-----|----------|
| [ADR-001](docs/decisions/ADR-001-RUNTIME-IS-THE-AUTHORITY.md) | Runtime is the Authority |
| [ADR-002](docs/decisions/ADR-002-ONE-METHODOLOGY-MULTIPLE-RUNTIMES.md) | One Methodology, Multiple Runtimes |
| [ADR-003](docs/decisions/ADR-003-WORKSPACE-ENGINE.md) | Workspace Engine |

## License

Apache 2.0