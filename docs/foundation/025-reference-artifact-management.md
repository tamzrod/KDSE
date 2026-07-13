# Reference Artifact Management

## Purpose

This document establishes **Reference Artifact Management** as a distinct first-class methodology concept. Reference Artifact Management is responsible for managing engineering evidence before analysis begins.

## The Need for Separation

### The Conceptual Problem

The KDSE methodology previously conflated two distinct responsibilities:

1. **Reference Artifact Management**: Managing evidence before analysis
2. **Collector (Knowledge Derivation)**: Analyzing evidence to derive knowledge

Conflating these responsibilities created ambiguity:

- "What does kdse collect do?" had no clear answer
- "Where does discovery end and analysis begin?" was unclear
- "What is a Collector?" overlapped with "What is a cataloging system?"

### The Solution

Separate these responsibilities into distinct methodology concepts:

| Responsibility | Question Answered | Primary Focus |
|---------------|-------------------|---------------|
| Reference Artifact Management | "What artifacts exist?" | Evidence inventory and classification |
| Collector | "What knowledge can be derived?" | Analysis and derivation |

These responsibilities shall never be merged.

## Definition

### Reference Artifact Management

**Reference Artifact Management** is responsible for managing engineering evidence before analysis begins.

**Definition:** The methodology phase responsible for discovering, cataloging, classifying, and maintaining Reference Artifacts throughout their lifecycle.

**Purpose:** Prepare Reference Artifacts for downstream analysis without interpreting or deriving knowledge from them.

## Responsibilities

Reference Artifact Management has a single, well-defined responsibility:

### Core Responsibilities

#### 1. Discovery

Find Reference Artifacts within the repository and external sources.

**Activities:**
- Scan directory structures
- Search for known artifact patterns
- Identify external artifacts (vendor docs, standards)
- Detect new or modified artifacts

#### 2. Inventory

Record the existence and basic properties of each Reference Artifact.

**Activities:**
- Assign unique identifiers
- Record file paths and metadata
- Track artifact versions
- Maintain inventory state

#### 3. Cataloging

Organize Reference Artifacts into meaningful categories.

**Activities:**
- Classify by artifact type (manual, standard, specification, etc.)
- Group by source or domain
- Establish relationships between artifacts
- Create searchable catalog entries

#### 4. Classification

Determine the nature and quality of each Reference Artifact.

**Activities:**
- Assess artifact type
- Evaluate artifact completeness
- Identify artifact authority level
- Determine artifact applicability

#### 5. Fingerprinting

Establish artifact integrity and identity.

**Activities:**
- Calculate content hashes
- Record modification timestamps
- Track artifact lineage
- Detect unauthorized changes

#### 6. Provenance

Maintain the origin and history of Reference Artifacts.

**Activities:**
- Record artifact sources
- Track acquisition methods
- Document modification history
- Establish authoritative versions

#### 7. Lifecycle Management

Manage Reference Artifacts throughout their existence.

**Activities:**
- Track artifact status (active, deprecated, superseded)
- Archive obsolete artifacts
- Handle artifact replacement
- Maintain artifact retention policies

## What Reference Artifact Management Is Not

Reference Artifact Management intentionally excludes certain responsibilities:

### Not Knowledge Derivation

Reference Artifact Management does NOT:
- Interpret artifact content
- Extract knowledge from artifacts
- Derive Domain Knowledge
- Assess evidence strength
- Identify contradictions
- Generate knowledge artifacts

These responsibilities belong to the Collector.

### Not Architecture

Reference Artifact Management does NOT:
- Organize knowledge into software structures
- Define component boundaries
- Establish architectural patterns
- Make organizational decisions

These responsibilities belong to the Architecture phase.

### Not Implementation

Reference Artifact Management does NOT:
- Realize knowledge using specific technologies
- Select programming languages
- Choose communication protocols
- Implement software

These responsibilities belong to the Implementation phase.

## The Relationship

### Reference Artifact Management and Collector

Reference Artifact Management and Collector have a producer-consumer relationship:

```
┌─────────────────────────────────────────┐
│    Reference Artifact Management        │
│                                         │
│  - Discovers artifacts                  │
│  - Catalogs and classifies              │
│  - Maintains inventory                  │
│  - Preserves provenance                 │
└─────────────────────────────────────────┘
                    │
                    │ Produces
                    ▼
        ┌───────────────────────┐
        │  Cataloged Reference  │
        │     Artifacts         │
        └───────────────────────┘
                    │
                    │ Consumed by
                    ▼
┌─────────────────────────────────────────┐
│           Collector                     │
│                                         │
│  - Analyzes artifacts                   │
│  - Derives Domain Knowledge             │
│  - Correlates evidence                  │
│  - Identifies gaps                      │
└─────────────────────────────────────────┘
```

### The Handoff

The handoff between Reference Artifact Management and Collector is well-defined:

| Reference Artifact Management Produces | Collector Consumes |
|--------------------------------------|-------------------|
| Artifact inventory | Artifact inventory |
| Classification metadata | Classification metadata |
| Provenance records | Provenance records |
| Integrity fingerprints | Integrity fingerprints |
| **NOT interpreted content** | **Interpreted content** |

The Collector receives cataloged artifacts with full provenance. The Collector does not perform discovery or cataloging.

## Reference Artifact Management in the Lifecycle

### Position in the KDSE Lifecycle

Reference Artifact Management is an explicit methodology phase:

```
Reference Artifact
        │
        ▼
┌───────────────────────────┐
│  Reference Artifact       │
│  Management               │
│                           │
│  - Discovery              │
│  - Inventory              │
│  - Cataloging             │
│  - Classification         │
│  - Fingerprinting         │
│  - Provenance             │
│  - Lifecycle              │
└───────────────────────────┘
        │
        ▼
┌───────────────────────────┐
│        Collector          │
│                           │
│  - Reference Analysis     │
│  - Domain Knowledge       │
│    Derivation             │
│  - Evidence Correlation   │
│  - Contradiction         │
│    Detection              │
│  - Knowledge Validation   │
│  - Gap Identification     │
│  - Question Classification│
└───────────────────────────┘
        │
        ▼
  Approved Domain Knowledge
        │
        ▼
     Architecture
        │
        ▼
   Implementation
```

### Phase Boundaries

| Phase | Entrypoint | Exit Criteria |
|-------|------------|---------------|
| Reference Artifact Management | Reference Artifact discovered | Inventory complete, provenance established |
| Collector | Cataloged artifacts received | Knowledge artifacts approved |
| Architecture | Knowledge artifacts finalized | Architecture decisions documented |
| Implementation | Architecture decisions approved | Implementation artifacts verified |

## Reference Artifact Inventory

### Inventory Contents

The Reference Artifact Inventory is the primary output of Reference Artifact Management:

```yaml
inventory:
  session_id: "RAM-2026-07-13-001"
  created_at: "2026-07-13T10:00:00Z"
  artifacts:
    - id: "RA-001"
      path: "artifacts/standards/iec-61850.pdf"
      category: "standard"
      hash: "abc123..."
      source: "vendor"
      discovered: "2026-07-13T10:00:00Z"
      provenance:
        origin: "IEC"
        acquired: "project-init"
        authority: "normative"
    - id: "RA-002"
      path: "artifacts/manuals/ops-guide.pdf"
      category: "manual"
      hash: "def456..."
      source: "vendor"
      discovered: "2026-07-13T10:00:00Z"
      provenance:
        origin: "Vendor Corp"
        acquired: "project-init"
        authority: "informative"
```

### Inventory Properties

Each cataloged Reference Artifact includes:

| Property | Description |
|----------|-------------|
| ID | Unique identifier |
| Path | Repository location |
| Category | Artifact classification |
| Hash | Content integrity fingerprint |
| Source | Origin of the artifact |
| Discovered | Discovery timestamp |
| Provenance | Origin, acquisition, authority |
| Status | Active, deprecated, superseded |

## Common Errors

### Error 1: Conflating Management with Analysis

**Incorrect:**
> "Reference Artifact Management includes analyzing artifacts to extract knowledge."

**Why Incorrect:** Analysis is the Collector's responsibility.

**Correct:**
> "Reference Artifact Management catalogs artifacts for later analysis."

### Error 2: Including Interpretation

**Incorrect:**
> "The Reference Artifact Manager interprets content and identifies key findings."

**Why Incorrect:** Interpretation requires understanding content meaning.

**Correct:**
> "The Reference Artifact Manager records that the artifact exists and classifies it."

### Error 3: Assigning Authority

**Incorrect:**
> "Reference Artifact Management assigns authority levels to artifacts."

**Why Incorrect:** Authority derives from Domain Knowledge, not artifact management.

**Correct:**
> "Reference Artifact Management records the artifact's source provenance, which informs authority assessment."

### Error 4: Deriving Knowledge

**Incorrect:**
> "Reference Artifact Management derives Domain Knowledge from artifacts."

**Why Incorrect:** Knowledge derivation is the Collector's responsibility.

**Correct:**
> "Reference Artifact Management provides cataloged artifacts to the Collector for analysis."

## Implementation Guidance

### Command Responsibility

Reference Artifact Management is a methodology concept, not a command definition.

The methodology does not prescribe:
- How many commands implement Reference Artifact Management
- Whether Reference Artifact Management is a single command or multiple commands
- What the CLI interface looks like

These are implementation decisions.

### Example Implementation Patterns

Reference Artifact Management can be implemented as:

| Pattern | Description |
|---------|-------------|
| Single Command | One command handles all management responsibilities |
| Modular Commands | Separate commands for discovery, cataloging, etc. |
| Integrated | Reference Artifact Management is part of a larger workflow |
| External | Reference Artifact Management happens outside the tool |

The methodology is agnostic to implementation pattern.

## Validation Checklist

When implementing or auditing Reference Artifact Management:

- [ ] Does Reference Artifact Management discover artifacts without analyzing them?
- [ ] Does it maintain a complete artifact inventory?
- [ ] Does it preserve provenance information?
- [ ] Does it classify artifacts without interpreting content?
- [ ] Does it establish artifact integrity through fingerprints?
- [ ] Does it manage artifact lifecycle (active, deprecated, superseded)?
- [ ] Does it defer knowledge derivation to the Collector?
- [ ] Does it avoid mixing responsibilities?

## Relationship to Other Concepts

### Reference Artifact Management and Reference Artifacts

| Concept | Role |
|---------|------|
| Reference Artifact | The evidence being managed |
| Reference Artifact Management | The activity of managing evidence |

### Reference Artifact Management and Domain Knowledge

| Concept | Role |
|---------|------|
| Reference Artifact Management | Prepares artifacts for analysis |
| Domain Knowledge | Is derived by the Collector from artifacts |

### Reference Artifact Management and Architecture

| Concept | Role |
|---------|------|
| Reference Artifact Management | Manages artifacts before any analysis |
| Architecture | Organizes derived knowledge into software structures |

### Reference Artifact Management and Implementation

| Concept | Role |
|---------|------|
| Reference Artifact Management | Manages artifacts before any analysis |
| Implementation | Realizes architecture using specific technologies |

## Summary

Reference Artifact Management is a distinct methodology phase that:

- **Manages** engineering evidence before analysis
- **Discovers** Reference Artifacts without interpreting them
- **Catalogs** and classifies artifacts without deriving knowledge
- **Preserves** provenance and integrity fingerprints
- **Defers** analysis to the Collector
- **Remains independent** from Architecture and Implementation

Understanding this separation is essential for maintaining methodology clarity and avoiding responsibility confusion.

---

## Version

- **Document Version**: 1.0
- **Effective Date**: 2026-07-13
- **Change Note**: Initial release establishing Reference Artifact Management as a distinct methodology concept
