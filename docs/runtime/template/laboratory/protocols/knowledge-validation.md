# Knowledge Validation Protocol

**Protocol ID:** LAB-KNOW-001  
**Version:** 1.0  
**Purpose:** Validate knowledge before integration

---

## Overview

This protocol validates new knowledge to ensure it meets quality standards before being added to the knowledge base.

## Validation Steps

### 1. Syntax Validation

Check that the knowledge document is syntactically correct:

| Check | Expected | Status |
|-------|----------|--------|
| Valid Markdown | Parses without error | |
| Valid YAML frontmatter | Parses without error | |
| Required sections present | Yes | |
| File naming correct | kebab-case.md | |

### 2. Semantic Validation

Check that the knowledge is semantically correct:

| Check | Expected | Status |
|-------|----------|--------|
| Title matches content | Yes | |
| Type is valid | general/operational/developmental | |
| Examples are valid | Runnable or verifiable | |
| Cross-references exist | Links to related knowledge | |

### 3. Completeness Validation

Check that the knowledge is complete:

| Check | Expected | Status |
|-------|----------|--------|
| Purpose section | Present and clear | |
| Rationale section | Present and clear | |
| Context section | Present and clear | |
| At least one example | Present | |

### 4. Quality Validation

Check quality standards:

| Check | Expected | Status |
|-------|----------|--------|
| No placeholder text | "TODO", "FIXME" absent | |
| Consistent terminology | Terms used consistently | |
| Appropriate complexity | Not too simple or complex | |
| Actionable | Can be applied | |

### 5. Uniqueness Validation

Check that the knowledge is not duplicate:

| Check | Expected | Status |
|-------|----------|--------|
| No duplicate title | Unique | |
| No duplicate content | Substantially different | |
| No conflicting knowledge | Consistent with existing | |

---

## Validation Result

```yaml
validation:
  protocol: "LAB-KNOW-001"
  artifact: "{path}"
  timestamp: "{timestamp}"
  
  results:
    syntax: "pass|fail"
    semantic: "pass|fail"
    completeness: "pass|fail"
    quality: "pass|fail"
    uniqueness: "pass|fail"
    
  overall: "pass|fail"
  
  findings:
    - type: "error|warning|info"
      location: "section or line"
      message: "description"
```

---

## Decision

| Result | Action |
|--------|--------|
| Pass | Integrate into knowledge base |
| Fail | Return to author with feedback |
| Conditional | Integrate with notes |

---

*This protocol is part of the KDSE Laboratory.*
