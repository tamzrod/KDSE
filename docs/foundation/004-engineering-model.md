# Engineering Model

## Lifecycle Overview

KDSE defines a lifecycle in which artifacts progress through defined stages. Each stage produces artifacts that serve as input to subsequent stages.

```
Reference Artifacts
    ↓
Reference Analysis
    ↓
Engineering Knowledge Derivation
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

Collect and organize existing sources of engineering information that will serve as evidence for knowledge derivation.

### Input

- Project documentation
- Implementation artifacts
- Vendor documentation
- Standards and specifications
- Commissioning records
- Communication artifacts

### Output

- Catalog of Reference Artifacts
- Artifact inventory with relevance assessment
- Context documentation

### Activities

- Identify available Reference Artifacts
- Assess relevance to engineering domain
- Document artifact provenance
- Organize artifacts for analysis

### Why This Stage Comes First

Reference Artifacts are the evidence from which Engineering Knowledge is derived. Without Reference Artifacts, knowledge derivation has no foundation.

## Stage 1: Reference Analysis

### Purpose

Examine Reference Artifacts to identify engineering evidence that supports knowledge derivation.

### Input

- Reference Artifacts from Stage 0
- Engineering domain scope

### Output

- Identified evidence from each artifact
- Evidence categorization
- Context for each piece of evidence

### Activities

- Analyze each Reference Artifact
- Extract factual statements, assertions, and decisions
- Identify constraints and requirements
- Document evidence with provenance

## Stage 2: Engineering Knowledge Derivation

### Purpose

Transform analyzed evidence into implementation-independent Engineering Knowledge.

### Input

- Evidence from Stage 1
- Engineering domain scope

### Output

- Engineering Knowledge statements
- Evidence Strength assessments
- Traceability links

### Activities

- Interpret evidence for engineering meaning
- Apply Engineering Independence Test
- Formulate implementation-independent statements
- Assess Evidence Strength
- Correlate evidence across sources
- Identify and preserve contradictions

### Why This Stage Follows Reference Analysis

Engineering Knowledge must be derived, not simply extracted. Analysis reveals evidence; derivation transforms evidence into authoritative knowledge.

## Stage 3: Architecture

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

### Why This Stage Comes After Knowledge Derivation

Architecture bridges knowledge and implementation. Architecture translates abstract knowledge into concrete structural decisions. Implementation without architecture guidance produces inconsistent systems.

## Stage 4: Implementation

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

## Stage 5: Verification

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

## Stage 6: Evolution

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
│  Reference Artifacts → Analysis →      │
│  Engineering Knowledge Derivation       │
│                                       │
└───────────────────────────────────────┘
```

## Derivation Principle

Each stage derives from its predecessor:

- Reference Analysis derives from Reference Artifacts
- Engineering Knowledge Derivation derives from Reference Analysis
- Architecture derives from Engineering Knowledge
- Implementation derives from Architecture
- Verification derives from Implementation (and confirms alignment with Knowledge and Architecture)

Derivation is not optional. Artifacts produced in later stages must trace to authorized artifacts from earlier stages.

## Authority Principle

Authority flows downward:

- Reference Artifacts provide evidence (not authority)
- Engineering Knowledge derives from evidence and carries authority
- Architecture derives authority from Knowledge
- Architecture authorizes Implementation
- Implementation is subject to Verification

Lower stages cannot contradict higher stages. Violations represent methodology non-compliance.

## Separation of Concerns

KDSE maintains strict separation between:

| Concern | Description |
|---------|-------------|
| Reference Artifact | Engineering evidence (project docs, implementation, vendor docs) |
| Engineering Knowledge | Implementation-independent understanding |
| Architecture | Organization of Engineering Knowledge into software |
| Implementation | Realization of Architecture using specific technologies |

These concepts remain independent throughout the entire methodology.
