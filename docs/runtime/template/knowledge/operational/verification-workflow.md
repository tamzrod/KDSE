# Verification Workflow

**Type:** Operational Knowledge  
**Template:** Baseline  
**Version:** 1.0

---

## Purpose

This document defines the verification workflow for KDSE-enabled projects. Verification ensures that all engineering artifacts meet quality standards before being incorporated into the project.

---

## Verification Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                      VERIFICATION WORKFLOW                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ARTIFACT CREATED                                                   │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ SELF-VERIFICATION                                             │   │
│  │  - Check against template                                    │   │
│  │  - Verify completeness                                       │   │
│  │  - Run basic checks                                          │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │ Pass                                                       │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ LABORATORY SUBMISSION                                        │   │
│  │  - Submit to Laboratory                                       │   │
│  │  - Run validation protocols                                   │   │
│  │  - Receive validation report                                 │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │ Pass                                                       │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ KNOWLEDGE INTEGRATION                                        │   │
│  │  - Update knowledge base                                     │   │
│  │  - Index for discovery                                       │   │
│  │  - Update project digest                                     │   │
│  └─────────────────────────────────────────────────────────────┘   │
│       │                                                             │
│       ▼                                                             │
│  ARTIFACT APPROVED                                                  │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Verification Levels

### Level 1: Self-Verification

Performed by the artifact author:

| Check | Description |
|-------|-------------|
| Template compliance | Follows document template |
| Completeness | All required sections present |
| Format | Correct file naming and structure |
| Links | Internal references valid |

### Level 2: Laboratory Validation

Automated and semi-automated checks:

| Check | Description |
|-------|-------------|
| Syntax | Valid YAML, Markdown, code |
| Semantics | Correct structure and relationships |
| Patterns | Matches established patterns |
| Cross-references | External references valid |

### Level 3: Expert Review

Human review for complex artifacts:

| Check | Description |
|-------|-------------|
| Correctness | Technically accurate |
| Applicability | Appropriate for context |
| Value | Adds value to knowledge base |

---

## Verification Commands

### Self-Verification

```bash
# Verify single artifact
kdse verify --artifact {path}

# Verify all artifacts in category
kdse verify --category knowledge

# Verify project state
kdse verify --project
```

### Laboratory Validation

```bash
# Submit to Laboratory
kdse lab submit --artifact {path}

# Run all validation protocols
kdse lab validate --all

# Run specific protocol
kdse lab validate --protocol knowledge-validation
```

### Full Verification

```bash
# Complete verification workflow
kdse verify --full {artifact}
```

---

## Verification Criteria by Type

### Knowledge Verification

| Criterion | Method |
|-----------|--------|
| Valid Markdown/YAML | Parser validation |
| Required sections | Template check |
| Cross-references | Link validation |
| Classification | Category check |
| Examples | Syntax validation |

### Decision Verification

| Criterion | Method |
|-----------|--------|
| Context provided | Field check |
| Options considered | Field check |
| Rationale provided | Content check |
| Evidence attached | Link validation |
| Impact documented | Field check |

### Evidence Verification

| Criterion | Method |
|-----------|--------|
| Artifact exists | File check |
| Relevance documented | Field check |
| Classification correct | Taxonomy check |
| Timestamp valid | Date format check |

---

## Verification Results

### Pass

```yaml
verification:
  result: "pass"
  artifact: "knowledge/patterns/circuit-breaker.md"
  checks:
    - name: "template"
      status: "pass"
    - name: "completeness"
      status: "pass"
    - name: "syntax"
      status: "pass"
    - name: "references"
      status: "pass"
  timestamp: "2026-07-15T14:08:53Z"
```

### Fail

```yaml
verification:
  result: "fail"
  artifact: "knowledge/patterns/circuit-breaker.md"
  checks:
    - name: "template"
      status: "pass"
    - name: "completeness"
      status: "fail"
      details: "Missing 'examples' section"
    - name: "syntax"
      status: "pass"
    - name: "references"
      status: "pass"
  timestamp: "2026-07-15T14:08:53Z"
  recommendations:
    - "Add examples section with at least one positive example"
```

---

## Integration with Other Processes

### Knowledge Derivation

```
Derivation → Verification → Integration
     │             │
     └─────────────┴──► Project Digest
```

### Experience Capture

```
Experience → Verification → Knowledge Base
     │             │
     └─────────────┴──► Project Digest
```

### Reverse Pull

```
Pull Request → Verification → Integration
     │             │
     └─────────────┴──► Project Digest
```

---

## Verification Schedule

| Verification | Frequency |
|-------------|----------|
| Self-verification | On artifact creation |
| Laboratory validation | On submission |
| Full project verification | On session completion |
| Periodic review | Weekly |

---

*This document is operational knowledge. It guides the verification workflow.*
