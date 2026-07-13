# KDSE `collect` Methodology Audit

**Audit Date:** 2026-07-13  
**Auditor:** KDSE Methodology Review  
**Status:** Complete

---

## Executive Summary

The current `kdse collect` command implementation is **fundamentally misaligned** with the KDSE methodology as defined in the Phase 2 documentation. The command performs artifact discovery and cataloging, but this responsibility does not match the methodology's definition of a "Collector."

**Key Finding:** `kdse collect` is an artifact inventory generator, not a methodology Collector. The methodology expects Collectors to analyze Reference Artifacts, derive Domain Knowledge, correlate evidence, and identify gaps—none of which the current implementation performs.

**Methodology Alignment Score: 2/10**

---

## Methodology Alignment Score

| Dimension | Score | Assessment |
|-----------|-------|------------|
| Purpose Alignment | 3/10 | Partial - discovers artifacts but doesn't derive knowledge |
| Reference Artifact Handling | 7/10 | Correct - treats inputs as Reference Artifacts |
| Domain Knowledge Derivation | 0/10 | Missing - no knowledge derivation occurs |
| Evidence Handling | 1/10 | Minimal - only catalogs metadata, no evidence analysis |
| Evidence Correlation | 0/10 | Missing - no correlation mechanism exists |
| Domain Interface Derivation | 0/10 | Missing - not addressed |
| Question Classification | 0/10 | Missing - no classification mechanism |
| Engineering Independence Test | 0/10 | Missing - no independence validation |
| Separation of Concerns | 2/10 | Partial - some separation, but Architecture/Implementation blur |
| Command Responsibility | 6/10 | Acceptable - single responsibility, but wrong responsibility |

**Overall Score: 2/10 (Critical Misalignment)**

---

## Findings

### Finding 1: Collector Definition Mismatch (CRITICAL)

**Location:** `internal/collect/collector.go`, `runtime/COLLECT.md`

**Current Behavior:**
- `kdse collect` discovers files in `artifacts/` directory
- Records metadata: path, size, hash, extension, modification time
- Categorizes files by extension and path patterns
- Generates inventory.json and collection report

**Expected Behavior (per `022-collector-philosophy.md`):**
A Collector must:
1. Analyze Reference Artifacts
2. Identify engineering evidence
3. Derive implementation-independent Domain Knowledge
4. Preserve traceability
5. Correlate evidence
6. Identify contradictions
7. Identify Domain Knowledge gaps

**Gap:** The current `kdse collect` performs none of these responsibilities. It is an artifact cataloging tool, not a methodology Collector.

**Impact:** The methodology defines Collectors with specific responsibilities, but no current command implements them. The gap between definition and implementation undermines methodology integrity.

---

### Finding 2: Evidence Pipeline Uses Outdated Terminology (HIGH)

**Location:** `runtime/COLLECT.md` lines 64-86

**Current Pipeline (from COLLECT.md):**
```
Engineering Evidence
        ↓
   kdse collect
        ↓
  Artifact Inventory
        ↓
     Executor
        ↓
Knowledge Extraction
        ↓
 Normalized Markdown
        ↓
 Operator Review
        ↓
Approved Knowledge
```

**Expected Pipeline (per `016-reference-analysis-knowledge-derivation.md`):**
```
Reference Artifact
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
```

**Gap:** The Evidence Pipeline uses "Knowledge Extraction" which the methodology explicitly replaces with "Knowledge Derivation." The pipeline stages do not match the defined Knowledge Derivation Lifecycle.

**Impact:** Documentation inconsistency creates confusion about methodology flow.

---

### Finding 3: Domain Knowledge Derivation Not Implemented (CRITICAL)

**Location:** `internal/collect/`, `internal/normalize/`

**Expected Behavior:**
Collectors derive implementation-independent Domain Knowledge from Reference Artifacts.

**Required Derivations:**
- Domain purpose and behavior
- Operating modes
- Engineering constraints
- Engineering assumptions
- Safety behavior
- Control philosophy
- Engineering state machines
- Domain Interfaces

**Current Behavior:**
`kdse normalize` extracts requirements, decisions, glossary terms, and constraints, but:
- Does not validate implementation-independence
- Does not filter programming language details
- Does not filter protocol details (Modbus, MQTT, etc.)
- Does not filter vendor-specific implementations
- Includes implementation details as knowledge

**Gap:** Domain Knowledge derivation with implementation-independence validation is not implemented.

---

### Finding 4: Evidence Strength Not Assessed (HIGH)

**Location:** `internal/collect/types.go`, `internal/normalize/types.go`

**Expected Behavior (per `021-evidence-and-strength.md`):**
Evidence Strength is determined by independent supporting engineering evidence:
- ★★★★★: Supported by multiple independent sources
- ★★★★☆: Supported by Project Doc + one additional source
- ★★★☆☆: Supported by Project Documentation only
- ★★☆☆☆: Supported by single source or vendor only
- ★☆☆☆☆: Inferred from indirect evidence

**Current Behavior:**
- `DiscoveredArtifact` has no Evidence Strength field
- `KnowledgeExtraction` has no Evidence Strength per knowledge item
- No correlation mechanism to strengthen knowledge with multiple sources

**Gap:** Evidence Strength assessment is not implemented.

---

### Finding 5: Evidence Correlation Not Implemented (HIGH)

**Location:** `internal/collect/`, `internal/normalize/`

**Expected Behavior (per `022-collector-philosophy.md`):**
Collectors strengthen knowledge through multiple sources:
- Compare evidence across artifacts
- Identify agreements
- Identify contradictions
- Assign Evidence Strength

**Current Behavior:**
- No mechanism to compare evidence across multiple Reference Artifacts
- No contradiction detection
- No evidence strengthening through correlation

**Gap:** Evidence correlation is not implemented.

---

### Finding 6: Contradiction Detection Not Implemented (HIGH)

**Location:** `internal/collect/`, `internal/normalize/`

**Expected Behavior (per `022-collector-philosophy.md`):**
Collectors preserve contradictions rather than resolving them:
- Document conflicting claims
- Assess engineering impact
- Classify resolution path
- Flag for operator review

**Current Behavior:**
- No contradiction detection mechanism
- No preservation of conflicts
- Silent resolution of any "conflicts" through deduplication

**Gap:** Contradiction detection and preservation is not implemented.

---

### Finding 7: Question Classification Not Implemented (MEDIUM)

**Location:** `internal/collect/`, `internal/normalize/`

**Expected Behavior (per `023-question-classification.md`):**
Before asking the operator, classify every unresolved item:
| Classification | Resolution Path |
|----------------|-----------------|
| Domain Knowledge Question | Ask during Knowledge Derivation |
| Architecture Question | Defer to Architecture Phase |
| Implementation Question | Attempt repository discovery first |

**Current Behavior:**
- No question classification mechanism
- No routing to appropriate phases
- No Repository First Principle implementation

**Gap:** Question classification is not implemented.

---

### Finding 8: Engineering Independence Test Not Implemented (MEDIUM)

**Location:** `internal/collect/`, `internal/normalize/`

**Expected Behavior (per `024-engineering-independence-test.md`):**
Every derived Domain Knowledge statement shall pass validation:
> "If the implementation were rewritten tomorrow using a different programming language, communication protocol, runtime, framework, vendor, or platform, would this statement still remain true?"

**Current Behavior:**
- No independence test validation
- Implementation details may appear in derived knowledge
- No filtering of programming language, protocol, or vendor specifics

**Gap:** Engineering Independence Test is not implemented.

---

### Finding 9: Domain Interface Derivation Not Implemented (MEDIUM)

**Location:** `internal/collect/`, `internal/normalize/`

**Expected Behavior (per `020-domain-interfaces.md`):**
Domain Interfaces are derived from Domain Knowledge:
- Define what information is exchanged
- Define what the information means
- Exclude communication technologies (REST, gRPC, MQTT, Modbus)

**Current Behavior:**
- No Domain Interface derivation
- Implementation details may be mixed with interface definitions
- No separation between knowledge and interface

**Gap:** Domain Interface derivation is not implemented.

---

### Finding 10: Architecture/Implementation Separation Blurred (MEDIUM)

**Location:** `internal/normalize/extraction.go`, `internal/normalize/generator.go`

**Expected Behavior:**
- Architecture explains how Domain Knowledge is organized into software
- Implementation explains how Architecture is realized using specific technologies

**Current Behavior:**
`kdse normalize` mixes concerns:
- Extracts "API endpoint" as a requirement (Implementation)
- Generates Architecture documents from unstructured content
- Does not separate organizational decisions from technology choices

**Gap:** Architecture/Implementation separation is not enforced.

---

### Finding 11: Repository First Principle Not Implemented (LOW)

**Location:** `internal/collect/`, `internal/normalize/`

**Expected Behavior (per `023-question-classification.md`):**
Before asking the operator:
1. Search all available Reference Artifacts
2. Analyze existing implementation
3. Examine project documentation
4. Review vendor materials

**Current Behavior:**
- No systematic repository search before questions
- No evidence that all artifacts were analyzed
- No documentation of repository analysis

**Gap:** Repository First Principle is not implemented.

---

## Outdated Concepts

### Outdated Concept 1: "Knowledge Extraction"

**Location:** `runtime/COLLECT.md`, `internal/normalize/extraction.go`

**Issue:** The methodology explicitly replaces "Knowledge Extraction" with "Knowledge Derivation."

**Rationale:** The phrase "Knowledge Extraction" is misleading. Collectors do not extract knowledge—they analyze Reference Artifacts and derive Domain Knowledge.

**Recommendation:** Replace all instances of "Knowledge Extraction" with "Knowledge Derivation."

---

### Outdated Concept 2: AI Confidence Scoring

**Location:** `runtime/COLLECT.md`, `internal/normalize/types.go`

**Issue:** The methodology replaces "AI confidence" with "Evidence Strength."

**Rationale:** Evidence Strength reflects engineering support rather than AI certainty.

**Recommendation:** Replace confidence scoring with Evidence Strength assessment based on supporting engineering evidence.

---

### Outdated Concept 3: Artifact-Type-Defined Collectors

**Location:** `runtime/COLLECT.md` artifact categories

**Issue:** Collectors should be defined by responsibility, not artifact type.

**Current:**
- "A requirements collector gathers requirements"
- "A design collector gathers design documents"
- "A code collector gathers source code"

**Expected (per `022-collector-philosophy.md`):**
A Collector analyzes any Reference Artifact that contains engineering evidence, regardless of type.

**Recommendation:** Reframe artifact categories as Reference Artifact sources, not collector types.

---

## Missing Concepts

### Missing Concept 1: Knowledge Derivation Lifecycle

**Not implemented:**
```
Reference Artifact
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
```

**Required:** Implement each phase as a distinct workflow step.

---

### Missing Concept 2: Evidence Correlation Pipeline

**Not implemented:**
- Cross-artifact evidence comparison
- Agreement detection
- Contradiction identification
- Evidence Strength assignment

**Required:** Implement evidence correlation between multiple Reference Artifacts.

---

### Missing Concept 3: Question Classification Pipeline

**Not implemented:**
- Classification of unresolved items
- Routing to appropriate phases
- Repository-first search
- Operator question minimization

**Required:** Implement question classification before operator interaction.

---

### Missing Concept 4: Engineering Independence Test

**Not implemented:**
- Validation that derived knowledge is implementation-independent
- Filtering of programming language, protocol, vendor details
- Pass/fail validation for each knowledge statement

**Required:** Implement Engineering Independence Test for all derived knowledge.

---

### Missing Concept 5: Domain Interface Derivation

**Not implemented:**
- Derivation of Domain Interfaces from Domain Knowledge
- Technology-agnostic interface definitions
- Separation of interface from implementation

**Required:** Implement Domain Interface derivation as part of knowledge processing.

---

## Responsibility Audit

### Current Command Responsibilities

| Command | Responsibility | Alignment |
|---------|-----------------|-----------|
| `kdse collect` | Discover and catalog files | ❌ Wrong responsibility |
| `kdse normalize` | Analyze docs, extract requirements | ⚠️ Partial |
| (missing) | Reference Analysis | ❌ Not implemented |
| (missing) | Knowledge Derivation | ❌ Not implemented |
| (missing) | Evidence Correlation | ❌ Not implemented |
| (missing) | Knowledge Validation | ❌ Not implemented |
| (missing) | Question Classification | ❌ Not implemented |

### Assessment

**Single Responsibility Principle: PASS**
Each command has a single responsibility, but the responsibilities are wrong for the methodology.

**Collector Responsibility (per methodology): FAIL**
The methodology defines a Collector with specific responsibilities that no current command implements.

---

## Separation of Concerns Audit

### Required Separation

```
Reference Artifact (evidence)
        ↓
Domain Knowledge (implementation-independent understanding)
        ↓
Architecture (organization into software)
        ↓
Implementation (specific technology realization)
```

### Current State

| Layer | Separation Status | Issues |
|-------|------------------|--------|
| Reference Artifacts | ✓ Implemented | Correct discovery and cataloging |
| Domain Knowledge | ❌ Not implemented | No implementation-independent derivation |
| Architecture | ⚠️ Partial | Mixed with Implementation in normalize |
| Implementation | ⚠️ Partial | Mixed with Architecture in normalize |

### Bleed-Through Issues

1. **Implementation in Knowledge:**
   - `kdse normalize` extracts "API endpoint: GET /api/v1/..." as a requirement
   - Implementation details appear in generated artifacts

2. **Architecture in Implementation:**
   - No clear separation between organizational decisions and technology choices
   - Generated documents mix concerns

3. **Evidence in Knowledge:**
   - Reference Artifact content may be copied verbatim
   - Evidence is not preserved separately from derived knowledge

---

## Recommendations

### Recommendation 1: Rename `kdse collect` Responsibility (CRITICAL)

**Current:** `kdse collect` discovers and catalogs files in `artifacts/` directory.

**Proposed:** Rename to `kdse inventory` or `kdse discover` and clarify its responsibility as "Reference Artifact discovery and cataloging."

**Rationale:** The current name suggests it performs the methodology's Collector role, which it does not.

---

### Recommendation 2: Create Collector Implementation (CRITICAL)

**Action:** Implement a new command or workflow that performs the Collector responsibility:

```
Reference Artifact Discovery
        ↓
Reference Artifact Inventory
        ↓
Reference Analysis
        ↓
Evidence Identification
        ↓
Domain Knowledge Derivation
        ↓
Engineering Independence Validation
        ↓
Evidence Correlation
        ↓
Contradiction Detection
        ↓
Evidence Strength Assessment
        ↓
Domain Knowledge Artifacts
        ↓
Gap Identification
        ↓
Question Classification
        ↓
Operator Questions (classified, minimal)
```

---

### Recommendation 3: Implement Evidence Correlation (HIGH)

**Action:** Add evidence correlation to the workflow:

1. Compare evidence across multiple Reference Artifacts
2. Identify agreements (strengthens knowledge)
3. Identify contradictions (preserve, don't resolve)
4. Assign Evidence Strength based on corroboration

---

### Recommendation 4: Implement Engineering Independence Test (MEDIUM)

**Action:** Add independence validation to knowledge derivation:

For each derived knowledge statement, ask:
> "If the implementation were rewritten tomorrow using a different programming language, communication protocol, runtime, framework, vendor, or platform, would this statement still remain true?"

**If NO:** Filter the statement as Architecture or Implementation, not Domain Knowledge.

---

### Recommendation 5: Implement Question Classification (MEDIUM)

**Action:** Before asking operator, classify unresolved items:

| Classification | Action |
|----------------|--------|
| Domain Knowledge Question | Ask during Knowledge Derivation |
| Architecture Question | Defer to Architecture Phase |
| Implementation Question | Repository discovery first |

**Repository First:** Search all available artifacts before asking questions.

---

### Recommendation 6: Update COLLECT.md Pipeline (HIGH)

**Current:**
```
Engineering Evidence → kdse collect → Artifact Inventory → Executor → Knowledge Extraction → ...
```

**Proposed:**
```
Reference Artifacts → kdse inventory → Reference Artifact Inventory → Collector Workflow → Domain Knowledge Artifacts → ...
```

---

### Recommendation 7: Separate Architecture and Implementation (MEDIUM)

**Action:** Create distinct phases for:
- **Architecture Phase:** How Domain Knowledge is organized into software components
- **Implementation Phase:** How Architecture is realized using specific technologies

**Validation:** Ensure no implementation details appear in Architecture artifacts.

---

## Proposed Internal Workflow for `kdse collect`

*Note: Renaming may be appropriate. The following describes the methodology-accurate workflow.*

### Phase 1: Reference Artifact Discovery

```
1.1 Scan repository for Reference Artifacts
1.2 Record artifact metadata (path, hash, category, timestamp)
1.3 Generate artifact inventory
1.4 Preserve artifacts unchanged
```

### Phase 2: Reference Analysis

```
2.1 Analyze each Reference Artifact
2.2 Identify factual statements
2.3 Identify assertions and claims
2.4 Identify constraints and decisions
2.5 Document evidence provenance
```

### Phase 3: Domain Knowledge Derivation

```
3.1 Transform evidence into Domain Knowledge
3.2 Apply Engineering Independence Test
3.3 Filter implementation details (language, protocol, vendor)
3.4 Formulate implementation-independent statements
3.5 Ensure knowledge remains valid across technology changes
```

### Phase 4: Evidence Correlation

```
4.1 Compare evidence across multiple artifacts
4.2 Identify corroborating evidence
4.3 Identify contradicting evidence
4.4 Assign Evidence Strength (★★★★★ to ★☆☆☆☆)
4.5 Preserve contradictions (do not resolve silently)
```

### Phase 5: Knowledge Validation

```
5.1 Validate traceability links
5.2 Verify independence test passes
5.3 Confirm evidence strength
5.4 Check completeness
5.5 Identify gaps
```

### Phase 6: Question Management

```
6.1 Classify unresolved items
6.2 Apply Repository First Principle
6.3 Defer Architecture Questions
6.4 Attempt repository discovery for Implementation Questions
6.5 Generate minimal, high-value operator questions
```

### Output Artifacts

```
- Reference Artifact Inventory
- Evidence Analysis Report
- Domain Knowledge Artifacts (validated)
- Evidence Correlation Report
- Gap Report
- Question Log (classified)
```

---

## Conclusion

The current `kdse collect` implementation does not fulfill the methodology's Collector responsibility. The command discovers and catalogs files, which is valuable but does not constitute knowledge derivation as defined by KDSE.

**Critical Actions Required:**

1. Clarify `kdse collect` responsibility as "Reference Artifact discovery and cataloging" (not knowledge derivation)
2. Implement the Collector responsibility as a distinct workflow or command
3. Implement Evidence Correlation
4. Implement Engineering Independence Test
5. Implement Question Classification
6. Update COLLECT.md documentation to reflect methodology-accurate terminology and pipeline

**Timeline Recommendation:** Before any additional implementation work, the methodology alignment must be corrected. Implementing features against outdated assumptions wastes effort and introduces technical debt.

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-07-13 | KDSE Audit | Initial audit |

---

*This document is an internal audit. It does not modify KDSE methodology or implementation.*
