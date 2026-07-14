# ADR-001: Development Experience as First-Class Artifact

**Date:** 2026-07-14
**Status:** Accepted
**Type:** Methodology Enhancement

---

## Context

KDSE currently defines five artifact types:
1. Knowledge
2. Architecture
3. ADR
4. Implementation
5. Verification

However, KDSE lacks a mechanism to capture verified development knowledge gained during runtime sessions. AI agents repeatedly encounter the same solvable problems across sessions, wasting time rediscovering solutions.

**The problem:** Development knowledge is lost between sessions.

**The opportunity:** Each solved problem in one session could prevent the same problem in future sessions.

---

## Decision

Introduce a sixth artifact type: **Development Experience**.

A Development Experience captures verified knowledge from actual development events, enabling knowledge transfer between AI sessions.

---

## Consequences

### Positive

- Development failures need not be repeated across sessions
- AI agents can learn from past problem-solving
- Development velocity improves over time
- Institutional knowledge persists beyond individual sessions

### Negative

- Additional documentation burden
- Schema must remain simple for adoption
- Evidence must be real, not hypothetical
- Quality depends on capture discipline

---

## Alternatives Considered

| Alternative | Rejected Because |
|------------|------------------|
| Add to Verification artifacts | Verification captures results, not learning |
| Add to Knowledge artifacts | Knowledge is domain-focused, not development-focused |
| Add to Implementation artifacts | Implementation is code, not process |
| External wiki/troubleshooting | Not integrated into KDSE flow |

---

## Schema

```yaml
experience:
  situation: string       # When to reuse (context, not command)
  solution: string       # What to do
  verify: string         # How to confirm
  tags: string[]         # Discovery keywords
  confidence: enum?      # Experimental|Low|Medium|High|Proven
  evidence: string[]?    # Commits, screenshots, logs
```

---

## Knowledge Architecture Position

Development Experience is **lateral** to the main artifact chain:

```
Knowledge → Architecture → ADR → Implementation → Verification

Development Experience (available at any phase)
```

---

## See Also

- [Development Experience Specification](../foundation/DEVELOPMENT-EXPERIENCE-SPECIFICATION.md)
- [Confidence Model](../foundation/CONFIDENCE-MODEL.md)
- [Capture Rules](../foundation/CAPTURE-RULES.md)
- [Retrieval Rules](../foundation/RETRIEVAL-RULES.md)
