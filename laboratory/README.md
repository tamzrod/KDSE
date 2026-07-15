# KDSE Laboratory

## Purpose

The Laboratory validates the KDSE methodology. It does **NOT** test software; it validates engineering methodology.

## Structure

```
laboratory/
├── README.md              # This file
├── KDSE_PRINCIPLES.md    # Constitutional principles
├── core/                 # KDSE Slim implementation
│   ├── artifacts/
│   │   ├── registry.py       # Artifact Registry
│   │   ├── decision_engine.py # Decision Engine (STATE → GAPS → DECIDE)
│   │   ├── derivation.py     # Derivation Orchestrator (GATHER → DERIVE → VALIDATE)
│   │   └── foundation.py     # Foundation Skeleton Generator
│   ├── laboratory.py        # Laboratory Runner
│   └── __init__.py
├── scenarios/            # Laboratory scenarios
│   └── LAB-001.md        # Inventory Management System
├── reports/             # Experiment reports
│   └── LAB-001-REPORT.md
├── results/             # Machine-readable results
│   └── LAB-001-REPORT.json
└── baselines/           # Baseline expectations
```

## Quick Start

### Run LAB-001

```bash
python3 run_lab.py --scenario LAB-001
```

This will:
1. Create Foundation Skeleton (6 documents)
2. Identify Known Facts
3. Identify Assumptions
4. Identify Knowledge Gaps
5. Record Evidence Strength
6. Generate report

## Constitutional Principles

The Laboratory enforces these non-negotiable principles:

| Principle | Description |
|-----------|-------------|
| **Knowledge Precedes Architecture** | Every architectural decision traces to knowledge |
| **Authority Flows Downward** | Lower layers cannot contradict higher |
| **Evidence Supports, Never Authorizes** | Evidence ≠ Authority |
| **Traceability Is Absolute** | Every decision traces to knowledge |
| **Knowledge Is Implementation-Independent** | Survives technology changes |
| **Verification Is Continuous** | Alignment must be continuously checked |

## Core Components

### Artifact Registry

Tracks all engineering artifacts by type, state, and owner.

**Lifecycle States:**
- DRAFT → REVIEWED → APPROVED → ARCHIVED

**Artifact Types:**
- KNOWLEDGE (highest authority)
- ARCHITECTURE (traces to knowledge)
- IMPLEMENTATION (traces to architecture)
- VERIFICATION (confirms alignment)
- EVIDENCE (supports knowledge)

### Decision Engine

Evaluates state, identifies gaps, and produces decisions.

**Three-Step Cycle:**
1. **STATE** - What exists? What traces to what?
2. **GAPS** - What is missing? What contradicts?
3. **DECIDE** - What should happen next?

### Derivation Orchestrator

Guides the knowledge derivation process.

**Three Stages:**
1. **GATHER** - Find and classify evidence sources
2. **DERIVE** - Transform evidence into knowledge statements
3. **VALIDATE** - Check derivation quality, assign strength

### Evidence Strength

All knowledge has evidence strength:

- **●●● (STRONG)** - Multiple independent sources confirm
- **●● (MODERATE)** - Some corroboration exists
- **● (WEAK)** - Single source or inference only

### Foundation Skeleton

Created immediately after receiving an objective:

- **PROBLEM.md** - What problem are we solving?
- **SPEC.md** - What is the project specification?
- **REQUIREMENTS.md** - What are the functional requirements?
- **ASSUMPTIONS.md** - What are the reasonable assumptions?
- **CONSTRAINTS.md** - What are the constraints?
- **GLOSSARY.md** - What domain terminology is used?

## Knowledge Distinction

KDSE distinguishes between:

| Category | Definition | Evidence |
|---------|------------|----------|
| **Known Facts** | Backed by multiple evidence | ●●● or ●● |
| **Assumptions** | Documented, to be validated | ● |
| **Knowledge Gaps** | Unknown, need investigation | None |

**Critical Rule:** Nothing silently becomes a fact.

## Running Scenarios

### Run a Single Scenario

```bash
python3 run_lab.py --scenario LAB-001
```

### Run with Custom Workspace

```bash
python3 run_lab.py --scenario LAB-001 --workspace /path/to/workspace
```

### Output

Reports are generated in:
- `results/LAB-001-REPORT.json` - Machine-readable
- `results/LAB-001-REPORT.md` - Human-readable

## Creating New Scenarios

Create a scenario file in `scenarios/`:

```markdown
# LAB-002: Your Scenario

**Scenario ID:** LAB-002  
**Title:** Your Title  
**Objective:** Your objective

## Expected Behavior

### KDSE SHOULD:
✓ Create Foundation Skeleton
✓ [Other expectations]

### KDSE MUST NOT:
✗ Generate source code
✗ [Other prohibitions]

## Success Criteria

| Criterion | Requirement |
|-----------|-------------|
| [Criterion] | [Requirement] |
```

Then add to `run_lab.py`:

```python
def create_scenario(scenario_id: str):
    if scenario_id == "LAB-002":
        return LaboratoryScenario(
            id="LAB-002",
            title="Your Title",
            objective="Your objective",
            # ...
        )
```

## Validation Criteria

**PASS:** All critical checkpoints passed AND no forbidden behavior detected

**FAIL:** Any critical checkpoint failed OR forbidden behavior detected

## Success Metrics

The experiment PASSES only if KDSE behaves according to its methodology:

| Metric | Target |
|--------|--------|
| Foundation Skeleton | 6 documents created |
| Known Facts | Distinguished from assumptions |
| Assumptions | Explicitly marked |
| Knowledge Gaps | Documented with investigation |
| Evidence Strength | Assigned to all knowledge |
| No Premature Implementation | No source code, frameworks, databases |

---

*Laboratory validates methodology, not software.*
