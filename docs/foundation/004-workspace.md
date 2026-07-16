# KDSE Workspace

## Purpose

The `.kdse/` directory is the isolated workspace where KDSE owns and manages all its operational state. It ensures KDSE never pollutes the user's repository with methodology artifacts.

## Isolation Principle

KDSE is a guest in the repository. It operates exclusively within `.kdse/` and must never create files outside this boundary.

```
repository-root/
├── .kdse/           # KDSE owns everything here
│   ├── knowledge/   # Knowledge artifacts
│   ├── evidence/    # Evidence files
│   ├── sessions/    # Session state
│   ├── artifacts/   # Derived artifacts
│   └── agreement.json
├── src/             # User's code - KDSE does not touch
├── docs/            # User's docs - KDSE does not touch
└── tests/           # User's tests - KDSE does not touch
```

## Directory Structure

| Directory | Purpose |
|-----------|---------|
| `knowledge/` | Structured knowledge artifacts derived from evidence |
| `evidence/` | Raw evidence files (screenshots, benchmarks, test results) |
| `sessions/` | Session state and handoff information |
| `artifacts/` | Derived artifacts (normalized docs, reports) |
| `runtime/` | Runtime verification and invariants |
| `reports/` | Generated reports and metrics |

## Key Files

### agreement.json

The Current Agreement captures project state between sessions:

```json
{
  "project_name": "my-project",
  "current_phase": "Problem",
  "methodology_version": "1.0",
  "runtime_version": "1.0.0",
  "assumptions": [],
  "context": {}
}
```

### Workspace Initialization

```bash
kdse init
```

Creates the `.kdse/` directory structure if it doesn't exist.

### Workspace Status

```bash
kdse status
```

Shows current workspace state:
- Agreement status
- Active sessions
- Knowledge artifact count
- Evidence count

## Migration

If legacy KDSE directories exist at the repository root (from older versions), migrate them:

```bash
kdse migrate
```

This moves:
- `foundation/` → `.kdse/foundation/`
- `knowledge/` → `.kdse/knowledge/`
- `context/` → `.kdse/context/`
- `artifacts/` → `.kdse/artifacts/`

## Enforcement

The workspace package enforces isolation through:

1. **Path Resolution**: All KDSE paths are resolved relative to `.kdse/`
2. **Boundary Checks**: Operations are rejected if they would write outside `.kdse/`
3. **Migration Reports**: Legacy directories are detected and reported

## Session Handoff

When transferring between sessions, the workspace captures:

1. Current agreement state
2. Active knowledge artifacts
3. Pending evidence
4. Phase context

This enables delta-based communication—subsequent sessions load the agreement and resume from where the previous session ended.

## Metrics

The workspace tracks:
- Artifact counts by type
- Session duration
- Knowledge derivation velocity
- Traceability coverage

These metrics are stored in `.kdse/reports/` and can be queried via CLI.
