# Knowledge Derivation

**Type:** Operational Knowledge  
**Template:** Baseline  
**Version:** 1.0

---

## Purpose

This document defines the process for deriving new knowledge from engineering activities. Knowledge derivation ensures that insights, patterns, and lessons learned are captured systematically.

---

## Knowledge Derivation Process

```
┌─────────────────────────────────────────────────────────────────────┐
│                   KNOWLEDGE DERIVATION PROCESS                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ENGINEERING ACTIVITY                                               │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 1. OBSERVE                                                    │   │
│  │    - What happened?                                           │   │
│  │    - What was the outcome?                                    │   │
│  │    - What context existed?                                    │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 2. ANALYZE                                                    │   │
│  │    - What patterns emerged?                                  │   │
│  │    - What principles apply?                                   │   │
│  │    - What can be generalized?                                  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 3. ABSTRACT                                                   │   │
│  │    - Extract essential elements                               │   │
│  │    - Remove project-specific details                          │   │
│  │    - Identify applicable context                               │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 4. FORMALIZE                                                  │   │
│  │    - Write according to template                              │   │
│  │    - Include examples                                        │   │
│  │    - Add references                                           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 5. VALIDATE                                                   │   │
│  │    - Submit to Laboratory                                    │   │
│  │    - Receive validation results                               │   │
│  │    - Revise if needed                                         │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ 6. INTEGRATE                                                  │   │
│  │    - Add to knowledge base                                    │   │
│  │    - Update indexes                                           │   │
│  │    - Update project digest                                    │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │                                                             │
│       ▼                                                             │
│  NEW KNOWLEDGE                                                      │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Knowledge Sources

### From Sessions

| Source | Knowledge Type |
|--------|---------------|
| Problem solving | Operational knowledge |
| Decision making | Decision rationale |
| Pattern recognition | Developmental patterns |
| Error recovery | Anti-patterns |

### From Projects

| Source | Knowledge Type |
|--------|---------------|
| Architecture decisions | Architecture patterns |
| Implementation approaches | Implementation patterns |
| Testing strategies | Testing patterns |
| Documentation formats | Documentation patterns |

### From Failures

| Source | Knowledge Type |
|--------|---------------|
| Bug analysis | Root cause patterns |
| Incident postmortems | Risk patterns |
| Project retrospectives | Process improvements |

---

## Knowledge Classification

### General Knowledge

Applies across all projects:

| Criterion | Example |
|-----------|---------|
| Universal applicability | "Always validate inputs" |
| Technology agnostic | "Document assumptions" |
| Process universal | "Review before merge" |

### Operational Knowledge

Applies to KDSE workflow:

| Criterion | Example |
|-----------|---------|
| KDSE-specific | "Use Knowledge Discovery for new domains" |
| Session-related | "Review digest after each session" |
| Tool-specific | "Run lab validation before approval" |

### Developmental Knowledge

Applies to software development:

| Criterion | Example |
|-----------|---------|
| Language-specific | "Use type hints in Python" |
| Framework-specific | "Use dependency injection in Angular" |
| Domain-specific | "Implement circuit breaker for microservices" |

---

## Knowledge Template

```yaml
# {category}/{topic}.md
# {type}-knowledge.md

**Type:** [General/Operational/Developmental] Knowledge  
**Template:** Baseline  
**Version:** 1.0  
**Source:** [Session ID, Project, External]  
**Derived:** [Date]

---

## Knowledge Statement

[Clear, actionable statement of the knowledge]

## Rationale

[Why this knowledge matters]

## Context

[When this knowledge applies]

## Examples

### Good Example

[Example of applying this knowledge]

### Poor Example

[Example of not applying this knowledge]

## Related Knowledge

- [Link to related knowledge]

## References

- [External references if any]

---

## Metadata

- Created: [Date]
- Validated: [Date, Laboratory ID]
- Status: [active/deprecated]
```

---

## Quality Criteria

| Criterion | Description |
|-----------|-------------|
| Actionable | Can be directly applied |
| Contextual | Clear when to apply |
| Complete | Includes rationale and examples |
| Validated | Passed Laboratory validation |
| Indexed | Properly categorized and linked |

---

## Anti-Patterns

### Insufficient Generalization

```markdown
# BAD: Too specific
"This API endpoint returned a 500 error because the database
connection pool was exhausted when we had more than 100 concurrent
users on a Tuesday afternoon during the product launch event."

# GOOD: Generalized
"Database connection pools can exhaust under high concurrency.
Always configure appropriate pool size and implement connection
retry logic with exponential backoff."
```

### Missing Context

```markdown
# BAD: No context
"Always use pagination."

# GOOD: With context
"When returning collections that may exceed memory limits or
network MTU, implement pagination to prevent timeouts and ensure
reliable data transfer. Applicable to list endpoints returning
user-generated content or historical records."
```

---

*This document is operational knowledge. It guides the knowledge derivation process.*
