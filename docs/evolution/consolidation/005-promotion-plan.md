# KDSE Consolidation: Promotion Plan

**Document Version:** 0.1  
**Date:** 2026-07-13  
**Type:** Promotion Plan  
**Status:** DRAFT - For Operator Review  

---

## Purpose

This document describes what should eventually be promoted into the KDSE Foundation. It provides a structured plan for incorporating the discoveries from the methodology exploration into Foundation documents. **No Foundation documents shall be modified in this document.**

---

## 1. Overview

### 1.1 Promotion Objectives

| Objective | Description |
|-----------|-------------|
| **Preserve Clarity** | Ensure all discoveries are properly captured |
| **Maintain Integrity** | Follow evidence chain model |
| **Enable Validation** | Provide mechanisms to verify implementation |
| **Minimize Disruption** | Avoid breaking existing processes |

### 1.2 Promotion Principles

1. **No Direct Modification:** This document describes what should change, but does not modify Foundation documents
2. **Evidence-Based:** Every promotion traces to documented evidence
3. **Prioritized:** Changes are ordered by importance and dependencies
4. **Verified:** Promotion includes validation criteria

---

## 2. Priority 1: Critical Updates (Required)

These updates address fundamental methodology gaps that should be resolved.

### 2.1 Elevate "Evidence Strengthens but Does Not Authorize"

**Discovery Source:** Phase 2.0 - Evidence vs Authority

**Current State:** Implicit in Evidence Strength documentation

**Recommended State:** Explicit Core Principle

**Action Required:**
- Add to Core Principles (003-core-principles.md) as Principle 16
- State: "Evidence Strength reflects confidence, not authority. Authority derives from structured derivation, not evidence quantity."

**Validation Criteria:**
- Principle appears in Core Principles
- Principle is referenced in Evidence Strength documentation
- Principle is referenced in Authority Resolution documentation

**Dependencies:** None

---

### 2.2 Integrate Reference Artifact Management into Lifecycle

**Discovery Source:** Phase 2.1 - Reference Artifact Management

**Current State:** Reference Artifact Management is implied but not explicitly modeled

**Recommended State:** Explicit Stage 1 in Knowledge Derivation Lifecycle

**Action Required:**
- Update 004-engineering-model.md to include Reference Artifact Management as Stage 1
- Update 016-reference-analysis-knowledge-derivation.md lifecycle diagram
- Add new document: 025-reference-artifact-management.md (already created)

**Validation Criteria:**
- Lifecycle diagram shows Reference Artifact Management as distinct phase
- Handoff contract is documented
- Collector definition references consumption of cataloged artifacts

**Dependencies:** None

---

### 2.3 Add Engineering Independence Test Mechanism

**Discovery Source:** Phase 2.0 - Implementation Independence

**Current State:** Engineering Independence Test is documented in 024-engineering-independence-test.md

**Recommended State:** Referenced from Core Principles and Engineering Knowledge documentation

**Action Required:**
- Add explicit reference to Engineering Independence Test from Principle 12
- Add reference from 009-engineering-knowledge.md
- Ensure test question is consistent across documentation

**Validation Criteria:**
- Principle 12 references the test
- Engineering Knowledge definition references the test
- Test question appears in appropriate documents

**Dependencies:** 2.1

---

## 3. Priority 2: Important Updates (Recommended)

These updates improve methodology clarity and consistency.

### 3.1 Update Glossary with Terminology Changes

**Discovery Source:** Multiple phases

**Terminology Changes:**
| Term | Change | Status |
|------|--------|--------|
| Owner | → Steward | COMPLETE |
| Knowledge Extraction | → Knowledge Derivation | COMPLETE |
| Compliance Score | → Assessment Score (early phase) | COMPLETE |
| Collector discovers | → Collector consumes cataloged artifacts | TO PROMOTE |

**Action Required:**
- Verify Glossary (007-glossary.md) includes all updated terminology
- Add Reference Artifact Management definition if missing
- Add Evidence Strength scale if missing
- Add Engineering Independence Test definition if missing

**Validation Criteria:**
- All terminology changes are reflected in Glossary
- Cross-references are accurate
- Examples are consistent

**Dependencies:** 2.2

---

### 3.2 Clarify Architecture/Implementation Boundary

**Discovery Source:** Phase 2.0 - Architecture/Implementation Separation

**Current State:** Architecture and Implementation are defined separately but boundary is implicit

**Recommended State:** Explicit boundary statement

**Action Required:**
- Add explicit boundary statement to 018-architecture-phase.md
- Add explicit boundary statement to 019-implementation-phase.md
- Statement: "Architecture decisions are organization; Implementation decisions are realization"

**Validation Criteria:**
- Boundary statement appears in both documents
- Examples clarify the distinction
- Chain of Authority references are consistent

**Dependencies:** None

---

### 3.3 Add Repository First as Explicit Principle

**Discovery Source:** Phase 2.0 - Question Classification

**Current State:** Repository First is Principle 14

**Recommended State:** Verify explicit documentation with question classification

**Action Required:**
- Verify 014-engineering-review-process.md references Repository First
- Add explicit "before asking the operator" sequence if missing
- Include question classification as a practice

**Validation Criteria:**
- Repository First is referenced in question handling documentation
- Sequence (search artifacts → analyze → examine docs → review materials) is explicit

**Dependencies:** 2.1

---

### 3.4 Document Phase-Aware Recommendations

**Discovery Source:** Phase 1.5 - Lifecycle Awareness

**Current State:** Phase detection is in execution documentation

**Recommended State:** Reference from Chain of Authority documentation

**Action Required:**
- Verify 006-chain-of-authority.md references phase-appropriate recommendations
- Add note about phase detection in recommendations

**Validation Criteria:**
- Chain of Authority documentation references phase awareness
- Recommendation criteria include phase context

**Dependencies:** None

---

## 4. Priority 3: Supporting Updates (If Time Permits)

These updates provide additional clarity but are not critical.

### 4.1 Document Artifact Lifecycle

**Discovery Source:** Phase 1.5 - Artifact Lifecycle

**Current State:** Artifact lifecycle states are implied

**Recommended State:** Explicit lifecycle states in documentation

**Action Required:**
- Add artifact lifecycle states to 005-engineering-artifacts.md
- Include state transition criteria

**Validation Criteria:**
- Lifecycle states are documented
- Transition criteria are clear

**Dependencies:** None

---

### 4.2 Add Evidence-Driven Debugging Capability

**Discovery Source:** Debug Runtime Session

**Current State:** Debugging workflow is in session workspace

**Recommended State:** Consider as optional capability

**Action Required:**
- If promoting: Add 026-debugging-capability.md to documentation
- If not promoting: Document as separate guidance

**Validation Criteria:**
- If promoted: Capability is integrated into methodology
- If not promoted: Separate guidance is available

**Dependencies:** None

---

## 5. Documents Requiring Updates

### 5.1 Summary Table

| Document | Priority | Changes Required |
|----------|----------|------------------|
| 003-core-principles.md | 1 | Add Principle 16 |
| 004-engineering-model.md | 1 | Add Reference Artifact Management to lifecycle |
| 007-glossary.md | 2 | Add terminology definitions |
| 009-engineering-knowledge.md | 1 | Reference Engineering Independence Test |
| 012-traceability.md | 3 | Add Evidence Strength reference |
| 013-authority-resolution.md | 1 | Reference Evidence Strengthens principle |
| 016-reference-analysis-knowledge-derivation.md | 1 | Update lifecycle diagram |
| 018-architecture-phase.md | 2 | Add boundary statement |
| 019-implementation-phase.md | 2 | Add boundary statement |
| 022-collector-philosophy.md | 1 | Update to consumption model |
| 023-question-classification.md | 2 | Reference Repository First |
| 024-engineering-independence-test.md | 1 | Ensure consistency |

### 5.2 Documents Created (Already Complete)

| Document | Purpose | Status |
|----------|---------|--------|
| 025-reference-artifact-management.md | Define Reference Artifact Management | COMPLETE |

---

## 6. Promotion Sequence

### 6.1 Phase 1: Critical Foundation Updates

| Step | Action | Validation |
|------|--------|------------|
| 1.1 | Add Principle 16 to 003-core-principles.md | Principle visible in list |
| 1.2 | Add Reference Artifact Management to 004-engineering-model.md | Lifecycle shows new stage |
| 1.3 | Update 016-reference-analysis-knowledge-derivation.md | Lifecycle diagram updated |
| 1.4 | Update 022-collector-philosophy.md | Consumption model documented |
| 1.5 | Reference Engineering Independence Test in 009-engineering-knowledge.md | Test referenced |

**Exit Criteria:** Core principles and lifecycle are complete

### 6.2 Phase 2: Important Clarifications

| Step | Action | Validation |
|------|--------|------------|
| 2.1 | Update 007-glossary.md with new terminology | All terms defined |
| 2.2 | Add boundary statements to 018-architecture-phase.md and 019-implementation-phase.md | Boundaries clear |
| 2.3 | Add Repository First reference to 023-question-classification.md | Principle referenced |
| 2.4 | Add Evidence Strengthens reference to 013-authority-resolution.md | Principle referenced |

**Exit Criteria:** Terminology is consistent across documentation

### 6.3 Phase 3: Supporting Updates

| Step | Action | Validation |
|------|--------|------------|
| 3.1 | Add artifact lifecycle to 005-engineering-artifacts.md | States documented |
| 3.2 | Consider debug capability promotion | Decision documented |

**Exit Criteria:** Supporting documentation is complete

---

## 7. Validation Plan

### 7.1 Consistency Validation

| Check | Method | Pass Criteria |
|------|--------|---------------|
| Principle Consistency | Compare all principle references | All references match |
| Terminology Consistency | Search for old terminology | No old terminology remains |
| Lifecycle Consistency | Compare lifecycle diagrams | All diagrams match |
| Cross-Reference Validation | Check all internal links | All links valid |

### 7.2 Completeness Validation

| Check | Method | Pass Criteria |
|------|--------|---------------|
| Discovery Coverage | Compare to 001-discovery.md | All discoveries addressed |
| Principle Coverage | Check all new principles | All principles defined |
| Example Coverage | Verify examples | Examples reflect changes |

### 7.3 Quality Validation

| Check | Method | Pass Criteria |
|------|--------|---------------|
| Readability Review | Human review of updated documents | Clear and understandable |
| Accuracy Review | Expert review of technical content | Technically correct |
| Evidence Review | Verify evidence chain | All changes traceable |

---

## 8. Rollback Plan

If issues are discovered after promotion:

| Scenario | Rollback Action |
|----------|-----------------|
| Principle conflicts | Revert to previous principle list |
| Lifecycle inconsistency | Use previous lifecycle until resolved |
| Terminology confusion | Maintain aliases during transition |
| Structural issues | Revert to previous architecture |

---

## 9. Open Items for Operator Decision

### 9.1 Items Requiring Decision

| Item | Question | Options |
|------|----------|---------|
| Debug Capability | Should evidence-driven debugging be promoted to Foundation? | Promote as 026, Keep as separate, Discard |
| Confidence Thresholds | Should debugging confidence thresholds be part of methodology? | Include, Defer to domain guidance, Discard |
| Phase 0 Initialization | Should bootstrap process be in Foundation or Runtime docs? | Foundation, Runtime only, Both |

### 9.2 Recommended Decisions

| Item | Recommendation | Rationale |
|------|----------------|-----------|
| Debug Capability | Keep as separate guidance | Domain-specific; not universal |
| Confidence Thresholds | Defer to domain guidance | Needs validation across domains |
| Phase 0 Initialization | Runtime only | Implementation-specific detail |

---

## 10. Success Criteria

### 10.1 Minimum Success Criteria

| Criterion | Verification |
|-----------|--------------|
| Principle 16 added | Visible in 003-core-principles.md |
| Reference Artifact Management integrated | Visible in 004-engineering-model.md |
| Engineering Independence Test referenced | Referenced from 009 and 012 |
| Terminology updated | Glossary contains new terms |
| Collector updated | 022 reflects consumption model |

### 10.2 Ideal Success Criteria

| Criterion | Verification |
|-----------|--------------|
| All Priority 1 updates complete | All validation criteria pass |
| All Priority 2 updates complete | All validation criteria pass |
| No inconsistencies identified | Full review complete |
| All discoveries addressed | 001-discovery.md coverage verified |

---

## Version History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 0.1 | 2026-07-13 | KDSE Methodology Team | Initial promotion plan |

---

*This document describes what should be promoted into the Foundation. Foundation documents are NOT modified by this document. Actual modification should follow the promotion sequence after operator approval.*
