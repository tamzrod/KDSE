# Documentation Standards

**Type:** General Knowledge  
**Template:** Baseline  
**Version:** 1.0

---

## Purpose

This document defines the documentation standards for KDSE-enabled projects. Consistent documentation enables knowledge sharing, reduces onboarding time, and ensures project continuity.

---

## Documentation Types

### Project Documentation

| Type | Purpose | Location |
|------|---------|----------|
| README | Project overview | Repository root |
| Architecture | System design | .kdse/foundation/architecture/ |
| Decisions | Decision records | .kdse/traceability/decision-log/ |
| Knowledge | Engineering knowledge | .kdse/knowledge/ |

### Engineering Documentation

| Type | Purpose | Location |
|------|---------|----------|
| Requirements | Feature specifications | .kdse/foundation/requirements/ |
| Designs | Technical specifications | .kdse/foundation/designs/ |
| APIs | Interface documentation | docs/api/ |
| Runbooks | Operational procedures | docs/runbooks/ |

### Quality Documentation

| Type | Purpose | Location |
|------|---------|----------|
| Audits | Compliance verification | .kdse/reports/audits/ |
| Test Plans | Testing strategy | .kdse/foundation/testing/ |
| Incident Reports | Problem analysis | .kdse/reports/incidents/ |

---

## Documentation Structure

### Document Header

Every document must include:

```yaml
---
title: "Document Title"
type: "engineering-knowledge"
version: "1.0"
created: "2026-07-15"
updated: "2026-07-15"
author: "Engineering Team"
status: "active"
---
```

### Document Sections

| Section | Required | Description |
|---------|----------|-------------|
| Purpose | Yes | Why this document exists |
| Scope | Yes | What this document covers |
| Definitions | No | Term definitions |
| Details | Yes | Main content |
| References | No | Related documents |
| History | No | Change log |

---

## Quality Criteria

### Completeness

| Criterion | Check |
|-----------|-------|
| All required sections present | Yes/No |
| No placeholder text | Yes/No |
| Examples provided where needed | Yes/No |

### Correctness

| Criterion | Check |
|-----------|-------|
| Facts accurate | Yes/No |
| Cross-references valid | Yes/No |
| Code examples runnable | Yes/No |

### Consistency

| Criterion | Check |
|-----------|-------|
| Terminology consistent | Yes/No |
| Format consistent | Yes/No |
| Style consistent | Yes/No |

---

## Review Process

```
┌─────────────────────────────────────────────────────────────────────┐
│                    DOCUMENT REVIEW PROCESS                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  DRAFT                                                              │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ TECHNICAL REVIEW                                            │   │
│  │  - Accuracy check                                           │   │
│  │  - Completeness check                                       │   │
│  │  - Consistency check                                        │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │ Pass                                                     │
│       ▼                                                           │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ STAKEHOLDER REVIEW                                          │   │
│  │  - Relevance check                                          │   │
│  │  - Usability check                                          │   │
│  │  - Approval                                                 │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │ Pass                                                     │
│       ▼                                                           │
│  PUBLISHED                                                         │
│       │                                                           │
│       ▼                                                           │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ PERIODIC REVIEW                                             │   │
│  │  - Currency check                                           │   │
│  │  - Relevance check                                          │   │
│  │  - Update if needed                                         │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Writing Guidelines

### Voice and Tone

| Aspect | Guideline |
|--------|-----------|
| Voice | Active preferred |
| Tone | Professional, clear |
| Jargon | Define when first used |
| Complexity | Prefer simple over complex |

### Structure

| Aspect | Guideline |
|--------|-----------|
| Paragraphs | One idea per paragraph |
| Lists | Use for parallel items |
| Headings | Logical hierarchy |
| Tables | Use for comparisons |

### Formatting

| Element | Standard |
|---------|----------|
| File names | kebab-case.md |
| Code | Monospace |
| Emphasis | Bold for critical |
| Links | Descriptive text |

---

*This document is baseline knowledge. Project-specific documentation standards may be added in .kdse/foundation/documentation/.*
