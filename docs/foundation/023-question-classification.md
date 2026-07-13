# Question Classification

## Purpose

This document establishes a framework for **Question Classification** in KDSE. Before asking the operator, classify every unresolved item to ensure it is addressed by the appropriate methodology phase.

## The Need for Classification

During the knowledge derivation process, questions arise that cannot be answered from available Reference Artifacts.

These questions fall into different categories, each with a different resolution path.

### Without Classification

- Questions are asked at random
- Questions are addressed by the wrong phase
- Operator is asked unnecessary questions
- Methodology discipline breaks down

### With Classification

- Questions are routed correctly
- Each phase addresses appropriate questions
- Operator questions are minimized
- Methodology integrity is maintained

## Question Classification Framework

### Classification Categories

| Classification | Description | Resolution Path |
|----------------|-------------|-----------------|
| Engineering Knowledge Question | Cannot be derived from available Reference Artifacts | Ask during Knowledge Derivation |
| Architecture Question | Relates to software organization | Defer to Architecture Phase |
| Implementation Question | Relates to implementation technology | Attempt repository discovery first |

## Category 1: Engineering Knowledge Questions

### Definition

An **Engineering Knowledge Question** is a question that cannot be derived from available Reference Artifacts and relates to engineering understanding.

### Characteristics

Engineering Knowledge Questions:

1. **Require understanding** of engineering purpose, behavior, or constraints
2. **Cannot be answered** from existing documentation or implementation
3. **Affect Engineering Knowledge** derivation
4. **Are essential** for downstream decisions

### Examples

**Question**: "What is the system's fault response time requirement?"

**Analysis**:
- Is this about engineering behavior? Yes
- Can it be derived from artifacts? No (not documented)
- Does it affect Engineering Knowledge? Yes

**Classification**: Engineering Knowledge Question

**Resolution**: Ask during Knowledge Derivation

---

**Question**: "What safety functions does the system implement?"

**Analysis**:
- Is this about engineering purpose? Yes
- Can it be derived from artifacts? No (partially documented)
- Does it affect Engineering Knowledge? Yes

**Classification**: Engineering Knowledge Question

**Resolution**: Ask during Knowledge Derivation

---

**Question**: "What are the operating mode transitions?"

**Analysis**:
- Is this about system behavior? Yes
- Can it be derived from artifacts? No (not specified)
- Does it affect Engineering Knowledge? Yes

**Classification**: Engineering Knowledge Question

**Resolution**: Ask during Knowledge Derivation

### Resolution Process

```
Engineering Knowledge Question Identified
        ↓
Document the question
        ↓
Assess impact on Engineering Knowledge
        ↓
Formulate clear question for operator
        ↓
Ask operator during Knowledge Derivation
        ↓
Incorporate answer into Engineering Knowledge
```

## Category 2: Architecture Questions

### Definition

An **Architecture Question** is a question that relates to software organization, structure, or component relationships.

### Characteristics

Architecture Questions:

1. **Relate to software structure** rather than engineering purpose
2. **Concern component organization** rather than system behavior
3. **Affect Architecture** rather than Engineering Knowledge
4. **Can be deferred** until Architecture Phase

### Examples

**Question**: "How should the control service be organized?"

**Analysis**:
- Is this about engineering purpose? No
- Is this about software organization? Yes
- Can it be deferred? Yes (Architecture Phase)

**Classification**: Architecture Question

**Resolution**: Defer to Architecture Phase

---

**Question**: "Should the historian be a separate service?"

**Analysis**:
- Is this about engineering meaning? No
- Is this about service boundaries? Yes
- Can it be deferred? Yes (Architecture Phase)

**Classification**: Architecture Question

**Resolution**: Defer to Architecture Phase

---

**Question**: "What communication pattern should services use?"

**Analysis**:
- Is this about engineering information? No
- Is this about communication architecture? Yes
- Can it be deferred? Yes (Architecture Phase)

**Classification**: Architecture Question

**Resolution**: Defer to Architecture Phase

### Resolution Process

```
Architecture Question Identified
        ↓
Document the question
        ↓
Confirm it relates to organization, not engineering
        ↓
Flag for Architecture Phase
        ↓
Defer resolution
```

## Category 3: Implementation Questions

### Definition

An **Implementation Question** is a question that relates to implementation technology, technology choices, or specific realization details.

### Characteristics

Implementation Questions:

1. **Concern technology choices** rather than engineering understanding
2. **Relate to specific implementations** rather than organizational structure
3. **Can often be answered** from repository discovery
4. **Should be deferred** if not answerable from artifacts

### Examples

**Question**: "What protocol does the sensor use for communication?"

**Analysis**:
- Is this about engineering information? No
- Is this about implementation technology? Yes
- Can it be discovered? Yes (from implementation)

**Classification**: Implementation Question

**Resolution**: Repository discovery first, defer if unresolved

---

**Question**: "What Modbus registers contain the measurements?"

**Analysis**:
- Is this about engineering meaning? No
- Is this about specific register addresses? Yes
- Can it be discovered? Yes (from implementation)

**Classification**: Implementation Question

**Resolution**: Repository discovery first, defer if unresolved

---

**Question**: "What database does the historian use?"

**Analysis**:
- Is this about engineering requirements? No
- Is this about technology choice? Yes
- Can it be discovered? Yes (from configuration)

**Classification**: Implementation Question

**Resolution**: Repository discovery first, defer if unresolved

### Resolution Process

```
Implementation Question Identified
        ↓
Attempt repository discovery
        │
├── If answer found ──→ Incorporate into Implementation findings
│
└── If answer not found ──→ Flag for Implementation Phase
```

## Classification Decision Tree

Use this decision tree to classify questions:

```
Question Identified
        │
        ▼
Does this question ask about engineering purpose, behavior,
constraints, or intent?
        │
        ├── YES ──→ Engineering Knowledge Question
        │           Ask during Knowledge Derivation
        │
        └── NO
            │
            ▼
        Does this question ask about software organization,
        component structure, or communication patterns?
            │
            ├── YES ──→ Architecture Question
            │           Defer to Architecture Phase
            │
            └── NO
                │
                ▼
            Does this question ask about specific technology,
            protocol, or implementation detail?
                │
                ├── YES ──→ Implementation Question
                │           Attempt repository discovery
                │
                └── NO
                    │
                    ▼
                Reclassify or seek clarification
```

## The Repository First Principle

Before asking the operator, attempt to answer Implementation Questions through repository discovery.

### Sources to Search

Search the following sources before requesting operator input:

1. **Project documentation**: P&IDs, specifications, manuals
2. **Existing implementation**: Source code, configurations
3. **Diagrams**: Architecture diagrams, sequence diagrams
4. **Configuration files**: Environment configs, deployment configs
5. **Source code**: Comments, variable names, constants
6. **Vendor documentation**: API docs, protocol specs
7. **Engineering calculations**: Design docs, sizing calcs
8. **Commissioning records**: Test results, commissioning reports
9. **Existing knowledge artifacts**: Previous analysis results

### Discovery Process

```
Implementation Question
        │
        ▼
Search project documentation
        │
        ▼
Search implementation artifacts
        │
        ▼
Search configuration
        │
        ▼
Search vendor documentation
        │
        ├── Answer found ──→ Incorporate finding
        │
        └── Answer not found ──→ Defer to Implementation Phase
```

## Minimizing Operator Questions

The methodology shall minimize operator questions. Only questions with high engineering value shall be presented.

### Question Value Assessment

Before asking the operator, assess question value:

| Assessment | Questions |
|------------|-----------|
| Is this essential for Engineering Knowledge? | Ask if yes |
| Can this be deferred to a later phase? | Defer if yes |
| Does this affect safety or reliability? | Prioritize if yes |
| Is the answer available from artifacts? | Don't ask if yes |

### Value Threshold

Ask the operator only when:

1. **The question is essential** for Engineering Knowledge derivation
2. **The answer cannot be found** from available artifacts
3. **The answer cannot be deferred** to a later phase
4. **The answer affects** downstream decisions

### Example: Value Assessment

**Question**: "What is the grid code compliance level?"

**Assessment**:

- Essential for Engineering Knowledge? Yes (requirements)
- Available from artifacts? Partially (grid code reference exists)
- Can be deferred? No (affects core requirements)
- Value: HIGH → Ask operator

---

**Question**: "What color is the status indicator light?"

**Assessment**:

- Essential for Engineering Knowledge? No (operational detail)
- Available from artifacts? Likely (vendor docs)
- Can be deferred? Yes (not critical)
- Value: LOW → Don't ask

---

**Question**: "What MQTT topic structure is used?"

**Assessment**:

- Essential for Engineering Knowledge? No (implementation detail)
- Available from artifacts? Yes (implementation)
- Can be deferred? Yes (Architecture/Implementation Phase)
- Value: LOW → Don't ask

## Question Documentation

All questions shall be documented regardless of resolution path.

### Question Log Entry

```
| Q-ID | Question | Classification | Resolution Path | Status |
|------|----------|---------------|----------------|--------|
| Q-001 | Fault response time | Engineering Knowledge | Ask operator | Pending |
| Q-002 | Service boundaries | Architecture | Defer | Deferred |
| Q-003 | Modbus registers | Implementation | Repository search | In Progress |
```

### Question Detail Entry

```
## Q-001: Fault Response Time

**Question**: What is the required fault response time?

**Classification**: Engineering Knowledge Question

**Rationale**: Fault response time is an engineering requirement that
affects system behavior specification.

**Impact**: High - affects core Engineering Knowledge

**Resolution Path**: Ask during Knowledge Derivation

**Status**: Pending operator input

**Answer**: [To be filled when received]
```

## Common Classification Errors

### Error 1: Asking Implementation Questions

**Incorrect**:
> "What Modbus registers are used?" → Asked during Knowledge Derivation

**Why Incorrect**: This is an implementation detail, not Engineering Knowledge.

**Correct**:
> "What measurements does the sensor provide?" → Engineering Knowledge Question

### Error 2: Deferring Engineering Knowledge Questions

**Incorrect**:
> "Fault response time will be addressed in Architecture" → Deferring Engineering Knowledge

**Why Incorrect**: Fault response time is an engineering requirement, not architecture.

**Correct**:
> "Fault response time is an Engineering Knowledge Question" → Address during Knowledge Derivation

### Error 3: Not Attempting Repository Discovery

**Incorrect**:
> "What protocol does the sensor use?" → Asked immediately

**Why Incorrect**: Implementation questions should be attempted from artifacts first.

**Correct**:
> "What protocol does the sensor use?" → Search implementation first, ask if not found

## Summary

Question Classification ensures that:

- **Engineering Knowledge Questions** are asked during Knowledge Derivation
- **Architecture Questions** are deferred to the Architecture Phase
- **Implementation Questions** are attempted through repository discovery first
- **Operator questions are minimized** to high-value items only
- **All questions are documented** regardless of resolution path

Understanding classification is essential for maintaining methodology discipline.

---

## Version

- **Document Version**: 1.0
- **Effective Date**: 2026-07-13
- **Change Note**: Initial release establishing Question Classification framework
