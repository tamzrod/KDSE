# Traceability

**Type:** Template Structure  
**Version:** 1.0

---

## Purpose

The Traceability directory tracks all engineering decisions, their rationale, and their impact on the project.

---

## Structure

```
traceability/
├── README.md                    # This file
├── decision-log/                # Decision records
│   └── README.md
├── requirement-map/              # Requirement traceability
│   └── README.md
├── assumption-registry/          # Assumption tracking
│   └── README.md
└── index.yaml                   # Traceability index
```

---

## Decision Log

All significant decisions are recorded:

```yaml
# decision-log/{date}-{id}.yaml
decision:
  id: "DEC-001"
  date: "2026-07-15"
  title: "Decision title"
  status: "accepted"
  
  context: |
    What prompted this decision
    
  options_considered:
    - "Option 1"
    - "Option 2"
    
  decision: "Chosen option"
  
  rationale: |
    Why this option was chosen
    
  knowledge_basis:
    - "Reference to supporting knowledge"
    
  evidence:
    - "Evidence supporting decision"
    
  impact:
    - "What this affects"
```

---

## Requirement Map

Maps requirements to artifacts:

| Requirement | Design | Implementation | Test |
|-------------|--------|----------------|------|
| REQ-001 | DES-001 | IMP-001 | TST-001 |

---

## Assumption Registry

Tracks assumptions:

```yaml
# assumption-registry/{id}.yaml
assumption:
  id: "ASM-001"
  description: "Assumption statement"
  category: "technical|business|environmental"
  
  risk_if_false: |
    Impact if assumption is incorrect
    
  mitigation: |
    How to reduce risk
    
  validation: |
    How to verify assumption
```

---

## Traceability Index

The index provides a summary:

```yaml
# index.yaml
traceability_index:
  version: "1.0"
  updated: "2026-07-15"
  
  counts:
    decisions: 0
    requirements: 0
    assumptions: 0
    
  coverage:
    traceable_requirements: 0
    total_requirements: 0
```

---

*This directory is copied from the KDSE template during initialization.*
