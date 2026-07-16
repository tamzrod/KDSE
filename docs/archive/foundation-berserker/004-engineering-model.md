# Engineering Model

## Lifecycle Overview

KDSE defines a lifecycle in which artifacts progress through defined stages. Each stage produces artifacts that serve as input to subsequent stages.

```
Reference Artifacts
    ↓
Reference Artifact Management
    ↓
Reference Analysis
    ↓
Domain Knowledge Derivation
    ↓
Evidence Correlation
    ↓
Knowledge Validation
    ↓
Approved Domain Knowledge
    ↓
Architecture
    ↓
Implementation
    ↓
Verification
    ↓
Evolution
```

## Stage 0: Reference Artifacts

### Purpose

Raw engineering evidence that exists in various forms throughout the project and external sources.

### Input

- Project documentation
- Implementation artifacts
- Vendor documentation
- Standards and specifications
- Commissioning records
- Communication artifacts

### Output

- Raw artifacts (unchanged)

### Note

Reference Artifacts are the evidence from which Domain Knowledge is derived. Without Reference Artifacts, knowledge derivation has no foundation.

## Stage 1: Reference Artifact Management

### Purpose

Discover, catalog, classify, and manage Reference Artifacts before analysis begins.

### Input

- Reference Artifacts from Stage 0

### Output

- Catalog of Reference Artifacts
- Artifact inventory with classification
- Provenance records
- Integrity fingerprints
- Lifecycle status

### Activities

- Discover artifacts within the repository and external sources
- Assign unique identifiers
- Classify artifacts by type (manual, standard, specification, etc.)
- Record provenance (origin, acquisition method)
- Calculate content hashes for integrity
- Track lifecycle status (active, deprecated, superseded)

### Responsibilities

Reference Artifact Management:
- **DOES**: Discover, catalog, classify, maintain inventory, preserve provenance
- **DOES NOT**: Interpret content, derive knowledge, analyze meaning

### Why This Stage Is Explicit

Reference Artifact Management is distinct from the Collector. Managing evidence (inventory, classification, provenance) is fundamentally different from analyzing evidence (deriving knowledge). This separation ensures clarity of responsibility.

## Stage 2: Reference Analysis

### Purpose

Examine cataloged Reference Artifacts to identify engineering evidence that supports knowledge derivation.

### Input

- Cataloged Reference Artifacts from Stage 1
- Engineering domain scope

### Output

- Identified evidence from each artifact
- Evidence categorization
- Context for each piece of evidence

### Activities

- Analyze each cataloged Reference Artifact
- Extract factual statements, assertions, and decisions
- Identify constraints and requirements
- Document evidence with provenance

## Stage 3: Domain Knowledge Derivation

### Purpose

Transform analyzed evidence into implementation-independent Domain Knowledge.

### Input

- Evidence from Stage 2
- Engineering domain scope

### Output

- Domain Knowledge statements
- Evidence Strength assessments
- Traceability links

### Activities

- Interpret evidence for domain meaning
- Apply Engineering Independence Test
- Formulate implementation-independent statements
- Assess Evidence Strength
- Identify domain behavior
- Define domain constraints
- Document assumptions

### Why This Stage Follows Reference Analysis

Domain Knowledge must be derived, not simply extracted. Analysis reveals evidence; derivation transforms evidence into authoritative knowledge.

## Stage 4: Evidence Correlation

### Purpose

Strengthen Domain Knowledge through multiple independent Reference Artifacts.

### Input

- Domain Knowledge statements from Stage 3
- Multiple Reference Artifacts

### Output

- Correlated Domain Knowledge
- Evidence Strength assignments
- Identified contradictions

### Activities

- Compare evidence across multiple artifacts
- Identify corroborating evidence
- Identify contradicting evidence
- Assign Evidence Strength based on corroboration
- Preserve contradictions (do not resolve silently)

### Evidence Strength Scale

- ★★★★★: Supported by multiple independent sources
- ★★★★☆: Supported by Project Doc + one additional source
- ★★★☆☆: Supported by Project Documentation only
- ★★☆☆☆: Supported by single source or vendor only
- ★☆☆☆☆: Inferred from indirect evidence

## Stage 5: Knowledge Validation

### Purpose

Validate that Domain Knowledge meets quality criteria before approval.

### Input

- Correlated Domain Knowledge from Stage 4

### Output

- Validated Domain Knowledge
- Validation report
- Identified gaps

### Activities

- Review traceability links
- Verify Engineering Independence Test passes
- Confirm Evidence Strength is appropriate
- Check completeness
- Identify remaining gaps

## Stage 6: Approved Domain Knowledge

### Purpose

Domain Knowledge that has passed validation and is authorized for downstream use.

### Input

- Validated Domain Knowledge from Stage 5

### Output

- Approved Domain Knowledge artifacts
- Authority records

## Stage 7: Architecture

### Purpose

Translate knowledge into structural decisions that guide implementation.

### Input

- Authorized knowledge artifacts
- Previously established architecture (for continuation)

### Output

- Architecture artifacts defining structure
- Architecture artifacts defining constraints
- Architecture Decision Records documenting decisions

### Activities

- Derive architecture from knowledge
- Document architectural approaches
- Record decisions with rationale
- Review architecture against knowledge
- Authorize architecture artifacts

### Why This Stage Comes After Knowledge Validation

Architecture bridges knowledge and implementation. Architecture translates abstract knowledge into concrete structural decisions. Implementation without architecture guidance produces inconsistent systems.

## Stage 8: Implementation

### Purpose

Realize architecture through code that traces to authorized artifacts.

### Input

- Authorized architecture artifacts
- Architecture Decision Records
- Previously established implementation (for continuation)

### Output

- Implementation artifacts (code, configuration, data)
- Implementation-specific documentation
- Traceability records linking implementation to architecture

### Activities

- Implement according to architecture
- Maintain traceability to architecture
- Resolve implementation questions against architecture
- Document implementation decisions not requiring architecture change

### Why This Stage Comes After Architecture

Implementation realizes architecture. Implementation decisions must align with architectural direction. Implementation that contradicts architecture represents drift requiring correction.

## Stage 9: Verification

### Purpose

Confirm that implementation aligns with architecture and that architecture aligns with knowledge.

### Input

- Authorized knowledge artifacts
- Authorized architecture artifacts
- Implementation artifacts

### Output

- Verification artifacts confirming alignment
- Reports documenting verification results
- Identified gaps requiring correction

### Activities

- Derive verification criteria from knowledge
- Execute verification against implementation
- Document verification results
- Identify and report misalignments
- Determine root cause of identified misalignments
- Classify misalignments by source artifact (Knowledge, Architecture, or Implementation)

### Why This Stage Comes After Implementation

Verification confirms alignment. Without prior stages producing artifacts, verification has no basis. Verification cannot confirm alignment that does not exist.

## Stage 10: Evolution

### Purpose

Manage changes to artifacts over time while maintaining traceability and authority.

### Input

- All prior stage artifacts
- Change requests
- Feedback from operation

### Output

- Updated artifacts reflecting changes
- Records of artifact evolution
- Maintained traceability chains

### Activities

- Evaluate change requests against artifact hierarchy
- Propagate changes through stages as necessary
- Maintain artifact versioning
- Retire obsolete artifacts
- Archive knowledge for future reference

### Why This Stage Is Continuous

Software systems evolve. Artifacts must evolve with them. Evolution is not a terminal stage but a continuous process that may trigger re-entry to any prior stage.

## Stage Relationships

```
        ┌───────────────────────────────────────┐
        │                                       │
        │           Evolution                   │
        │                                       │
        └───────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
        ▼               ▼               ▼
    Knowledge      Architecture   Implementation
        │               │               │
        │               │               │
        └───────────────┴───────────────┘
                        │
                        ▼
                Verification
        │
        ▼
┌───────────────────────────────────────┐
│                                       │
│  Reference Artifacts                   │
│       ↓                               │
│  Reference Artifact Management          │
│       ↓                               │
│  Reference Analysis                     │
│       ↓                               │
│  Domain Knowledge Derivation           │
│       ↓                               │
│  Evidence Correlation                  │
│       ↓                               │
│  Knowledge Validation                   │
│       ↓                               │
│  Approved Domain Knowledge             │
│                                       │
└───────────────────────────────────────┘
```

## Derivation Principle

Each stage derives from its predecessor:

- Reference Artifact Management derives from Reference Artifacts
- Reference Analysis derives from cataloged Reference Artifacts
- Domain Knowledge Derivation derives from Reference Analysis
- Evidence Correlation derives from Domain Knowledge Derivation
- Knowledge Validation derives from Evidence Correlation
- Architecture derives from Approved Domain Knowledge
- Implementation derives from Architecture
- Verification derives from Implementation (and confirms alignment with Knowledge and Architecture)

Derivation is not optional. Artifacts produced in later stages must trace to authorized artifacts from earlier stages.

## Authority Principle

Authority flows downward:

- Reference Artifacts provide evidence (not authority)
- Reference Artifact Management preserves evidence (not authority)
- Domain Knowledge derives from evidence and carries authority
- Architecture derives authority from Knowledge
- Architecture authorizes Implementation
- Implementation is subject to Verification

Lower stages cannot contradict higher stages. Violations represent methodology non-compliance.

## Separation of Concerns

KDSE maintains strict separation between:

| Concern | Description |
|---------|-------------|
| Reference Artifact | Raw engineering evidence |
| Reference Artifact Management | Discovery, cataloging, classification, provenance (does NOT analyze) |
| Domain Knowledge | Implementation-independent understanding |
| Architecture | Organization of Domain Knowledge into software |
| Implementation | Realization of Architecture using specific technologies |

These concepts remain independent throughout the entire methodology.
