# Authority & Traceability

## Purpose

This document defines how authority flows in KDSE and how traceability enables that authority.

## Authority Hierarchy

Authority flows downward in KDSE:

```
Knowledge (highest authority)
    ├──> Architecture (authorized by knowledge)
    └──> Verification Criteria (derives from knowledge)

Architecture
    ├──> ADRs (document decisions)
    └──> Implementation (guides)

Implementation
    └──> Verification (subject to)
```

**Key rule**: Lower artifacts cannot contradict higher artifacts. Implementation cannot contradict architecture. Architecture cannot contradict knowledge.

## What Traceability Is

Traceability is the ability to follow relationships between artifacts. It enables:
- Understanding why artifacts exist
- Impact analysis when knowledge changes
- Audit and compliance verification

### Traceability vs Documentation

| Documentation | Traceability |
|--------------|--------------|
| Records information | Records relationships |
| Can be incomplete | Requires explicit links |
| May be inaccurate | Can be verified |

### Traceability vs Lineage

| Lineage | Traceability |
|---------|--------------|
| Tracks origin | Tracks justification |
| "Where it came from" | "Why it exists" |

## What Can Be Traced

### Traceability Matrix

| From | To | Relationship | Required |
|------|-----|-------------|----------|
| Knowledge | Architecture | Authorizes | Yes |
| Knowledge | Verification Criteria | Derives | Yes |
| Architecture | ADR | Documents | Yes |
| Architecture | Implementation | Guides | Yes |
| ADR | Implementation | Justifies | Yes |
| Implementation | Verification | Subject to | Yes |
| Verification | Knowledge | Confirms alignment | Yes |

### Forward Traceability

Forward traceability follows relationships from higher-authority to lower-authority artifacts, confirming that lower artifacts derive from higher artifacts.

### Backward Traceability

Backward traceability follows relationships from lower-authority to higher-authority artifacts, confirming that higher artifacts justify lower artifacts.

## Traceability Depth

| Depth | When to Use |
|-------|-------------|
| **Deep** | Regulated industries, safety-critical systems |
| **Moderate** | Typical business systems (default) |
| **Light** | Small projects, rapid development |

## Traceability and Authority

Traceability enables authority by demonstrating justification:

```
Knowledge (authority)
    | Traceability demonstrates
    v
Architecture (inherits authority)
    | Traceability demonstrates
    v
Implementation (inherits authority)
```

**Without traceability:**
- Authority without traceability = assertion without justification
- Traceability without authority = documentation without meaning

## Key Insight

Traceability confirms relationships exist. It does **not** confirm:
- Relationships are complete
- Artifacts are correct
- Decisions are good

Well-traced poor decisions are still poor decisions.

## Verification

Traceability should be verified:
1. When artifacts are created
2. When artifacts change
3. Before release
4. During audit
