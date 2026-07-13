# Collector Philosophy

## Purpose

This document establishes the philosophy for **Collectors** in KDSE. A Collector is a methodology component defined by its responsibility, not by the type of Reference Artifact it analyzes.

**Important:** Collectors consume cataloged Reference Artifacts from Reference Artifact Management. They do not discover or catalog artifacts themselves.

## The Problem with Artifact-Type Definitions

### Traditional Approach

Traditional methodology often defines collectors by the artifact type they handle:

- "A requirements collector gathers requirements"
- "A design collector gathers design documents"
- "A code collector gathers source code"

### Problems with This Approach

1. **Boundaries Are Artificial**: Artifact types overlap and blur
2. **Responsibilities Are Unclear**: What exactly does each collector do?
3. **Duplication Occurs**: Multiple collectors may handle the same information
4. **Flexibility Is Limited**: New artifact types require new collectors

## The Responsibility-Based Approach

### Collector Definition

A **Collector** is a methodology component that:

1. **Consumes cataloged Reference Artifacts** (from Reference Artifact Management)
2. **Performs Reference Analysis** on artifacts
3. **Derives implementation-independent Domain Knowledge**
4. **Preserves traceability**
5. **Correlates evidence**
6. **Identifies contradictions**
7. **Validates Domain Knowledge**
8. **Identifies Domain Knowledge gaps**
9. **Classifies questions** for operator interaction

### What Collectors Do NOT Do

Collectors do NOT:
- **Discover Reference Artifacts** (this is Reference Artifact Management)
- **Catalog or classify artifacts** (this is Reference Artifact Management)
- **Maintain artifact inventory** (this is Reference Artifact Management)
- **Establish provenance** (this is Reference Artifact Management)

### Collector Definition By Responsibility

A Collector is defined by what it does, not what it touches.

The same Collector can analyze:

- Project documentation
- Implementation artifacts
- Vendor documentation
- Standards and specifications
- Communication artifacts

Because the Collector's responsibility is analysis and derivation, not artifact-type-specific collection.

## What Collectors Do

### Core Responsibilities

#### 1. Analyze Reference Artifacts

Collectors examine Reference Artifacts to identify relevant information.

**Activities**:

- Read and interpret artifacts
- Identify relevant sections
- Extract factual statements
- Note assertions and claims
- Document context

#### 2. Identify Engineering Evidence

Collectors identify evidence that supports Engineering Knowledge.

**Activities**:

- Distinguish facts from opinions
- Identify explicit claims
- Note implicit assumptions
- Capture constraints
- Document decisions

#### 3. Derive Implementation-Independent Knowledge

Collectors transform evidence into Engineering Knowledge.

**Activities**:

- Apply the Engineering Independence Test
- Abstract implementation details
- Formulate knowledge statements
- Ensure traceable derivation

#### 4. Preserve Traceability

Collectors maintain links between evidence and knowledge.

**Activities**:

- Document evidence sources
- Link knowledge to evidence
- Track dependencies
- Maintain provenance

#### 5. Correlate Evidence

Collectors strengthen knowledge through multiple sources.

**Activities**:

- Compare evidence across artifacts
- Identify agreements
- Identify contradictions
- Assign Evidence Strength

#### 6. Identify Contradictions

Collectors preserve contradictions rather than resolving them.

**Activities**:

- Document conflicting claims
- Assess engineering impact
- Classify resolution path
- Flag for operator review

#### 7. Identify Gaps

Collectors identify areas where knowledge is incomplete.

**Activities**:

- Detect missing evidence
- Note unresolved questions
- Classify gap type
- Flag for resolution

## What Collectors Do Not Do

### Excluded Responsibilities

Collectors do not:

1. **Generate Software Documentation**: Creating documentation is a separate activity
2. **Design Software**: Architecture is a separate phase
3. **Generate Architecture**: Architecture derives from Engineering Knowledge
4. **Implement Software**: Implementation is a separate phase

### Why These Are Excluded

These responsibilities belong to different phases and roles:

| Excluded Responsibility | Belongs To |
|------------------------|-----------|
| Discover Reference Artifacts | Reference Artifact Management |
| Catalog Reference Artifacts | Reference Artifact Management |
| Maintain artifact inventory | Reference Artifact Management |
| Generate documentation | Documentation practices |
| Design software | Architecture phase |
| Generate architecture | Architecture phase |
| Implement software | Implementation phase |

## The Collector in the KDSE Lifecycle

### Position in Lifecycle

Collectors consume cataloged Reference Artifacts from Reference Artifact Management:

```
┌─────────────────────────────────────────┐
│    Reference Artifact Management        │
│                                         │
│  - Discovery                           │
│  - Inventory                           │
│  - Cataloging                          │
│  - Classification                      │
│  - Provenance                          │
└─────────────────────────────────────────┘
                    │
                    │ Produces cataloged artifacts
                    ▼
┌─────────────────────────────────────────┐
│           Collector                     │
│                                         │
│  - Reference Analysis                  │
│  - Domain Knowledge Derivation         │
│  - Evidence Correlation                │
│  - Contradiction Detection             │
│  - Knowledge Validation                │
│  - Gap Identification                 │
│  - Question Classification             │
└─────────────────────────────────────────┘
                    │
                    ▼
        Approved Domain Knowledge
```

### The Handoff

The handoff between Reference Artifact Management and Collector is well-defined:

| Reference Artifact Management Produces | Collector Consumes |
|--------------------------------------|-------------------|
| Artifact inventory | Artifact inventory |
| Classification metadata | Classification metadata |
| Provenance records | Provenance records |
| Integrity fingerprints | Integrity fingerprints |
| **NOT analyzed content** | **Analyzed content** |

### Collector Workflow

```
┌─────────────────────────────────────┐
│         Collector Workflow          │
├─────────────────────────────────────┤
│ 1. Receive cataloged artifacts      │
│ 2. Perform Reference Analysis       │
│ 3. Identify evidence                 │
│ 4. Derive Domain Knowledge          │
│ 5. Apply Independence Test          │
│ 6. Correlate evidence                │
│ 7. Detect contradictions             │
│ 8. Validate knowledge               │
│ 9. Identify gaps                     │
│ 10. Classify questions               │
│ 11. Produce knowledge artifacts     │
└─────────────────────────────────────┘
```

## Relationship to Reference Artifact Management

### Distinction

| Aspect | Reference Artifact Management | Collector |
|--------|------------------------------|-----------|
| Question Answered | "What artifacts exist?" | "What knowledge can be derived?" |
| Focus | Evidence inventory and classification | Evidence analysis and derivation |
| Input | Raw artifacts | Cataloged artifacts |
| Output | Artifact inventory | Domain Knowledge artifacts |
| Content Interpretation | None | Full interpretation |

### Why This Separation Matters

1. **Clarity**: Each phase has exactly one primary responsibility
2. **Reusability**: The artifact inventory can be used by multiple collectors
3. **Separation of Concerns**: Discovery and analysis are independent
4. **Testing**: Each phase can be tested independently
5. **Scalability**: Different approaches can be used for each phase

## Artifact Type Independence

### Why Collectors Are Artifact-Type Independent

A Collector analyzes any Reference Artifact that contains engineering evidence.

This includes:

| Artifact Type | Example | Analysis Focus |
|---------------|---------|----------------|
| Documentation | P&ID, manuals | Extracting specifications |
| Implementation | Source code, configs | Understanding realized behavior |
| Communication | Meeting notes, emails | Capturing decisions |
| Standards | Grid codes, specs | Identifying requirements |
| Commissioning | Test records | Verifying behavior |

### Benefits of Independence

1. **Flexibility**: New artifact types can be added without new collectors
2. **Consistency**: Same analysis approach regardless of source
3. **Completeness**: All artifacts can be analyzed uniformly
4. **Simplicity**: Single collector definition

## Collector Qualifications

### Required Capabilities

A Collector must be capable of:

1. **Reading and interpreting** various artifact types
2. **Identifying engineering relevance** in technical content
3. **Applying the Engineering Independence Test**
4. **Assessing Evidence Strength**
5. **Formulating knowledge statements**
6. **Maintaining traceability**

### Not Required

A Collector does not require:

1. **Programming skills** (unless analyzing code)
2. **Specific tool expertise** (tool support is implementation)
3. **Domain expertise** (unless deriving complex knowledge)
4. **Architecture skills** (architecture is a separate phase)

## Collector Output

### Primary Output: Engineering Knowledge

Collectors produce Engineering Knowledge Artifacts:

```
| Knowledge ID | Statement | Evidence | Strength | Status |
|--------------|-----------|----------|----------|--------|
| EK-001 | System supports grid-forming | RA-001, RA-002, RA-003 | ★★★★★ | Approved |
| EK-002 | Grid-forming via switch | RA-001, RA-002 | ★★★★☆ | Draft |
```

### Secondary Output: Analysis Reports

Collectors produce analysis reports:

```
## Analysis Report: Control System

### Evidence Identified
- RA-001: Project P&ID (5 statements)
- RA-002: Node-RED implementation (3 statements)
- RA-003: Vendor manual (2 statements)

### Contradictions
- Default mode: RA-001 vs RA-002 (Low impact)

### Gaps
- Fault response behavior not documented

### Recommendations
- Request operator clarification on default mode
- Search for fault response documentation
```

### Gap Reports

Collectors identify gaps for resolution:

```
## Gap Report: Control System

| Gap ID | Description | Type | Resolution Path |
|--------|-------------|------|-----------------|
| GAP-001 | Fault response time unknown | Engineering Knowledge | Ask operator |
| GAP-002 | Redundancy architecture unclear | Architecture | Defer to Architecture |
| GAP-003 | Modbus register mapping | Implementation | Repository discovery |
```

## Collector-Operator Interaction

### When to Ask the Operator

Collectors ask the operator only when:

1. **Engineering Knowledge gaps** cannot be resolved from artifacts
2. **Contradictions** affect Engineering Knowledge
3. **High-value questions** that cannot be deferred

### Repository First Principle

Before asking the operator, collectors must:

1. Search all available Reference Artifacts
2. Analyze existing implementation
3. Examine project documentation
4. Review vendor materials

### Question Classification

Before asking, classify the question:

| Classification | Action |
|----------------|--------|
| Engineering Knowledge | Ask during Knowledge Derivation |
| Architecture | Defer to Architecture Phase |
| Implementation | Attempt repository discovery first |

## Implementation Concerns

### Tool Support vs Methodology

Support for specific Reference Artifact formats is an implementation concern.

| Methodology Concern | Implementation Concern |
|--------------------|-----------------------|
| Collector responsibility | How to parse PDF documents |
| Knowledge derivation process | How to store knowledge artifacts |
| Evidence correlation | How to search artifacts |
| Traceability requirements | Tool selection and configuration |

### Collector Implementation Examples

Possible implementations include:

- **AI-Assisted Collectors**: Use AI to analyze artifacts and suggest knowledge
- **Rule-Based Collectors**: Apply structured rules to extract information
- **Hybrid Collectors**: Combine AI assistance with human validation
- **Manual Collectors**: Human analysts following methodology

All implementations serve the same responsibility; they differ only in execution.

## Collector Validation

### Validation Criteria

Collector output is validated against:

1. **Traceability**: All knowledge traces to evidence
2. **Independence**: All knowledge passes Engineering Independence Test
3. **Corroboration**: Evidence Strength is assessed
4. **Completeness**: All relevant artifacts analyzed
5. **Contradictions Preserved**: No silent resolutions

### Validation Process

```
┌─────────────────────────────────────┐
│       Validation Process           │
├─────────────────────────────────────┤
│ 1. Review traceability links        │
│ 2. Verify independence test passes  │
│ 3. Confirm evidence strength        │
│ 4. Check completeness               │
│ 5. Verify contradictions preserved   │
└─────────────────────────────────────┘
```

## Common Errors

### Error 1: Defining Collectors by Artifact Type

**Incorrect**:
> "We need a code collector for Node-RED flows and a documentation collector for P&IDs."

**Why Incorrect**: This creates artificial boundaries and duplication.

**Correct**:
> "We need a collector that can analyze any Reference Artifact to derive Engineering Knowledge."

### Error 2: Assigning Architecture Responsibilities

**Incorrect**:
> "The collector should also generate the software architecture."

**Why Incorrect**: Architecture is a separate phase with separate responsibilities.

**Correct**:
> "The collector derives Engineering Knowledge. Architecture derives from that knowledge."

### Error 3: Including Implementation Decisions

**Incorrect**:
> "The collector found that the system uses Modbus TCP on port 502."

**Why Incorrect**: This is implementation detail, not Engineering Knowledge.

**Correct**:
> "The collector found that sensor data is exchanged between components. Implementation technology is a separate finding."

## Summary

The Collector philosophy establishes that:

- **Collectors are defined by responsibility, not artifact type**
- **Collectors consume cataloged Reference Artifacts from Reference Artifact Management**
- **Collectors do NOT discover or catalog artifacts (this is Reference Artifact Management)**
- **Collectors analyze Reference Artifacts to derive Domain Knowledge**
- **Collectors preserve traceability, correlate evidence, and identify gaps**
- **Collectors do not generate documentation, design software, or implement systems**
- **Support for specific formats is implementation, not methodology**
- **Collectors ask operators only after repository-first analysis**

Understanding this philosophy ensures consistent, effective knowledge derivation.

---

## Version

- **Document Version**: 2.0
- **Effective Date**: 2026-07-13
- **Change Note**: Separated Collector from Reference Artifact Management; Collector now consumes cataloged artifacts rather than discovering them
