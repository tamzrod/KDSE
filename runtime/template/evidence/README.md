# Evidence

**Type:** Template Structure  
**Version:** 1.0

---

## Purpose

The Evidence directory stores engineering evidence—verifiable documentation that supports claims, decisions, and validations.

---

## Structure

```
evidence/
├── README.md                    # This file
├── collection/                  # Evidence collection templates
│   └── README.md
├── classification/              # Evidence taxonomy
│   └── README.md
└── index.yaml                   # Evidence index
```

---

## Evidence Types

| Type | Description | Examples |
|------|-------------|----------|
| Document | Written evidence | Specs, reports, emails |
| Metric | Quantitative evidence | Performance data, test results |
| Test | Verification evidence | Unit tests, integration tests |
| Log | Historical evidence | System logs, audit trails |
| Reference | External evidence | Standards, documentation |

---

## Evidence Collection

Evidence is collected during engineering sessions:

1. **Identify** - What evidence supports this claim?
2. **Collect** - Gather the evidence
3. **Classify** - Assign to evidence type
4. **Index** - Record in evidence index
5. **Store** - Place in appropriate directory

---

## Quality Criteria

| Criterion | Description |
|-----------|-------------|
| Verifiable | Can be independently confirmed |
| Relevant | Directly supports the claim |
| Complete | Contains all necessary details |
| Current | Up to date and accurate |

---

*This directory is copied from the KDSE template during initialization.*
