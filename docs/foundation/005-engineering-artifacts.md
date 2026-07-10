# Engineering Artifacts

## Overview

KDSE defines six canonical artifact types. Each artifact type has defined purpose, ownership, lifetime, authority, dependencies, and deliverables.

## Artifact Type 1: Knowledge

### Purpose

Capture and authorize the knowledge necessary to address a problem space.

### Description

Knowledge artifacts represent authoritative understanding about the problem domain, requirements, constraints, and context. Knowledge artifacts are the highest-authority artifact type in KDSE.

### Owner

Knowledge Owner, as designated by governance

### Lifetime

Persists until explicitly retired. Knowledge may be archived but is never deleted from traceability records.

### Authority

Highest authority in the artifact hierarchy. All other artifacts derive authority from knowledge.

### Dependencies

None (originates the artifact hierarchy)

### Deliverables

- Requirements knowledge
- Non-functional requirements knowledge
- Domain knowledge
- Context knowledge
- Assumption knowledge

### Characteristics

- Authoritative
- Traceable
- Versioned
- Reviewed and approved

## Artifact Type 2: Architecture

### Purpose

Translate knowledge into structural decisions that guide implementation.

### Description

Architecture artifacts represent the structural decisions derived from knowledge. Architecture artifacts define how the system is structured, what components exist, and how they interact.

### Owner

Architecture Owner, as designated by governance

### Lifetime

Persists for the system lifecycle. May evolve with version control.

### Authority

Derives authority from knowledge. Authorizes implementation.

### Dependencies

- Knowledge artifacts (required)

### Deliverables

- System architecture documentation
- Component architecture documentation
- Interface definitions
- Architecture Decision Records (ADRs)
- Architectural constraints

### Characteristics

- Derived from knowledge
- Authorizes implementation
- Traceable to knowledge
- Versioned

## Artifact Type 3: Architecture Decision Record (ADR)

### Purpose

Document significant architectural decisions with rationale and consequences.

### Description

ADRs capture architectural decisions that merit documentation. An ADR captures what was decided, why it was decided, and what alternatives were considered.

### Owner

Architecture Owner, as designated by governance

### Lifetime

Persists for system lifecycle. Superseded ADRs are archived, not deleted.

### Authority

Derives authority from architecture, which derives authority from knowledge.

### Dependencies

- Architecture artifacts (required)
- Knowledge artifacts (indirect)

### Deliverables

- Decision title and identifier
- Status (proposed, accepted, deprecated, superseded)
- Context and constraints
- Decision
- Rationale
- Consequences (positive, negative, neutral)
- Related decisions

### Characteristics

- Immutable once accepted (superseded decisions are not modified)
- Traceable to knowledge and architecture
- Versioned
- Reviewed and approved

## Artifact Type 4: Implementation

### Purpose

Realize architecture through code and related artifacts.

### Description

Implementation artifacts represent the physical realization of the system. Implementation artifacts include code, configuration, scripts, and related files.

### Owner

Development Team, as designated by governance

### Lifetime

Persists through system lifecycle. May be replaced through evolution.

### Authority

Derives authority from architecture.

### Dependencies

- Architecture artifacts (required)
- ADRs (when decisions are relevant)
- Knowledge artifacts (indirect, through architecture)

### Deliverables

- Source code
- Configuration files
- Build scripts
- Deployment artifacts
- Implementation-specific documentation

### Characteristics

- Traceable to architecture
- Subject to verification
- Version controlled
- Must not contradict architecture

## Artifact Type 5: Verification

### Purpose

Confirm that implementation aligns with architecture and that architecture aligns with knowledge.

### Description

Verification artifacts document the results of verification activities. Verification artifacts confirm that artifacts at lower levels of the hierarchy satisfy requirements at higher levels.

### Owner

Verification Team, as designated by governance

### Lifetime

Persists for audit purposes. Verification artifacts are archived, not deleted.

### Authority

Verification authority derives from knowledge and architecture. Verification does not create new authority but confirms existing authority alignment.

### Dependencies

- Knowledge artifacts (for verification criteria)
- Architecture artifacts (for verification criteria)
- Implementation artifacts (for what is verified)

### Deliverables

- Verification plans
- Test cases
- Test results
- Verification reports
- Gap analysis
- Non-conformance reports

### Characteristics

- Objective
- Traceable to knowledge and architecture
- Reproducible (where applicable)
- Documented and retained

## Artifact Type 6: Governance

### Purpose

Establish authority, ownership, and process for all other artifact types.

### Description

Governance artifacts define the rules by which other artifacts are created, reviewed, approved, changed, and retired. Governance artifacts establish ownership, authority delegation, and process compliance requirements.

### Owner

Governing Body, as designated by organizational governance

### Lifetime

Persists for the duration of the engineering initiative. May evolve through defined change processes.

### Authority

Governance artifacts derive authority from organizational governance. Governance authorizes the processes that produce all other artifacts.

### Dependencies

- Organizational governance (external)
- All other artifact types (governance applies to all)

### Deliverables

- Ownership assignments
- Authority delegations
- Process definitions
- Review criteria
- Approval workflows
- Change management procedures
- Compliance requirements

### Characteristics

- Authoritative over process
- Stable but changeable through defined procedures
- Applied consistently across artifact types

## Artifact Dependency Graph

```
                    ┌─────────────┐
                    │ Governance  │
                    └──────┬──────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
           ▼               ▼               ▼
      ┌─────────┐     ┌───────────┐   ┌─────────────┐
      │Knowledge│────▶│Architecture│──▶│Implementation│
      └─────────┘     └─────┬─────┘   └──────┬──────┘
                             │                │
                             ▼                ▼
                         ┌─────────┐   ┌─────────────┐
                         │   ADR   │   │ Verification│
                         └─────────┘   └─────────────┘
```

## Artifact Characteristics Summary

| Artifact | Purpose | Authority Level | Originates |
|----------|---------|-----------------|------------|
| Governance | Process authority | Meta | Organizational |
| Knowledge | Problem understanding | Highest | Stakeholders |
| Architecture | Structural decisions | High | Knowledge |
| ADR | Decision documentation | Medium | Architecture |
| Implementation | Physical realization | Low | Architecture |
| Verification | Alignment confirmation | Applied | All above |
