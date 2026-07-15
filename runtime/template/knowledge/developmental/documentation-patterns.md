# Documentation Patterns

**Type:** Developmental Knowledge  
**Template:** Baseline  
**Version:** 1.0

---

## Purpose

This document defines patterns for creating effective engineering documentation. Consistent documentation improves readability, maintainability, and knowledge sharing.

---

## Pattern: README Structure

### Template

```markdown
# Project Name

Brief description (1-2 sentences)

## Overview

High-level understanding of the project

## Quick Start

Steps to get running quickly

## Features

- Feature 1
- Feature 2
- Feature 3

## Architecture

Architecture diagram or link to detailed docs

## Development

### Prerequisites

### Setup

### Running

### Testing

## Deployment

## Contributing

## License
```

### When to Use

- Repository root
- Project overview
- Getting started guide

---

## Pattern: Decision Record

### Template

```markdown
# {Number}: {Title}

**Status:** {Proposed|Accepted|Deprecated|Superseded}

**Date:** {YYYY-MM-DD}

## Context

{What is the issue that we're seeing that is motivating this decision?}

## Decision

{What is the decision that we're proposing?}

## Consequences

### Positive
- {Benefit 1}
- {Benefit 2}

### Negative
- {Downside 1}
- {Downside 2}

### Neutral
- {Trade-off 1}
```

### When to Use

- Architectural decisions
- Technology choices
- Process changes

---

## Pattern: API Documentation

### Template

```markdown
## Endpoint Name

**URL:** `/api/v1/resource`

**Method:** GET|POST|PUT|DELETE

### Description

What this endpoint does

### Request

#### Headers

| Header | Required | Description |
|--------|----------|-------------|
| Authorization | Yes | Bearer token |

#### Body (for POST/PUT)

```json
{
  "field": "description"
}
```

### Response

#### Success (200)

```json
{
  "data": {}
}
```

#### Error (400)

```json
{
  "error": "Validation error",
  "details": []
}
```

### Example

```bash
curl -X GET /api/v1/resource -H "Authorization: Bearer token"
```
```

### When to Use

- API endpoints
- Service interfaces
- Integration points

---

## Pattern: Runbook

### Template

```markdown
# Runbook: {Title}

**Author:** {Name}
**Last Updated:** {Date}
**On-Call:** {Link to on-call schedule}

## Purpose

What this runbook is for

## Prerequisites

- Access level required
- Tools needed

## Procedures

### Normal Operations

1. Step 1
2. Step 2

### Incident Response

1. Step 1
2. Step 2

## Troubleshooting

| Symptom | Solution |
|---------|----------|
| Issue 1 | Fix 1 |
| Issue 2 | Fix 2 |

## Rollback

How to rollback if something goes wrong

## Contacts

- Primary: {Name}
- Escalation: {Name}
```

### When to Use

- Operational procedures
- Incident response
- Deployment guides

---

## Pattern: Knowledge Document

### Template

```markdown
# Title

**Type:** General|Operational|Developmental Knowledge
**Version:** 1.0
**Created:** {Date}

---

## Knowledge Statement

The main knowledge being documented

## Rationale

Why this knowledge matters

## Context

When to apply this knowledge

## Examples

### Good

```example
Correct usage
```

### Poor

```example
Incorrect usage
```

## Related

- [Link to related knowledge]

## References

- [External references]
```

### When to Use

- Engineering knowledge
- Best practices
- Lessons learned

---

## Pattern: Architecture Diagram

### ASCII Style

```
┌─────────────────────────────────────────────────────────────────────┐
│                         COMPONENT NAME                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐        │
│  │   Layer A    │────▶│   Layer B    │────▶│   Layer C    │        │
│  └──────────────┘     └──────────────┘     └──────────────┘        │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### When to Use

- System architecture
- Component relationships
- Data flow

---

## Documentation Quality Checklist

| Check | Description |
|-------|-------------|
| Audience | Clear who this is for |
| Purpose | Clear why this exists |
| Completeness | All essential information present |
| Accuracy | Information is correct |
| Currency | Up to date |
| Examples | Real examples provided |
| Maintenance | Clear how to update |

---

*This document is developmental knowledge. Project-specific documentation patterns may be added in .kdse/knowledge/developmental/documentation/.*
