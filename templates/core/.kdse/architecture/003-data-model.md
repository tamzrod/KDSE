# Data Model

## Purpose

This document defines the system data structures and relationships.

## Entity Relationship Diagram

```
┌──────────────┐       ┌──────────────┐
│   Entity A   │───────│   Entity B   │
└──────────────┘       └──────────────┘
       │
       │
┌──────────────┐
│   Entity C   │
└──────────────┘
```

## Entities

### Entity A

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | Unique identifier |
| name | string | NOT NULL | Entity name |
| created_at | timestamp | | Creation time |

### Entity B

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | Unique identifier |
| entity_a_id | UUID | FK | Reference to Entity A |
| value | number | | Numeric value |

## Relationships

| From | To | Type | Description |
|------|----|------|-------------|
| Entity A | Entity B | 1:N | One A has many B |
| Entity A | Entity C | 1:1 | One A has one C |

## Data Storage

| Entity | Storage | Retention |
|--------|---------|-----------|
| Entity A | PostgreSQL | 7 years |
| Entity B | PostgreSQL | 7 years |

---

**Related Documents**:
- [Components](002-components.md)
- [Knowledge: Requirements](../knowledge/003-requirements.md)

**Evidence Sources**: [List source artifacts]
