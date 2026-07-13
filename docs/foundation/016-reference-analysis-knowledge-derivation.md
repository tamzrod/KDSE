# Reference Analysis and Knowledge Derivation

## Purpose

This document describes the process by which Reference Artifacts are transformed into Engineering Knowledge. It replaces the concept of "Knowledge Extraction" with a more accurate description of the analytical process.

The phrase "Knowledge Extraction" is misleading. Collectors do not extract knowledge. Collectors analyze Reference Artifacts. Engineering Knowledge is derived from engineering evidence.

## Why "Knowledge Extraction" Is Misleading

### Problems with the Extraction Metaphor

The term "Knowledge Extraction" suggests:

1. **Passive Collection**: Knowledge exists somewhere and is simply extracted, like mining ore
2. **Complete Information**: The extraction retrieves the complete knowledge without transformation
3. **Direct Transfer**: Reference Artifacts contain knowledge that transfers directly to artifacts

### The Reality

The reality is different:

1. **Active Analysis**: Reference Artifacts require active analysis to identify engineering evidence
2. **Incomplete Information**: Reference Artifacts contain fragments that must be synthesized
3. **Transformative Process**: Raw evidence must be transformed into implementation-independent statements

### The Correct Framing

The correct framing is:

- **Collectors analyze Reference Artifacts**
- **Analysis identifies engineering evidence**
- **Evidence is correlated across multiple sources**
- **Engineering Knowledge is derived through reasoning**
- **Derived knowledge is validated before authorization**

## The Knowledge Derivation Lifecycle

The transformation from Reference Artifacts to Engineering Knowledge follows a defined lifecycle:

```
Reference Artifact
        ↓
Reference Analysis
        ↓
Engineering Knowledge Derivation
        ↓
Evidence Correlation
        ↓
Knowledge Validation
        ↓
Approved Engineering Knowledge
```

### Phase 1: Reference Artifact

The lifecycle begins with Reference Artifacts. These are existing sources of engineering information, including:

- Project documentation
- Implementation artifacts
- Vendor documentation
- Standards and specifications
- Commissioning records
- Communication artifacts

### Phase 2: Reference Analysis

Reference Analysis is the process of examining Reference Artifacts to identify engineering evidence.

#### Activities in Reference Analysis

1. **Artifact Inventory**: Catalog all available Reference Artifacts
2. **Relevance Assessment**: Determine which artifacts are relevant to the engineering domain
3. **Evidence Identification**: Identify factual statements, assertions, decisions, and constraints
4. **Context Capture**: Document the circumstances under which each artifact was created
5. **Completeness Assessment**: Evaluate the completeness of evidence within each artifact

#### Reference Analysis Outputs

- Inventory of relevant Reference Artifacts
- Identified evidence from each artifact
- Context for each piece of evidence
- Assessment of completeness and reliability

### Phase 3: Engineering Knowledge Derivation

Engineering Knowledge Derivation is the process of transforming analyzed evidence into implementation-independent Engineering Knowledge statements.

#### Activities in Engineering Knowledge Derivation

1. **Interpretation**: Determine the engineering meaning of identified evidence
2. **Generalization**: Abstract implementation-specific details into engineering principles
3. **Synthesis**: Combine evidence from multiple sources into coherent knowledge statements
4. **Independence Validation**: Apply the Engineering Independence Test to confirm implementation-independence

#### The Engineering Independence Test

Every derived Engineering Knowledge statement shall pass the following validation:

> "If the implementation were rewritten tomorrow using a different programming language, communication protocol, runtime, framework, vendor, or platform, would this statement still remain true?"

If YES: The statement is Engineering Knowledge.

If NO: The statement is Architecture or Implementation.

#### Engineering Knowledge Derivation Outputs

- Engineering Knowledge statements (implementation-independent)
- Identification of statements that are Architecture or Implementation
- Traceability links to source evidence
- Evidence Strength assessment

### Phase 4: Evidence Correlation

Evidence Correlation is the process of strengthening Engineering Knowledge through multiple independent Reference Artifacts.

#### Activities in Evidence Correlation

1. **Agreement Analysis**: Identify where multiple artifacts agree
2. **Contradiction Analysis**: Identify where artifacts disagree
3. **Gap Analysis**: Identify areas lacking evidence
4. **Strength Assessment**: Assign Evidence Strength based on corroboration

#### Evidence Strength Scale

Evidence Strength reflects engineering support rather than AI certainty:

| Strength | Description | Criteria |
|----------|-------------|----------|
| ★★★★★ | Strong | Supported by Project Documentation, Node-RED, and Vendor Manual |
| ★★★★☆ | Good | Supported by Project Documentation and one additional source |
| ★★★☆☆ | Moderate | Supported by Project Documentation only |
| ★★☆☆☆ | Weak | Supported by Vendor Manual only or single implementation artifact |
| ★☆☆☆☆ | Minimal | Inferred from indirect evidence |

#### Handling Contradictions

When Reference Artifacts disagree:

1. **Preserve the Contradiction**: Do not silently resolve
2. **Document Both Positions**: Record what each artifact claims
3. **Assess Engineering Impact**: Determine significance for system behavior
4. **Classify Resolution Path**: Engineering Knowledge, Architecture, or Implementation question
5. **Request Operator Review**: Only when contradictions affect Engineering Knowledge

### Phase 5: Knowledge Validation

Knowledge Validation is the process of confirming that derived Engineering Knowledge is accurate, complete, and properly structured.

#### Activities in Knowledge Validation

1. **Completeness Check**: Confirm all necessary evidence has been captured
2. **Consistency Check**: Confirm knowledge statements do not contradict each other
3. **Independence Check**: Confirm all statements pass the Engineering Independence Test
4. **Traceability Check**: Confirm all statements trace to Reference Artifacts
5. **Peer Review**: Subject knowledge to expert review

#### Validation Criteria

| Criterion | Question | Response |
|-----------|----------|----------|
| Accuracy | Does the knowledge accurately represent the evidence? | Proceed or revise |
| Completeness | Is the knowledge complete for its intended purpose? | Add or flag gaps |
| Independence | Does the knowledge remain valid if implementation changes? | Proceed or reclassify |
| Traceability | Can every statement trace to evidence? | Add links or flag gaps |

### Phase 6: Approved Engineering Knowledge

Approved Engineering Knowledge has completed the derivation lifecycle and received authorization for use in subsequent phases.

#### Authorization Criteria

Knowledge becomes Approved when:

1. **Evidence is Sufficient**: Reference Analysis identified adequate evidence
2. **Derivation is Sound**: Engineering Knowledge Derivation followed proper process
3. **Correlation is Complete**: Evidence was correlated across sources
4. **Validation Passed**: All validation criteria were satisfied
5. **Review is Complete**: Expert review confirmed accuracy and completeness

#### Approved Engineering Knowledge Characteristics

Approved Engineering Knowledge is:

- **Implementation-Independent**: Valid regardless of implementation technology
- **Traceable**: Links to source Reference Artifacts documented
- **Evidence-Correlated**: Corroborated across multiple sources or contradictions preserved
- **Validated**: Subject to and passed validation process
- **Authorized**: Reviewed and approved by qualified personnel

## Detailed Phase Activities

### Reference Analysis Detailed Activities

#### Step 1: Artifact Inventory

Create an inventory of all available Reference Artifacts:

```
| Artifact ID | Artifact Type | Source | Date | Relevance |
|-------------|---------------|--------|------|-----------|
| RA-001 | Project Documentation | P&ID | 2024-01-15 | High |
| RA-002 | Implementation | Node-RED | 2024-03-20 | High |
| RA-003 | Vendor Manual | ABB | 2023-11-01 | Medium |
```

#### Step 2: Evidence Extraction

Extract evidence from each artifact:

```
| Evidence ID | Artifact ID | Evidence Type | Content |
|-------------|-------------|---------------|---------|
| E-001 | RA-001 | Fact | Inverter rated at 500kW |
| E-002 | RA-001 | Assertion | Supports grid-forming mode |
| E-003 | RA-002 | Decision | Grid-forming enabled via switch |
```

#### Step 3: Context Documentation

Document the context for each piece of evidence:

```
| Evidence ID | Creator | Creation Date | Review Status | Currency |
|-------------|---------|---------------|---------------|----------|
| E-001 | Project Engineer | 2024-01-15 | Reviewed | Current |
| E-002 | Project Engineer | 2024-01-15 | Reviewed | Current |
| E-003 | Commissioning Tech | 2024-03-20 | Not Reviewed | Current |
```

### Engineering Knowledge Derivation Detailed Activities

#### Step 1: Interpretation

Interpret each piece of evidence:

```
| Evidence ID | Interpretation | Confidence | Rationale |
|-------------|---------------|------------|-----------|
| E-001 | System includes 500kW inverter | High | Explicit specification |
| E-002 | Inverter capable of grid-forming operation | High | Explicit statement |
| E-003 | Grid-forming mode is operator-selectable | Medium | Implementation detail, may vary |
```

#### Step 2: Generalization

Generalize interpretations to engineering principles:

```
| Evidence ID | Generalization | Implementation Independence |
|-------------|---------------|---------------------------|
| E-001 | Plant includes power conversion equipment rated for specific capacity | Yes - capacity is engineering specification |
| E-002 | Inverter supports grid-forming control mode | Yes - capability is engineering specification |
| E-003 | Grid-forming mode selection is controlled externally | Yes - control interface is engineering specification |
```

#### Step 3: Engineering Knowledge Statement Formulation

Formulate Engineering Knowledge statements:

```
| Statement ID | Engineering Knowledge Statement | Evidence | Independence |
|--------------|-------------------------------|----------|--------------|
| EK-001 | The plant includes power conversion equipment rated at 500kW | E-001 | Validated |
| EK-002 | The inverter supports grid-forming control mode as a system capability | E-002 | Validated |
| EK-003 | Grid-forming mode selection is controlled through an external interface | E-003 | Validated |
```

### Evidence Correlation Detailed Activities

#### Step 1: Agreement Matrix

Create an agreement matrix:

```
| Statement | RA-001 | RA-002 | RA-003 | Agreement |
|-----------|--------|--------|--------|-----------|
| EK-001 | ✓ | - | ✓ | Corroborated |
| EK-002 | ✓ | ✓ | ✓ | Strongly Corroborated |
| EK-003 | - | ✓ | - | Single Source |
```

#### Step 2: Evidence Strength Assignment

Assign Evidence Strength:

```
| Statement | Sources | Strength | Rationale |
|-----------|---------|----------|-----------|
| EK-001 | RA-001, RA-003 | ★★★★☆ | Project doc + Vendor manual |
| EK-002 | RA-001, RA-002, RA-003 | ★★★★★ | All three sources agree |
| EK-003 | RA-002 | ★★☆☆☆ | Single implementation source |
```

### Knowledge Validation Detailed Activities

#### Step 1: Completeness Checklist

```
| Criterion | Status | Notes |
|-----------|--------|-------|
| All relevant evidence identified | ✓ | RA-001, RA-002, RA-003 |
| Evidence interpretation documented | ✓ | All statements have rationale |
| Engineering knowledge derived | ✓ | EK-001, EK-002, EK-003 |
| Evidence strength assessed | ✓ | Scale applied to all |
| Contradictions preserved | ✓ | None identified |
```

#### Step 2: Consistency Verification

```
| Check | Result | Action |
|-------|--------|--------|
| EK-001 consistent with EK-002? | ✓ | No conflict |
| EK-002 consistent with EK-003? | ✓ | No conflict |
| All statements independent? | ✓ | Passes Engineering Independence Test |
```

## The Role of the Collector

A **Collector** is a methodology component responsible for executing the knowledge derivation lifecycle.

### Collector Responsibilities

A Collector shall:

- Analyze Reference Artifacts
- Identify engineering evidence
- Derive implementation-independent Engineering Knowledge
- Preserve traceability
- Correlate evidence
- Identify contradictions
- Identify Engineering Knowledge gaps

### Collector Limitations

A Collector does not:

- Generate software documentation
- Design software
- Generate architecture
- Implement software

### Collector Definition

A Collector is defined by its responsibility, not by the type of Reference Artifact it analyzes.

Support for specific Reference Artifact formats is an implementation concern.

## Question Classification During Derivation

During the derivation lifecycle, questions arise that cannot be answered from Reference Artifacts.

### Classification Framework

Before asking the operator, classify every unresolved item:

| Classification | Description | Action |
|----------------|-------------|--------|
| Engineering Knowledge Question | Cannot be derived from available Reference Artifacts | Ask during Knowledge Derivation |
| Architecture Question | Relates to software organization | Defer to Architecture Phase |
| Implementation Question | Relates to implementation technology | Attempt repository discovery first; defer if unresolved |

### The Repository First Principle

Before asking the operator:

1. Search project documentation
2. Examine existing implementation
3. Review diagrams
4. Check configuration
5. Analyze source code
6. Consult vendor documentation
7. Review engineering calculations
8. Examine commissioning records
9. Check existing knowledge artifacts

If sufficient engineering evidence exists, derive the Engineering Knowledge. Do not ask the operator.

### Minimizing Operator Questions

The methodology shall minimize operator questions. Only questions with high engineering value shall be presented.

Ask yourself:

1. Is this question essential for Engineering Knowledge?
2. Can this question be answered from available Reference Artifacts?
3. Is this an Architecture or Implementation question that should be deferred?
4. What is the engineering risk if this question remains unanswered?

## Summary

The Knowledge Derivation Lifecycle transforms Reference Artifacts into Approved Engineering Knowledge through systematic process:

1. **Reference Artifact**: Collect existing engineering information
2. **Reference Analysis**: Examine artifacts to identify evidence
3. **Engineering Knowledge Derivation**: Transform evidence into implementation-independent statements
4. **Evidence Correlation**: Strengthen knowledge through multiple sources
5. **Knowledge Validation**: Confirm accuracy and completeness
6. **Approved Engineering Knowledge**: Authorize validated knowledge for use

This process replaces "Knowledge Extraction" with a more accurate description of analytical transformation.

---

## Version

- **Document Version**: 1.0
- **Effective Date**: 2026-07-13
- **Change Note**: Initial release replacing "Knowledge Extraction" with formal derivation lifecycle
