# KDSE Phase 2.1 Methodology Correction

## Reference Artifact Management Separation

**Document Version:** 1.0  
**Effective Date:** 2026-07-13  
**Type:** Methodology Correction

---

## Executive Summary

This document describes a methodology correction that separates **Reference Artifact Management** from the **Collector** responsibility. This correction eliminates conceptual ambiguity in the KDSE methodology.

### What Changed

| Before | After |
|--------|-------|
| "Collector" was expected to discover AND analyze Reference Artifacts | Reference Artifact Management handles discovery; Collector consumes cataloged artifacts |
| Lifecycle did not explicitly show Reference Artifact Management | Lifecycle now explicitly includes Reference Artifact Management as a distinct phase |
| "What does kdse collect do?" had no clear answer | "kds e collect" discovers artifacts; Collector derives knowledge |

### Why This Correction Was Necessary

The Phase 2 methodology introduced the concept of a "Collector" with responsibilities for:
1. Analyzing Reference Artifacts
2. Deriving Domain Knowledge
3. Correlating evidence
4. Identifying contradictions

However, the previous implementation conflated **managing evidence** (discovery, cataloging, classification, provenance) with **analyzing evidence** (deriving knowledge). This created ambiguity:

- "What does kdse collect do?" could mean either managing artifacts or analyzing them
- The distinction between discovery and analysis was unclear
- Implementation could choose either interpretation

The audit correctly identified this as a methodology issue, not an implementation issue.

---

## The Conceptual Problem

### Before: Conflated Responsibilities

```
┌─────────────────────────────────────┐
│            Collector                │
│                                     │
│  - Discover artifacts               │
│  - Catalog artifacts               │
│  - Analyze artifacts               │
│  - Derive knowledge               │
│  - Correlate evidence             │
└─────────────────────────────────────┘
```

**Problem:** "Collector" did too many things. Discovery and analysis are fundamentally different activities.

### After: Clear Separation

```
┌─────────────────────────────────────┐
│  Reference Artifact Management      │
│                                     │
│  - Discover artifacts              │
│  - Catalog artifacts               │
│  - Classify artifacts             │
│  - Preserve provenance            │
└─────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────┐
│            Collector                 │
│                                     │
│  - Analyze cataloged artifacts     │
│  - Derive knowledge               │
│  - Correlate evidence             │
│  - Identify gaps                  │
└─────────────────────────────────────┘
```

**Solution:** Each responsibility has a distinct owner with clear boundaries.

---

## What is Reference Artifact Management?

### Definition

**Reference Artifact Management** is responsible for managing engineering evidence before analysis begins.

**Definition:** The methodology phase responsible for discovering, cataloging, classifying, and maintaining Reference Artifacts throughout their lifecycle.

### Responsibilities

| Responsibility | Description |
|---------------|-------------|
| Discovery | Find Reference Artifacts within the repository and external sources |
| Inventory | Record the existence and basic properties of each artifact |
| Cataloging | Organize Reference Artifacts into meaningful categories |
| Classification | Determine the nature and quality of each artifact |
| Fingerprinting | Establish artifact integrity and identity |
| Provenance | Maintain the origin and history of artifacts |
| Lifecycle Management | Manage artifacts throughout their existence |

### What Reference Artifact Management Does NOT Do

- **Does NOT** interpret artifact content
- **Does NOT** extract knowledge from artifacts
- **Does NOT** derive Domain Knowledge
- **Does NOT** assess Evidence Strength
- **Does NOT** identify contradictions

These responsibilities belong to the Collector.

---

## What is a Collector?

### Updated Definition

A **Collector** is a methodology component that consumes cataloged Reference Artifacts from Reference Artifact Management.

### Responsibilities

| Responsibility | Description |
|---------------|-------------|
| Consume cataloged artifacts | Receive artifacts from Reference Artifact Management |
| Perform Reference Analysis | Examine artifacts to identify evidence |
| Derive Domain Knowledge | Transform evidence into implementation-independent knowledge |
| Apply Engineering Independence Test | Validate that knowledge is technology-independent |
| Correlate evidence | Strengthen knowledge through multiple sources |
| Detect contradictions | Identify and preserve conflicts |
| Validate knowledge | Ensure quality criteria are met |
| Identify gaps | Find areas where knowledge is incomplete |
| Classify questions | Route unresolved items to appropriate phases |

### What Collectors Do NOT Do

- **Do NOT** discover Reference Artifacts
- **Do NOT** catalog or classify artifacts
- **Do NOT** maintain artifact inventory
- **Do NOT** establish provenance

---

## Updated Lifecycle

### The Knowledge Derivation Lifecycle

```
Reference Artifact
        │
        ▼
┌───────────────────────────┐
│  Reference Artifact       │
│  Management               │
│                           │
│  - Discovery             │
│  - Inventory             │
│  - Cataloging            │
│  - Classification        │
│  - Provenance            │
│  - Lifecycle             │
└───────────────────────────┘
        │
        ▼
┌───────────────────────────┐
│  Reference Analysis       │
│                           │
│  - Examine artifacts     │
│  - Identify evidence     │
│  - Document context      │
└───────────────────────────┘
        │
        ▼
┌───────────────────────────┐
│  Domain Knowledge        │
│  Derivation              │
│                           │
│  - Apply independence     │
│    test                  │
│  - Formulate statements  │
│  - Assess strength       │
└───────────────────────────┘
        │
        ▼
┌───────────────────────────┐
│  Evidence Correlation     │
│                           │
│  - Compare sources        │
│  - Identify agreements    │
│  - Preserve conflicts    │
└───────────────────────────┘
        │
        ▼
┌───────────────────────────┐
│  Knowledge Validation     │
│                           │
│  - Review traceability    │
│  - Verify independence    │
│  - Check completeness     │
└───────────────────────────┘
        │
        ▼
  Approved Domain Knowledge
```

---

## The Handoff

The handoff between Reference Artifact Management and Collector is well-defined:

| Reference Artifact Management Produces | Collector Consumes |
|--------------------------------------|-------------------|
| Artifact inventory | Artifact inventory |
| Classification metadata | Classification metadata |
| Provenance records | Provenance records |
| Integrity fingerprints | Integrity fingerprints |
| **NOT analyzed content** | **Analyzed content** |

---

## Command Responsibility

### Before: Ambiguous

The phrase "`kdse collect`" was ambiguous:
- Did it collect (discover) artifacts?
- Did it collect (derive) knowledge?
- Did it do both?

### After: Clear

The methodology now defines:

| Phase | Methodology Responsibility | Command (Implementation) |
|-------|--------------------------|------------------------|
| Discovery | Reference Artifact Management | kdse inventory (or similar) |
| Analysis | Collector | Collector workflow |

**Note:** Whether these are one command or two is an **implementation decision**, not a methodology requirement.

---

## Migration Guide

### For Methodology Users

1. **Understand the separation:**
   - Reference Artifact Management handles discovery and cataloging
   - Collector handles analysis and knowledge derivation

2. **Ask the right question:**
   - "What artifacts exist?" → Reference Artifact Management
   - "What knowledge can be derived?" → Collector

3. **Expect clear outputs:**
   - Reference Artifact Management → Artifact inventory
   - Collector → Domain Knowledge artifacts

### For Implementation Teams

1. **Update documentation:**
   - Replace conflated Collector descriptions with clear separation
   - Update lifecycle diagrams

2. **Update command documentation:**
   - Clarify that commands may implement one or both phases
   - Document the handoff between phases

3. **Update audit criteria:**
   - Reference Artifact Management has distinct criteria
   - Collector has distinct criteria

---

## Summary of Changes

### Documents Updated

| Document | Change |
|----------|--------|
| 004-engineering-model.md | Added Stage 1: Reference Artifact Management |
| 022-collector-philosophy.md | Updated to reflect consumption of cataloged artifacts |
| 007-glossary.md | Added Reference Artifact Management definition |
| README.md | Updated lifecycle and concepts |

### Documents Created

| Document | Purpose |
|----------|---------|
| 025-reference-artifact-management.md | Define Reference Artifact Management responsibility |

### Terminology Changes

| Before | After |
|--------|-------|
| Collector discovers and analyzes | Reference Artifact Management discovers; Collector consumes cataloged artifacts |
| "Knowledge Extraction" | "Knowledge Derivation" (already corrected in Phase 2) |

---

## Lessons Learned

### What Went Wrong

1. **Conflation:** We combined two distinct responsibilities under one name
2. **Ambiguity:** The phrase "Collector" could mean discovery, analysis, or both
3. **Implementation confusion:** Teams couldn't determine which responsibility they were implementing

### How We Fixed It

1. **Separation:** Defined Reference Artifact Management as distinct from Collector
2. **Clarity:** Each phase has exactly one primary responsibility
3. **Handoff:** Defined clear handoff between phases

### Prevention

Before introducing new methodology concepts:

1. Define the responsibility boundary clearly
2. Identify what the concept does NOT include
3. Document the handoff with other concepts
4. Ensure single-responsibility principle is maintained

---

## Related Documents

| Document | Relationship |
|----------|-------------|
| 004-engineering-model.md | Defines the lifecycle with Reference Artifact Management |
| 022-collector-philosophy.md | Defines Collector responsibility |
| 025-reference-artifact-management.md | Defines Reference Artifact Management |
| 007-glossary.md | Contains terminology definitions |

---

## Version History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-07-13 | KDSE Methodology Team | Initial migration notes |

---

*This document explains a methodology correction. It does not modify the KDSE methodology itself.*
