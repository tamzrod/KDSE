# KDSE Normalization

**Document Version:** 1.0  
**Type:** Informative Reference Implementation  
**Effective Date:** 2026-07-12

---

## Purpose

KDSE Normalization enables existing repositories with documentation to adopt the KDSE Runtime without requiring documentation to be rewritten. It creates a normalized engineering layer while preserving all original documentation as the historical engineering record.

---

## Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Existing Repository                              │
│                                                                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  README.md  │  │ architecture│  │    docs/    │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│         │                │                │                        │
│         └────────────────┼────────────────┘                        │
│                          ▼                                           │
│              ┌───────────────────────┐                             │
│              │   Documentation       │                             │
│              │   (Various formats)  │                             │
│              └───────────────────────┘                             │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              │ Normalize
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      KDSE Normalization                             │
│                                                                     │
│  1. Discover documentation                                           │
│  2. Analyze content                                                │
│  3. Extract knowledge                                              │
│  4. Generate KDSE artifacts                                        │
│  5. Build traceability                                             │
│  6. Produce report                                                 │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     KDSE-Normalized Repository                      │
│                                                                     │
│  Original Documentation (Preserved)                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  README.md  │  │ architecture│  │    docs/    │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
│                                                                     │
│  KDSE Engineering Layer (Generated)                                 │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  .kdse/normalized/                                          │   │
│  │                                                             │   │
│  │  ├── SPEC.md         ← KDSE Specification                  │   │
│  │  ├── ARCHITECTURE.md ← KDSE Architecture                   │   │
│  │  ├── GLOSSARY.md     ← KDSE Glossary                       │   │
│  │  ├── DECISIONS.md    ← KDSE Decisions                       │   │
│  │  ├── GOVERNANCE.md   ← KDSE Governance                      │   │
│  │  └── ARTIFACTS.md    ← Artifact Index                       │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Command

### Usage

```bash
kdse normalize
```

### Options

| Option | Description |
|--------|-------------|
| `--help`, `-h` | Show help for normalize command |
| `--types` | Comma-separated list of doc types to include |

### Examples

```bash
# Normalize all documentation
kdse normalize

# Normalize specific documentation types
kdse normalize --types=README,Architecture,Design
```

---

## Process

### Step 1: Discover Documentation

The normalizer scans the repository for documentation files:

| Pattern | Type |
|---------|------|
| `README*`, `*.md` | README, Knowledge |
| `architecture/*`, `ARCHITECTURE*` | Architecture |
| `design/*`, `DESIGN*` | Design |
| `docs/*`, `api/*` | Reference, API |
| `adr/*`, `ADR-*` | Architecture Decisions |
| `wiki/*` | Wiki |
| `examples/*` | Examples |
| `CONTRIBUTING*`, `LICENSE*`, `CHANGELOG*` | Governance |

### Step 2: Analyze Documentation

Each discovered document is analyzed for:

- **Type classification**: README, Architecture, Design, API, etc.
- **Title extraction**: First heading or filename
- **Content analysis**: Patterns, requirements, decisions
- **Confidence scoring**: How confident the system is in the classification

### Step 3: Extract Knowledge

The normalizer extracts structured knowledge:

| Knowledge Type | Description | Sources |
|----------------|-------------|---------|
| Purpose | Project purpose and overview | README |
| Domain | Technical domain | README, Architecture |
| Stakeholders | Authors, maintainers | README, Contributing |
| Requirements | Functional requirements | All documentation |
| Decisions | Architectural decisions | ADRs, Architecture |
| Constraints | Project constraints | Architecture |
| Glossary | Term definitions | All documentation |

### Step 4: Generate KDSE Artifacts

Based on extracted knowledge, the normalizer generates KDSE-standard artifacts:

#### SPEC.md
KDSE-compliant specification document containing:
- Purpose statement
- Domain identification
- Stakeholders
- Requirements
- Constraints
- Traceability block

#### ARCHITECTURE.md
KDSE-compliant architecture document containing:
- System overview
- Design decisions
- Component descriptions
- Traceability block

#### GLOSSARY.md
KDSE-compliant glossary containing:
- Term definitions
- Source attribution
- Traceability block

#### DECISIONS.md
KDSE-compliant decisions index containing:
- Architectural decision records
- Status, context, decision, rationale
- Source attribution

#### GOVERNANCE.md
KDSE-compliant governance document containing:
- Contributing guidelines (reference)
- License information (reference)
- Changelog reference
- Traceability block

#### ARTIFACTS.md
KDSE artifacts index containing:
- List of all normalized artifacts
- Traceability matrix
- Link to original documentation

### Step 5: Build Traceability

Every generated artifact includes full traceability:

```yaml
traceability:
  derived_from:
    - path: "README.md"
      doc_type: "README"
      confidence: "95%"
    - path: "architecture/overview.md"
      doc_type: "Architecture"
      confidence: "90%"
  original_documents:
    - "README.md"
    - "architecture/overview.md"
    - "docs/api.md"
  authority: "Project Team"
  version: "1.0.0"
  normative_refs:
    - "docs/foundation/005-engineering-artifacts.md"
    - "docs/foundation/004-engineering-model.md"
  transformation_type: "Normalization"
```

### Step 6: Produce Report

The normalizer generates a comprehensive normalization report containing:
- Summary statistics
- Discovered documentation list
- Generated artifacts list
- Extracted knowledge summary
- Traceability information
- Next steps recommendations

---

## Non-Destructive Process

### Preservation Guarantee

KDSE Normalization **never modifies** original documentation:

- ✅ Original files remain unchanged
- ✅ Original files remain the historical record
- ✅ KDSE artifacts are created in separate location
- ✅ Operators can reference both views

### Location Separation

```
repository/
├── README.md              ← Original (unchanged)
├── docs/
│   └── api.md             ← Original (unchanged)
├── architecture/
│   └── overview.md        ← Original (unchanged)
└── .kdse/                 ← KDSE artifacts (generated)
    └── normalized/
        ├── SPEC.md
        ├── ARCHITECTURE.md
        ├── GLOSSARY.md
        ├── DECISIONS.md
        ├── GOVERNANCE.md
        └── ARTIFACTS.md
```

---

## Optional Process

KDSE Normalization is **always optional**. The runtime never requires normalization before other commands.

### Command Independence

| Command | Requires Normalize |
|---------|-------------------|
| `kdse run` | No |
| `kdse status` | No |
| `kdse report` | No |
| `kdse audit` | No |
| `kdse normalize` | No (it is the command) |

### Workflow Flexibility

Operators choose their workflow:

```
Workflow A: Immediate Session
  kdse run  →  Start working immediately

Workflow B: Normalize First
  kdse normalize  →  kdse run  →  Work with KDSE artifacts

Workflow C: Normalize Later
  kdse run  →  Work  →  kdse normalize  →  Continue with KDSE artifacts
```

---

## Integration with KDSE Session

After normalization, operators can use `kdse run` to start a KDSE session. The normalized artifacts provide:

- Complete engineering context
- Full traceability to original documentation
- KDSE-standard artifact structure
- Clear derivation chain

### Session Benefits

A normalized repository gains:

| Benefit | Description |
|---------|-------------|
| Context | KDSE-standard artifact descriptions |
| Traceability | Clear derivation from original docs |
| Compliance | KDSE artifact format |
| Organization | Centralized artifact index |

---

## Verification

### Verify Normalization

```bash
# Check generated artifacts
ls -la .kdse/normalized/

# View artifact index
cat .kdse/normalized/ARTIFACTS.md

# View normalization report
kdse normalize  # Re-run or check saved report
```

### Verify Traceability

Each generated artifact contains a traceability section:

```bash
# View traceability in SPEC.md
head -30 .kdse/normalized/SPEC.md
```

---

## Troubleshooting

### No Documentation Found

If normalization finds no documentation:

```
No documentation found in repository.
Consider:
  - Adding README.md
  - Creating docs/ directory
  - Documenting architecture decisions
```

### Partial Results

If some artifacts were not generated:

- Check the normalization report for details
- Verify original documentation content
- Some artifacts require specific documentation types

### Conflicting Documentation

If original documentation conflicts:

- Original documentation remains authoritative
- KDSE artifacts note derived sources
- Operators should resolve conflicts manually

---

## Document Relationships

```
NORMALIZE.md
    │
    ├── Defines: Normalization process, commands, integration
    │
    └── Referenced by:
        ├── EXECUTION_MODEL.md
        ├── COMMANDS.md
        ├── REPORT_SPEC.md
        └── WORKFLOW.md
```

---

## Design Principles

### Non-Destructive
Original documentation is never modified.

### Optional
Normalization never blocks other commands.

### Traceable
Every artifact traces to original sources.

### Standardized
Generated artifacts follow KDSE standards.

### Operator-Controlled
Operators choose when and how to normalize.

---

## Future Enhancements

Potential improvements for future versions:

- Incremental normalization (update only changed docs)
- Conflict detection and resolution
- Validation of extracted knowledge
- Bidirectional sync with original docs
- Custom normalization templates

---

*This document is an informative reference implementation. It describes the KDSE Normalization process, not KDSE requirements.*
