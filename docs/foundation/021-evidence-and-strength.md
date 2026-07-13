# Evidence and Strength

## Purpose

This document establishes **Evidence Strength** as a replacement for AI confidence in KDSE. Evidence Strength reflects engineering support rather than AI certainty. Evidence Strength is determined by independent supporting engineering evidence.

## Why Replace AI Confidence

### Problems with AI Confidence

AI confidence metrics have significant limitations:

1. **Uncertainty About Uncertainty**: AI confidence reflects model uncertainty, not engineering uncertainty
2. **Calibration Issues**: AI confidence often miscalibrated with actual accuracy
3. **Lack of Engineering Context**: AI confidence ignores the nature of supporting evidence
4. **Black Box**: Confidence values are not interpretable in engineering terms

### The Evidence Strength Alternative

Evidence Strength is based on:

1. **Engineering Evidence**: Reference Artifacts that support the knowledge
2. **Independence**: Whether evidence comes from independent sources
3. **Corroboration**: Whether multiple sources agree
4. **Engineering Significance**: What the evidence represents

## Evidence Strength Scale

Evidence Strength reflects engineering support rather than AI certainty.

### Strength Levels

| Strength | Label | Description |
|----------|-------|-------------|
| ★★★★★ | Strong | Supported by multiple independent sources |
| ★★★★☆ | Good | Supported by project documentation and one additional source |
| ★★★☆☆ | Moderate | Supported by project documentation only |
| ★★☆☆☆ | Weak | Supported by single source or vendor documentation only |
| ★☆☆☆☆ | Minimal | Inferred from indirect evidence |

### Evidence Strength Examples

#### ★★★★★ - Strong Evidence

**Supported by**:

- Project Documentation
- Node-RED (Implementation)
- Vendor Manual

**Example**:
> "The inverter supports grid-forming control mode."

**Evidence**:

- Project P&ID specifies grid-forming capability
- Node-RED flow shows grid-forming switch implementation
- Vendor manual confirms grid-forming feature

#### ★★★★☆ - Good Evidence

**Supported by**:

- Project Documentation
- Node-RED (Implementation)

**Example**:
> "Grid-forming mode is selected via operator interface."

**Evidence**:

- Project HMI specification shows mode selection control
- Node-RED flow implements mode switching logic

#### ★★★☆☆ - Moderate Evidence

**Supported by**:

- Project Documentation only

**Example**:
> "The system requires 100ms fault response time."

**Evidence**:

- Project requirements document specifies 100ms response

#### ★★☆☆☆ - Weak Evidence

**Supported by**:

- Vendor Manual only

**Example**:
> "The inverter supports DNP3 communication."

**Evidence**:

- Vendor manual lists DNP3 as supported protocol

#### ★☆☆☆☆ - Minimal Evidence

**Supported by**:

- Inferred from implementation
- Derived from vendor documentation without project confirmation

**Example**:
> "System operates in grid-connected mode by default."

**Evidence**:

- Node-RED flow initializes with grid-connected state
- Inferred: this is the default operating mode

## Evidence Correlation

Engineering Knowledge is strengthened through multiple independent Reference Artifacts.

### Correlation Principles

#### Principle 1: Reference Artifacts Support, Not Replace

Reference Artifacts support Engineering Knowledge. They do not replace it.

The existence of a Reference Artifact does not automatically create Engineering Knowledge. Evidence must be analyzed, interpreted, and validated.

#### Principle 2: Agreement Strengthens

When multiple Reference Artifacts agree, Engineering Knowledge becomes stronger.

| Agreement Pattern | Effect |
|-------------------|--------|
| Multiple independent sources agree | ★★★★★ |
| Project doc + implementation agree | ★★★★☆ |
| Single source | ★★★☆☆ or below |

#### Principle 3: Disagreement Requires Preservation

When Reference Artifacts disagree, the contradiction shall be preserved.

Contradictions shall never be silently resolved.

#### Principle 4: Source Independence Matters

Evidence from independent sources is more valuable than evidence from the same source.

| Independence Pattern | Effect |
|---------------------|--------|
| Project doc + vendor manual + implementation | ★★★★★ |
| Project doc + vendor manual | ★★★★☆ |
| Multiple implementations from same project | ★★★☆☆ |

## Handling Contradictions

### Step 1: Document the Contradiction

Preserve the contradiction in the analysis:

```
| Source A Claims | Source B Claims | Point of Disagreement |
|-----------------|-----------------|----------------------|
| Grid-forming supported | Grid-forming planned for future | Capability status |
```

### Step 2: Assess Engineering Impact

Evaluate the engineering significance:

| Impact Level | Criteria | Action |
|--------------|----------|--------|
| High | Affects safety, reliability, or critical function | Requires resolution |
| Medium | Affects non-critical function | Document and monitor |
| Low | Minor detail | Accept uncertainty |

### Step 3: Determine Resolution Path

Classify the contradiction:

| Classification | Description | Action |
|----------------|-------------|--------|
| Engineering Knowledge Question | Core capability or behavior disputed | Ask during Knowledge Derivation |
| Architecture Question | Organization or structure disputed | Defer to Architecture Phase |
| Implementation Question | Technology choice disputed | Attempt repository discovery |

### Step 4: Request Operator Review

Operator review is required only when contradictions affect Engineering Knowledge.

Operator review may be deferred when contradictions affect only Architecture or Implementation.

## Evidence Collection Checklist

When assessing Evidence Strength, verify:

- [ ] All relevant Reference Artifacts identified?
- [ ] Evidence from multiple independent sources?
- [ ] Agreements documented?
- [ ] Contradictions preserved?
- [ ] Strength rating assigned?

## Evidence Strength and Authority

Evidence Strength affects confidence but not authority.

### Authority Depends on Process

Authority derives from:

1. **Validation**: Knowledge was validated through structured process
2. **Review**: Knowledge was reviewed by qualified personnel
3. **Traceability**: Knowledge traces to evidence
4. **Independence**: Knowledge passes the Engineering Independence Test

### Confidence Depends on Evidence

Evidence Strength affects confidence:

1. **Strong evidence**: High confidence in knowledge accuracy
2. **Weak evidence**: Lower confidence, may require additional validation
3. **Contradictions**: Flagged for resolution, may reduce confidence

### Example: Authority vs Confidence

**High Authority + Strong Evidence**:

> "The system shall support grid-forming mode."
>
> - Validated through review
> - Passed Engineering Independence Test
> - Corroborated by multiple sources
> - Confidence: High

**High Authority + Weak Evidence**:

> "Grid-forming requires external synchronization source."
>
> - Validated through review
> - Passed Engineering Independence Test
> - Supported by vendor manual only
> - Confidence: Medium
> - Recommendation: Verify with additional evidence

## Evidence Assessment Process

### Phase 1: Identify Evidence Sources

List all Reference Artifacts that contain relevant evidence:

```
| Artifact ID | Artifact Type | Relevant Evidence |
|-------------|---------------|------------------|
| RA-001 | Project Documentation | P&ID specifications |
| RA-002 | Implementation | Node-RED flows |
| RA-003 | Vendor Documentation | Product manual |
```

### Phase 2: Extract Evidence

Extract specific evidence from each artifact:

```
| Evidence ID | Artifact ID | Evidence Content |
|-------------|-------------|------------------|
| E-001 | RA-001 | "Grid-forming mode supported" |
| E-002 | RA-002 | Grid-forming switch in flow |
| E-003 | RA-003 | "Grid-forming capability listed" |
```

### Phase 3: Assess Agreement

Compare evidence across sources:

```
| Statement | RA-001 | RA-002 | RA-003 | Agreement |
|-----------|--------|--------|--------|----------|
| Grid-forming supported | ✓ | ✓ | ✓ | Full agreement |
| Grid-forming via switch | ✓ | ✓ | - | Partial agreement |
```

### Phase 4: Assign Strength

Assign Evidence Strength based on assessment:

```
| Statement | Sources | Agreement | Strength |
|-----------|---------|-----------|----------|
| Grid-forming supported | All three | Full | ★★★★★ |
| Grid-forming via switch | Two | Partial | ★★★★☆ |
```

### Phase 5: Document Contradictions

Document any contradictions:

```
| Contradiction | Source A | Source B | Impact |
|---------------|----------|----------|--------|
| Default mode | RA-001: Grid-following | RA-002: Grid-forming | Low |
```

## Evidence Strength Documentation

### In Engineering Knowledge Artifacts

Document Evidence Strength for each knowledge statement:

```
## Engineering Knowledge: Grid-Forming Support

**Statement**: The inverter shall support grid-forming control mode.

**Evidence Strength**: ★★★★★

**Supporting Evidence**:
- RA-001 (Project P&ID): "Grid-forming mode specified"
- RA-002 (Node-RED): Grid-forming switch implementation
- RA-003 (Vendor Manual): Grid-forming feature listed

**Traceability**: EK-001 traces to RA-001, RA-002, RA-003
```

### In Analysis Reports

Document Evidence Strength in analysis reports:

```
## Analysis: Control Mode Implementation

| Statement | Evidence | Strength | Notes |
|-----------|----------|---------|-------|
| Grid-forming supported | RA-001, RA-002, RA-003 | ★★★★★ | Full corroboration |
| Grid-forming default | RA-002 | ★★☆☆☆ | Single source, needs verification |
| Fallback behavior | RA-003 | ★☆☆☆☆ | Inferred, requires project confirmation |
```

## Common Errors

### Error 1: Equating AI Confidence to Evidence Strength

**Incorrect**:
> "AI confidence is 95%, so Evidence Strength is ★★★★★"

**Why Incorrect**: AI confidence reflects model certainty, not engineering evidence.

**Correct**:
> "Three independent sources corroborate this statement, so Evidence Strength is ★★★★★"

### Error 2: Counting Sources Without Independence

**Incorrect**:
> "Three implementation artifacts support this, so Evidence Strength is ★★★★★"

**Why Incorrect**: Multiple artifacts from the same source don't provide independent evidence.

**Correct**:
> "Project documentation, implementation, and vendor manual corroborate, so Evidence Strength is ★★★★★"

### Error 3: Silently Resolving Contradictions

**Incorrect**:
> "Source A and B disagree, but we decided source A is correct"

**Why Incorrect**: Contradictions shall be preserved.

**Correct**:
> "Source A claims X, Source B claims Y. Contradiction preserved. Resolution pending operator review."

### Error 4: Ignoring Weak Evidence

**Incorrect**:
> "Only vendor manual supports this, but we'll treat it as confirmed"

**Why Incorrect**: Weak evidence should be flagged for verification.

**Correct**:
> "Only vendor manual supports this statement. Evidence Strength: ★★☆☆☆. Recommend verification with project documentation."

## Summary

Evidence Strength replaces AI confidence with a measure based on engineering evidence:

- **Reflects** engineering support, not AI certainty
- **Based on** independent supporting Reference Artifacts
- **Corroborated** across multiple sources
- **Preserves** contradictions rather than resolving them
- **Documented** for each Engineering Knowledge statement
- **Separate from** authority, which depends on process

Understanding Evidence Strength is essential for maintaining honest assessment of knowledge confidence.

---

## Version

- **Document Version**: 1.0
- **Effective Date**: 2026-07-13
- **Change Note**: Initial release establishing Evidence Strength as replacement for AI confidence
